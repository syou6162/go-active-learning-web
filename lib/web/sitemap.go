package web

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/syou6162/go-active-learning/lib/util"
)

func (s *server) SitemapTop() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sm := stm.NewSitemap(1)
		sm.SetDefaultHost("https://www.machine-learning.news")
		sm.SetCompress(true)
		sm.SetVerbose(true)

		sm.Create()

		sm.Add(stm.URL{{"loc", "/list/general"}, {"changefreq", "daily"}})
		sm.Add(stm.URL{{"loc", "/list/article"}, {"changefreq", "daily"}})
		sm.Add(stm.URL{{"loc", "/list/github"}, {"changefreq", "daily"}})
		sm.Add(stm.URL{{"loc", "/list/arxiv"}, {"changefreq", "daily"}})
		sm.Add(stm.URL{{"loc", "/list/slide"}, {"changefreq", "daily"}})

		sm.Add(stm.URL{{"loc", "/recent-added-examples"}, {"changefreq", "daily"}})

		w.Header().Set("Content-Type", "application/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(sm.XMLContent())
	})
}

func (s *server) SitemapCategory() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		listName := queryValues.Get("category")

		examples, err := s.getUrlsFromList(listName)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintln(w, err.Error())
			return
		}

		sm := stm.NewSitemap(1)
		sm.SetDefaultHost("https://www.machine-learning.news")
		sm.SetCompress(true)
		sm.SetVerbose(true)

		sm.Create()
		for _, e := range examples {
			sm.Add(stm.URL{{"loc", "/example/" + url.PathEscape(e.FinalUrl)}, {"changefreq", "daily"}})
		}

		w.Header().Set("Content-Type", "application/xml; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(sm.XMLContent())
	})
}
