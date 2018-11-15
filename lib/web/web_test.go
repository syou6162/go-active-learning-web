package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"
	"log"
	"os"

	"net/url"
	"strings"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/syou6162/go-active-learning-web/lib/search"
	"github.com/syou6162/go-active-learning-web/lib/web"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func TestMain(m *testing.M) {
	err := db.Init()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	_, err = db.DeleteAllExamples()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = cache.Init()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cache.Close()

	ret := m.Run()
	os.Exit(ret)
}

func TestRecentAddedExamples(t *testing.T) {
	_, err := db.DeleteAllExamples()
	if err != nil {
		t.Error("Cannot delete examples")
	}

	req, err := http.NewRequest("GET", "/api/recent_added_examples", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	http.HandlerFunc(web.RecentAddedExamples).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	examples := example.Examples{}
	json.Unmarshal(w.Body.Bytes(), &examples)

	if len(examples) != 0 {
		t.Errorf("handler returned wrong length of examples: got %v want %v", len(examples), 0)
	}

	inputFilename := "../../tech_input_example.txt"
	train, err := file.ReadExamples(inputFilename)
	if err != nil {
		t.Error(err)
	}
	for _, example := range train {
		_, err = db.InsertOrUpdateExample(example)
		if err != nil {
			t.Error(err)
		}
	}
	cache.AttachMetadata(train, true, false)

	req, err = http.NewRequest("GET", "/api/recent_added_examples", nil)
	if err != nil {
		t.Error(err)
	}
	w = httptest.NewRecorder()
	http.HandlerFunc(web.RecentAddedExamples).ServeHTTP(w, req)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	result := web.RecentAddedExamplesResult{}
	json.Unmarshal(w.Body.Bytes(), &result)

	if len(result.PositiveExamples) == 0 {
		t.Error("Result must not be empty")
	}
}

func TestGetExamplesFromList(t *testing.T) {
	inputFilename := "../../tech_input_example.txt"
	train, err := file.ReadExamples(inputFilename)
	if err != nil {
		t.Error(err)
	}
	for _, example := range train {
		_, err = db.InsertOrUpdateExample(example)
		if err != nil {
			t.Error(err)
		}
	}
	cache.AddExamplesToList("general", train)

	req, err := http.NewRequest("GET", "/api/examples?listName=general", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	http.HandlerFunc(web.GetExamplesFromList).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	examplesFromList := web.ExamplesFromList{}
	json.Unmarshal(w.Body.Bytes(), &examplesFromList)

	if len(examplesFromList.Examples) == 0 {
		t.Error("Result must not be empty")
	}
}

func TestSearch(t *testing.T) {
	inputFilename := "../../tech_input_example.txt"
	train, err := file.ReadExamples(inputFilename)
	if err != nil {
		t.Error(err)
	}

	cache.AttachMetadata(train, true, true)

	for _, example := range train {
		_, err = db.InsertOrUpdateExample(example)
		if err != nil {
			t.Error(err)
		}
	}

	if err = search.Init(); err != nil {
		t.Error(err)
	}
	defer search.Close()

	values := url.Values{}
	values.Set("query", "blog")

	req, err := http.NewRequest("POST", "/api/search", strings.NewReader(values.Encode()))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	http.HandlerFunc(web.Search).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	result := web.SearchResult{}
	json.Unmarshal(w.Body.Bytes(), &result)

	if len(result.Examples) == 0 || result.Count == 0 {
		t.Error("Result must not be empty")
	}
}

func TestServerAvail(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/server_avail", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	http.HandlerFunc(web.ServerAvail).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestStats(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/stats", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	http.HandlerFunc(stats_api.Handler).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
