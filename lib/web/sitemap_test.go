package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/syou6162/go-active-learning-web/lib/web"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func TestSitemapTop(t *testing.T) {
	req, err := http.NewRequest("GET", "/sitemap/top", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	http.HandlerFunc(web.SitemapTop).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSitemapCategory(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/sitemap?category=general", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	http.HandlerFunc(web.SitemapCategory).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSitemapRecentPositiveExamples(t *testing.T) {
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
	cache.AttachMetadata(train, true, true)

	req, err := http.NewRequest("GET", "/sitemap/recent", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	http.HandlerFunc(web.SitemapRecentPositiveExamples).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
