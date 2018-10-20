package search

import (
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/db"
	"hash/fnv"
	"sync"
	"github.com/syou6162/go-active-learning/lib/example"
)

var (
	once sync.Once
	searcher = riot.Engine{}
	id2url = map[uint64]string{}
)

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func Init() error {
	var err error
	once.Do(func() {
		searcher.Init(types.EngineOpts{
			GseDict: "jp",
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
	return nil
}

func Search(query string) (example.Examples, error) {
	urls := make([]string, 0)
	for _, resp := range searcher.SearchDoc(types.SearchReq{
		Text: query,
		RankOpts: &types.RankOpts{MaxOutputs: 100},
		}).Docs {
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
