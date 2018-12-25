package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/syou6162/go-active-learning-web/lib/web"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func TestSitemapTop(t *testing.T) {
	req, err := http.NewRequest("GET", "/sitemap/top", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	app, err := service.NewDefaultApp()
	if err != nil {
		t.Error(err)
	}
	defer app.Close()

	svr := web.NewServer(app)
	http.Handler(svr.SitemapTop()).ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSitemapCategory(t *testing.T) {
	app, err := service.NewDefaultApp()
	if err != nil {
		t.Error(err)
	}
	defer app.Close()

	inputFilename := "../../tech_input_example.txt"
	train, err := file.ReadExamples(inputFilename)
	if err != nil {
		t.Error(err)
	}
	for _, example := range train {
		if err = app.UpdateOrCreateExample(example); err != nil {
			t.Error(err)
		}
	}
	if err := app.UpdateRecommendation("general", train); err != nil {
		t.Error(err)
	}

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
	app, err := service.NewDefaultApp()
	if err != nil {
		t.Error(err)
	}
	defer app.Close()

	inputFilename := "../../tech_input_example.txt"
	train, err := file.ReadExamples(inputFilename)
	if err != nil {
		t.Error(err)
	}

	app.Fetch(train)
	for _, example := range train {
		example.Score = 10.0
		if err = app.UpdateOrCreateExample(example); err != nil {
			t.Error(err)
		}
	}

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
