package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"time"

	"os"

	"sort"

	"syscall"

	"strconv"

	"github.com/getsentry/sentry-go"

	stats_api "github.com/fukata/golang-stats-api-handler"
	"github.com/gorilla/feeds"
	"github.com/syou6162/go-active-learning-web/lib/ahocorasick"
	"github.com/syou6162/go-active-learning-web/lib/search"
	web_util "github.com/syou6162/go-active-learning-web/lib/util"
	"github.com/syou6162/go-active-learning-web/lib/version"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/urfave/cli"
)

type Server interface {
	Handler() http.Handler
	SitemapTop() http.Handler
	SitemapCategory() http.Handler

	RecentAddedExamples() http.Handler
	RecentAddedReferringTweets() http.Handler
	GetExamplesFromList() http.Handler
	GetExampleById() http.Handler
	GetTweets() http.Handler
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
		s.app.AttachMetadata(positiveExamples, 0, 0)

		negativeExamples, err := s.app.SearchNegativeExamples(30)
		if err != nil {
			BadRequest(w, err.Error())
			fmt.Fprintln(w, err.Error())
			return
		}
		s.app.AttachMetadata(negativeExamples, 0, 0)

		unlabeledExamples, err := s.app.SearchUnlabeledExamples(30)
		if err != nil {
			BadRequest(w, err.Error())
			fmt.Fprintln(w, err.Error())
			return
		}
		s.app.AttachMetadata(unlabeledExamples, 0, 0)
		unlabeledExamples = util.FilterStatusCodeOkExamples(unlabeledExamples)

		JSON(w, http.StatusOK, RecentAddedExamplesResult{
			PositiveExamples:  positiveExamples,
			NegativeExamples:  negativeExamples,
			UnlabeledExamples: unlabeledExamples,
		})
	})
}

func (s *server) getListOfExampleWithTweet(tweets model.ReferringTweets) (model.Examples, error) {
	exampleIds := make([]int, 0)
	for _, t := range tweets.Tweets {
		exampleIds = append(exampleIds, t.ExampleId)
	}
	examples, err := s.app.SearchExamplesByIds(exampleIds)
	if err != nil {
		return nil, err
	}

	tweetsByExampleId := make(map[int][]*model.Tweet)
	for _, t := range tweets.Tweets {
		tweetsByExampleId[t.ExampleId] = append(tweetsByExampleId[t.ExampleId], t)
	}
	for _, e := range examples {
		tmp := tweetsByExampleId[e.Id]
		e.ReferringTweets = &model.ReferringTweets{
			Tweets: tmp,
			Count:  len(tmp),
		}
	}
	return examples, nil
}

func (s *server) RecentAddedReferringTweets() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit := 30
		positive_finished := make(chan bool)
		negative_finished := make(chan bool)
		unlabeled_finished := make(chan bool)
		positiveExamples := model.Examples{}
		negativeExamples := model.Examples{}
		unlabeledExamples := model.Examples{}
		scoreThreshold := -1.0
		tweetsLimitInSameExample := 3

		go func() {
			positiveTweets, err := s.app.SearchPositiveReferringTweets(scoreThreshold, tweetsLimitInSameExample, limit)
			if err != nil {
				ServerError(w, err.Error())
				return
			}
			positiveExamples, err = s.getListOfExampleWithTweet(positiveTweets)
			if err != nil {
				ServerError(w, err.Error())
				return
			}
			positive_finished <- true
		}()

		go func() {
			negativeTweets, err := s.app.SearchNegativeReferringTweets(scoreThreshold, tweetsLimitInSameExample, limit)
			if err != nil {
				ServerError(w, err.Error())
				return
			}
			negativeExamples, err = s.getListOfExampleWithTweet(negativeTweets)
			if err != nil {
				ServerError(w, err.Error())
				return
			}
			negative_finished <- true
		}()

		go func() {
			unlabeledTweets, err := s.app.SearchUnlabeledReferringTweets(scoreThreshold, tweetsLimitInSameExample, 50)
			if err != nil {
				ServerError(w, err.Error())
				return
			}
			unlabeledExamples, err = s.getListOfExampleWithTweet(unlabeledTweets)
			if err != nil {
				ServerError(w, err.Error())
				return
			}
			unlabeled_finished <- true
		}()

		<-positive_finished
		<-negative_finished
		<-unlabeled_finished

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
	s.app.AttachMetadata(examples, 0, 0)
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
		web_util.LightenExamples(examples)
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

func removeByFinalUrl(examples model.Examples, ex model.Example) model.Examples {
	result := model.Examples{}
	for _, e := range examples {
		if e.FinalUrl != ex.FinalUrl {
			result = append(result, e)
		}
	}
	return result
}

func bodyLikeStr(e model.Example) string {
	str := ""
	if e.OgDescription != "" {
		str = e.OgDescription
	} else if e.Description != "" {
		str = e.Description
	} else if e.Body != "" {
		str = e.Body
	}
	r := []rune(str)
	return string(r[0:web_util.Min(200, len(r))])
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

		s.app.AttachMetadata(model.Examples{ex}, 10, 10)
		if err != nil {
			BadRequest(w, err.Error())
			fmt.Fprintln(w, err.Error())
			return
		}
		web_util.LightenExamples(model.Examples{ex})

		similarExamples, keywords, err := search.SearchSimilarExamples(s.app, ex, 5)
		if err != nil {
			BadRequest(w, err.Error())
			return
		}
		similarExamplesWithoutOriginal := util.FilterStatusCodeOkExamples(removeByFinalUrl(similarExamples, *ex))

		JSON(w, http.StatusOK, ExampleWithSimilarExamples{
			Example:         ex,
			SimilarExamples: similarExamplesWithoutOriginal,
			Keywords:        keywords,
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
		web_util.LightenExamples(examples)

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
		web_util.LightenExamples(examples)

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

type ExamplesWithTweets struct {
	Examples model.Examples
}

func (s *server) GetTweets() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		limit := 50
		tweets, err := s.app.SearchRecentReferringTweetsWithHighScore(now.Add(time.Duration(-2*24)*time.Hour), 0.5, limit)
		if err != nil {
			ServerError(w, err.Error())
			return
		}
		examples, err := s.getListOfExampleWithTweet(tweets)
		if err != nil {
			ServerError(w, err.Error())
			return
		}
		JSON(w, http.StatusOK, ExamplesWithTweets{
			Examples: examples,
		})
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

	mux.Handle("/api/recent_added_examples", s.RecentAddedExamples())
	mux.Handle("/api/recent_added_tweets", s.RecentAddedReferringTweets())
	mux.Handle("/api/examples", s.GetExamplesFromList())
	mux.Handle("/api/example", s.GetExampleById())
	mux.Handle("/api/tweets", s.GetTweets())
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
	ahocorasick.Init()

	sentryDSN := os.Getenv("SENTRY_DSN")
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: sentryDSN,
	}); err != nil {
		log.Println(err.Error())
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err.Error())
			sentry.CaptureException(err)
			sentry.Flush(time.Second * 5)
		}
	}()

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
