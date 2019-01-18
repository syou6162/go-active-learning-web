package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"time"

	"io/ioutil"
	"os"
	"strings"

	"sort"

	"syscall"

	"strconv"

	"github.com/codegangsta/cli"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/gorilla/feeds"
	"github.com/syou6162/go-active-learning-web/lib/ahocorasick"
	"github.com/syou6162/go-active-learning-web/lib/search"
	"github.com/syou6162/go-active-learning-web/lib/version"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
)

type Server interface {
	Handler() http.Handler
	SitemapTop() http.Handler
	SitemapCategory() http.Handler

	RecentAddedExamples() http.Handler
	GetExamplesFromList() http.Handler
	GetExampleById() http.Handler
	Search() http.Handler
	GetFeed() http.Handler
	ServerAvail() http.Handler
}

func NewServer(app service.GoActiveLearningApp) Server {
	return &server{app: app}
}

type server struct {
	app service.GoActiveLearningApp
}

func checkAuth(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if ok == false {
		return false
	}
	return username == os.Getenv("BASIC_AUTH_USERNAME") && password == os.Getenv("BASIC_AUTH_PASSWORD")
}

func (s *server) registerTrainingData() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if checkAuth(r) == false {
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			return
		} else {
			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				fmt.Fprintln(w, err.Error())
				return
			}
			defer r.Body.Close()
			err = s.app.InsertExamplesFromReader(strings.NewReader(string(buf)))
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				fmt.Fprintln(w, err.Error())
				return
			}
		}
	})
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func lightenExamples(examples model.Examples) {
	for _, example := range examples {
		example.Fv = make([]string, 0)
		r := []rune(example.Body)
		example.Body = string(r[0:Min(500, len(r))])
	}
}

type RecentAddedExamplesResult struct {
	PositiveExamples  model.Examples
	NegativeExamples  model.Examples
	UnlabeledExamples model.Examples
}

func (s *server) RecentAddedExamples() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		positiveExamples, err := s.app.SearchPositiveExamples(30)
		if err != nil {
			BadRequest(w, err.Error())
			fmt.Fprintln(w, err.Error())
			return
		}
		s.app.AttachLightMetadata(positiveExamples)

		negativeExamples, err := s.app.SearchNegativeExamples(30)
		if err != nil {
			BadRequest(w, err.Error())
			fmt.Fprintln(w, err.Error())
			return
		}
		s.app.AttachLightMetadata(negativeExamples)

		unlabeledExamples, err := s.app.SearchUnlabeledExamples(30)
		if err != nil {
			BadRequest(w, err.Error())
			fmt.Fprintln(w, err.Error())
			return
		}
		s.app.AttachLightMetadata(unlabeledExamples)
		unlabeledExamples = util.FilterStatusCodeOkExamples(unlabeledExamples)

		JSON(w, http.StatusOK, RecentAddedExamplesResult{
			PositiveExamples:  positiveExamples,
			NegativeExamples:  negativeExamples,
			UnlabeledExamples: unlabeledExamples,
		})
	})
}

func (s *server) getUrlsFromList(listName string) (model.Examples, error) {
	examples, err := s.app.GetRecommendation(listName)
	if err != nil {
		return nil, err
	}
	s.app.AttachLightMetadata(examples)
	sort.Sort(sort.Reverse(examples))
	result := util.RemoveNegativeExamples(examples)
	return result, nil
}

type ExamplesFromList struct {
	Examples model.Examples
}

func (s *server) GetExamplesFromList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		listName := queryValues.Get("listName")

		examples, err := s.getUrlsFromList(listName)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintln(w, err.Error())
			return
		}

		examples = util.FilterStatusCodeOkExamples(examples)
		lightenExamples(examples)
		JSON(w, http.StatusOK, ExamplesFromList{
			Examples: examples,
		})
	})
}

type ExampleWithSimilarExamples struct {
	Example         *model.Example
	SimilarExamples model.Examples `json:"SimilarExamples"`
	Keywords        []string
}

func (s *server) GetExampleById() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		q := queryValues.Get("id")
		id, err := strconv.Atoi(q)
		if err != nil {
			BadRequest(w, fmt.Sprintf("Cannot parse id: %s", q))
			return
		}

		ex, err := s.app.FindExampleById(id)
		if err != nil {
			NotFound(w, fmt.Sprintf("No such id: %d", id))
			return
		}

		s.app.AttachLightMetadata(model.Examples{ex})
		if err != nil {
			BadRequest(w, err.Error())
			fmt.Fprintln(w, err.Error())
			return
		}

		//similarExamples, keywords, err := search.SearchSimilarExamples(s.app, ex.Title, 5)
		//if err != nil {
		//	BadRequest(w, err.Error())
		//	fmt.Fprintln(w, err.Error())
		//	return
		//}
		similarExamplesWithoutOriginal := model.Examples{}
		//for _, e := range similarExamples {
		//	if e.FinalUrl != ex.FinalUrl {
		//		similarExamplesWithoutOriginal = append(similarExamplesWithoutOriginal, e)
		//	}
		//}
		//similarExamplesWithoutOriginal = util.FilterStatusCodeOkExamples(similarExamplesWithoutOriginal)

		//w.Header().Set("X-Keywords", strings.Join(keywords, ","))
		JSON(w, http.StatusOK, ExampleWithSimilarExamples{
			Example:         ex,
			SimilarExamples: similarExamplesWithoutOriginal,
			Keywords:        ahocorasick.SearchKeywords(strings.ToLower(ex.Title)),
		})
	})
}

type SearchResult struct {
	Examples model.Examples
	Query    string `json:"Query"`
	Count    int    `json:"Count"`
}

func (s *server) Search() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			BadRequest(w, "Only POST method is supported")
			return
		}

		query := r.FormValue("query")

		examples, err := search.Search(s.app, query)
		if err != nil {
			BadRequest(w, err.Error())
			fmt.Fprintln(w, err.Error())
			return
		}

		examples = util.FilterStatusCodeOkExamples(examples)
		lightenExamples(examples)

		JSON(w, http.StatusOK, SearchResult{
			Examples: examples,
			Query:    query,
			Count:    len(examples),
		})
	})
}

func (s *server) GetFeed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		listName := queryValues.Get("listName")

		examples, err := s.getUrlsFromList(listName)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintln(w, err.Error())
			return
		}

		examples = util.FilterStatusCodeOkExamples(examples)
		lightenExamples(examples)

		now := time.Now()
		feed := &feeds.Feed{
			Title:       fmt.Sprintf("ML-News - %s", listName),
			Link:        &feeds.Link{Href: fmt.Sprintf("https://www.machine-learning.news/list/%s", listName)},
			Description: "機械学習に関連する人気のエントリを読むことができます",
			Author:      &feeds.Author{Name: "Yasuhisa Yoshida"},
			Created:     now,
		}

		for _, e := range examples {
			item := &feeds.Item{
				Title:       e.Title,
				Link:        &feeds.Link{Href: fmt.Sprintf("https://www.machine-learning.news/example/%d", e.Id)},
				Description: e.Description,
				Created:     e.CreatedAt,
			}
			feed.Items = append(feed.Items, item)
		}

		rss, err := feed.ToRss()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintln(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rss))
	})
}

func (s *server) ServerAvail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		if err := s.app.Ping(); err != nil {
			UnavaliableError(w, err.Error())
			fmt.Fprintln(w, err.Error())
			return
		}
		Ok(w, "OK, I'm fine")
	})
}

func (s *server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "text/plain")
		resp.Header().Set("X-Revision", version.GitCommit)
		fmt.Fprintln(resp, "I'm ML-News.")
	})

	mux.Handle("/api/register_training_data", s.registerTrainingData())
	mux.Handle("/api/recent_added_examples", s.RecentAddedExamples())
	mux.Handle("/api/examples", s.GetExamplesFromList())
	mux.Handle("/api/example", s.GetExampleById())
	mux.Handle("/api/search", s.Search())
	mux.Handle("/api/server_avail", s.ServerAvail())
	mux.HandleFunc("/api/stats", stats_api.Handler)
	mux.Handle("/sitemap", s.SitemapCategory())
	mux.Handle("/sitemap/top", s.SitemapTop())
	mux.Handle("/rss", s.GetFeed())
	return mux
}

func doServe(c *cli.Context) error {
	addr := c.String("addr")
	if addr == "" {
		addr = ":7778"
	}

	app, err := service.NewDefaultApp()
	if err != nil {
		return err
	}
	defer app.Close()

	srv := &http.Server{
		Addr:    addr,
		Handler: NewServer(app).Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	ahocorasick.Init()

	// SIGINTとSYSTERMが飛んできたらgraceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
	return nil
}

var CommandServe = cli.Command{
	Name:  "serve",
	Usage: "Run a server",
	Description: `
Run a web server.
`,
	Action: doServe,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "addr"},
	},
}
