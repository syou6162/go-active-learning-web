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

	"encoding/json"

	"sort"

	"syscall"

	"github.com/codegangsta/cli"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/syou6162/go-active-learning-web/lib/search"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util"
)

func checkAuth(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if ok == false {
		return false
	}
	return username == os.Getenv("BASIC_AUTH_USERNAME") && password == os.Getenv("BASIC_AUTH_PASSWORD")
}

func registerTrainingData(w http.ResponseWriter, r *http.Request) {
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
		err = db.InsertExamplesFromReader(strings.NewReader(string(buf)))
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintln(w, err.Error())
			return
		}
	}
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func lightenExamples(examples example.Examples) {
	for _, example := range examples {
		example.Fv = make([]string, 0)
		r := []rune(example.Body)
		example.Body = string(r[0:Min(500, len(r))])
	}
}

func RecentAddedExamples(w http.ResponseWriter, r *http.Request) {
	positiveExamples, err := db.ReadPositiveExamples(30)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	negativeExamples, err := db.ReadNegativeExamples(30)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	unlabeledExamples, err := db.ReadUnlabeledExamples(30)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	var examples example.Examples
	examples = append(examples, positiveExamples...)
	examples = append(examples, negativeExamples...)
	examples = append(examples, unlabeledExamples...)
	cache.AttachMetadata(examples, false, true)
	examples = util.FilterStatusCodeOkExamples(examples)
	lightenExamples(examples)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(examples)
}

func GetExamplesFromList(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	listName := queryValues.Get("listName")

	getUrlsFromList := func(listName string) (example.Examples, error) {
		urls, err := cache.GetUrlsFromList(listName, 0, 100)
		if err != nil {
			return nil, err
		}

		examples, err := db.SearchExamplesByUlrs(urls)
		if err != nil {
			return nil, err
		}

		cache.AttachMetadata(examples, false, true)
		sort.Sort(sort.Reverse(examples))
		result := util.RemoveNegativeExamples(examples)
		return result, nil
	}

	examples, err := getUrlsFromList(listName)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	examples = util.FilterStatusCodeOkExamples(examples)
	lightenExamples(examples)
	json.NewEncoder(w).Encode(examples)
}

func Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := r.FormValue("query")
	examples, err := search.Search(query)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	examples = util.FilterStatusCodeOkExamples(examples)
	lightenExamples(examples)
	json.NewEncoder(w).Encode(examples)
}

func ServerAvail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	if err := cache.Ping(); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	if err := search.Ping(); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK, I'm fine")
}

func doServe(c *cli.Context) error {
	addr := c.String("addr")
	if addr == "" {
		addr = ":7778"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/register_training_data", registerTrainingData)
	mux.HandleFunc("/api/recent_added_examples", RecentAddedExamples)
	mux.HandleFunc("/api/examples", GetExamplesFromList)
	mux.HandleFunc("/api/search", Search)
	mux.HandleFunc("/api/server_avail", ServerAvail)
	mux.HandleFunc("/api/stats", stats_api.Handler)

	srv := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	err := db.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	err = cache.Init()
	if err != nil {
		return err
	}
	defer cache.Close()

	search.Init()
	defer search.Close()

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
