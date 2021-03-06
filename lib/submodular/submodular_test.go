package submodular

import (
	"fmt"
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
)

func TestGetDF(t *testing.T) {
	e := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e.Title = "こんにちは、日本"
	dfMap := GetDF(*e)

	japan := "日本"
	if _, ok := dfMap[japan]; !ok {
		t.Error(fmt.Printf("Example must contain %s", japan))
	}
}

func TestGetIDF(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e1.Title = "こんにちは、日本"
	idfMap := GetIDF(model.Examples{e1})

	japan := "日本"
	if _, ok := idfMap[japan]; !ok {
		t.Error(fmt.Printf("Example must contain %s", japan))
	}
}

func TestSelectSubExamplesBySubModular(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e1.Title = "こんにちは、日本"
	e2 := example.NewExample("http://google.com", model.POSITIVE)
	e2.Title = "hello google"

	examples := SelectSubExamplesBySubModular(model.Examples{e1, e2}, 1, 1.0, 1.0, 1.0)

	if len(examples) != 1 {
		t.Error(fmt.Printf("Number of selected examples must be %d", len(examples)))
	}
}
