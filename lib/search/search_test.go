package search_test

import (
	"testing"

	"os"

	"github.com/syou6162/go-active-learning-web/lib/search"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	os.Exit(ret)
}

func TestSearch(t *testing.T) {
	err := db.Init()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	_, err = db.DeleteAllExamples()
	if err != nil {
		t.Error(err.Error())
	}

	err = cache.Init()
	if err != nil {
		t.Error(err.Error())
	}
	defer cache.Close()

	e1 := example.Example{Url: "https://www.yasuhisay.info/entry/2018/10/04/201000", Label: example.POSITIVE}
	e2 := example.Example{Url: "https://www.yasuhisay.info/entry/2018/10/01/090000", Label: example.POSITIVE}
	e3 := example.Example{Url: "https://www.yasuhisay.info/entry/mackerel_meetup_12_anomaly_detection", Label: example.POSITIVE}
	examples := example.Examples{&e1, &e2, &e3}

	for _, e := range examples {
		_, err = db.InsertOrUpdateExample(e)
		if err != nil {
			t.Error(err)
		}
	}
	cache.AttachMetadata(examples, true, true)
	search.Init()
	defer search.Close()

	result, err := search.Search("機械学習")
	if err != nil {
		t.Error(err.Error())
	}
	if len(result) == 0 {
		t.Error("Result must not be empty")
	}
}
