package search

import (
	"errors"
	"hash/fnv"
	"sync"

	"os"

	"sort"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util"
)

var (
	once               sync.Once
	riotDictPath       string
	searcher           = riot.Engine{}
	id2url             = map[uint64]string{}
	avail              = false
	documentFreqByword = map[string]int{}
)

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func setDictPathFromEnv() {
	dictPath, ok := os.LookupEnv("RIOT_DICT_PATH")
	if !ok {
		goPath := util.GetEnv("GOPATH", "~/go")
		dictPath = goPath + "/src/github.com/go-ego/gse/data/dict/jp/dict.txt"
	}
	riotDictPath = dictPath
}

func Init() error {
	var err error
	once.Do(func() {
		setDictPathFromEnv()
		searcher.Init(types.EngineOpts{
			GseDict:  riotDictPath,
			UseStore: false,
		})

		examples := example.Examples{}
		positiveExamples, err := db.ReadPositiveExamples(10000)
		if err != nil {
			return
		}
		examples = append(examples, positiveExamples...)
		unlabeledExamples, err := db.ReadUnlabeledExamples(10000)
		if err != nil {
			return
		}
		examples = append(examples, unlabeledExamples...)
		cache.AttachMetadata(examples, false, false)
		for _, e := range examples {
			id := hash(e.FinalUrl)
			id2url[id] = e.FinalUrl
			searcher.Index(id, types.DocData{Content: e.Title})
			for _, w := range getUniqueWords(e.Title) {
				documentFreqByword[w]++
			}
		}
		searcher.Flush()
	})
	if err != nil {
		return err
	}
	avail = true
	return nil
}

func Ping() error {
	if !avail {
		return errors.New("searcher cannot be available")
	}
	return nil
}

func Search(query string) (example.Examples, error) {
	urls := make([]string, 0)
	req := types.SearchReq{
		Text:     query,
		RankOpts: &types.RankOpts{MaxOutputs: 100},
	}
	for _, resp := range searcher.SearchDoc(req).Docs {
		url := id2url[resp.DocId]
		urls = append(urls, url)
	}
	examples, err := db.SearchExamplesByUlrs(urls)
	if err != nil {
		return nil, err
	}
	cache.AttachMetadata(examples, false, true)
	return examples, nil
}

func removeOneCharKeywords(keywords []string) []string {
	result := make([]string, 0)
	for _, k := range keywords {
		if len([]rune(k)) > 1 {
			result = append(result, k)
		}
	}
	return result
}

func getUniqueWords(s string) []string {
	return util.RemoveDuplicate(removeOneCharKeywords(searcher.Segment(s)))
}

func GetKeywordsInQuery(query string) []string {
	tokens := getUniqueWords(query)
	type kv struct {
		Key   string
		Value int
	}

	var dfInQuery []kv
	for _, k := range tokens {
		if cnt, ok := documentFreqByword[k]; ok {
			dfInQuery = append(dfInQuery, kv{k, cnt})
		}
	}

	sort.Slice(dfInQuery, func(i, j int) bool {
		return dfInQuery[i].Value < dfInQuery[j].Value
	})

	result := make([]string, 0)
	for idx, kv := range dfInQuery {
		if idx > 2 {
			break
		}
		result = append(result, kv.Key)
	}

	return result
}

func SearchSimilarExamples(query string) (example.Examples, []string, error) {
	keywords := GetKeywordsInQuery(query)
	req := types.SearchReq{
		Tokens:   keywords,
		Logic:    types.Logic{Should: true},
		RankOpts: &types.RankOpts{MaxOutputs: 10},
	}

	urls := make([]string, 0)
	for _, resp := range searcher.SearchDoc(req).Docs {
		url := id2url[resp.DocId]
		urls = append(urls, url)
	}
	examples, err := db.SearchExamplesByUlrs(urls)
	if err != nil {
		return nil, make([]string, 0), err
	}
	cache.AttachMetadata(examples, false, true)
	return examples, keywords, nil
}

func Close() {
	searcher.Close()
}
