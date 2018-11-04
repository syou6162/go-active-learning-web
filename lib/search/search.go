package search

import (
	"errors"
	"hash/fnv"
	"sync"

	"os"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util"
)

var (
	once         sync.Once
	riotDictPath string
	searcher     = riot.Engine{}
	id2url       = map[uint64]string{}
	avail        = false
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

func Close() {
	searcher.Close()
}
