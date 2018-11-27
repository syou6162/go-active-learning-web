package search_test

import (
	"testing"

	"os"

	"github.com/syou6162/go-active-learning-web/lib/search"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/repository"
	"github.com/syou6162/go-active-learning/lib/service"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	os.Exit(ret)
}

func TestSearch(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Error(err.Error())
	}
	app := service.NewApp(repo)
	defer app.Close()

	if err = app.DeleteAllExamples(); err != nil {
		t.Error(err.Error())
	}

	err = cache.Init()
	if err != nil {
		t.Error(err.Error())
	}
	defer cache.Close()

	e1 := model.Example{Url: "https://www.yasuhisay.info/entry/2018/10/04/201000", Label: model.POSITIVE}
	e2 := model.Example{Url: "https://www.yasuhisay.info/entry/2018/10/01/090000", Label: model.POSITIVE}
	e3 := model.Example{Url: "https://www.yasuhisay.info/entry/mackerel_meetup_12_anomaly_detection", Label: model.POSITIVE}
	examples := model.Examples{&e1, &e2, &e3}

	for _, e := range examples {
		if err = app.InsertOrUpdateExample(e); err != nil {
			t.Error(err)
		}
	}
	cache.AttachMetadata(examples, true, true)
	search.Init(app)
	defer search.Close()

	result, err := search.Search(app, "機械 学習")
	if err != nil {
		t.Error(err.Error())
	}
	if len(result) == 0 {
		t.Error("Result must not be empty")
	}
}
