package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/syou6162/go-active-learning-web/lib/web"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/repository"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func TestSitemapTop(t *testing.T) {
	req, err := http.NewRequest("GET", "/sitemap/top", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	repo, err := repository.New()
	if err != nil {
		t.Error(err)
	}
	app := service.NewApp(repo)
	defer app.Close()

	svr := web.NewServer(app)
	http.Handler(svr.SitemapTop()).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSitemapCategory(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Error(err)
	}
	app := service.NewApp(repo)
	defer app.Close()

	inputFilename := "../../tech_input_example.txt"
	train, err := file.ReadExamples(inputFilename)
	if err != nil {
		t.Error(err)
	}
	for _, example := range train {
		println(example.Url)
		if err = app.InsertOrUpdateExample(example); err != nil {
			t.Error(err)
		}
	}
	cache.AddExamplesToList("general", train)

	req, err := http.NewRequest("GET", "/sitemap?category=general", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	svr := web.NewServer(app)

	http.Handler(svr.SitemapCategory()).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSitemapRecentPositiveExamples(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Error(err)
	}
	app := service.NewApp(repo)
	defer app.Close()

	inputFilename := "../../tech_input_example.txt"
	train, err := file.ReadExamples(inputFilename)
	if err != nil {
		t.Error(err)
	}
	for _, example := range train {
		if err = app.InsertOrUpdateExample(example); err != nil {
			t.Error(err)
		}
	}
	cache.AttachMetadata(train, true, true)

	req, err := http.NewRequest("GET", "/sitemap/recent", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	svr := web.NewServer(app)
	http.Handler(svr.SitemapRecentPositiveExamples()).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
