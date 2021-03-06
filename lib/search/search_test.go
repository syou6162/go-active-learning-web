package search_test

import (
	"testing"

	"os"

	"github.com/syou6162/go-active-learning-web/lib/search"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	os.Exit(ret)
}

func TestSearch(t *testing.T) {
	app, err := service.NewDefaultApp()
	if err != nil {
		t.Error(err.Error())
	}
	defer app.Close()

	if err = app.DeleteAllExamples(); err != nil {
		t.Error(err.Error())
	}

	e1 := example.NewExample("https://www.yasuhisay.info/entry/2018/10/04/201000", model.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info/entry/2018/10/01/090000", model.POSITIVE)
	e3 := example.NewExample("https://www.yasuhisay.info/entry/mackerel_meetup_12_anomaly_detection", model.POSITIVE)
	examples := model.Examples{e1, e2, e3}

	app.Fetch(examples)
	for _, e := range examples {
		if err = app.UpdateOrCreateExample(e); err != nil {
			t.Error(err)
		}
	}

	result, err := search.Search(app, "機械 学習")
	if err != nil {
		t.Error(err.Error())
	}
	if len(result) == 0 {
		t.Error("Result must not be empty")
	}
}
