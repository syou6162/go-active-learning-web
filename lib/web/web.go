package web

import (
	"fmt"
	"net/http"

	"io/ioutil"
	"os"
	"strings"

	"encoding/json"

	"sort"

	"github.com/codegangsta/cli"
	_ "github.com/lib/pq"
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
		buf, _ := ioutil.ReadAll(r.Body)
		err := db.InsertExamplesFromReader(strings.NewReader(string(buf)))
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintln(w, err.Error())
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

func recentAddedExamples(w http.ResponseWriter, r *http.Request) {
	positiveExamples, err := db.ReadPositiveExamples(30)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	cache.AttachMetadata(positiveExamples, false, true)

	negativeExamples, err := db.ReadNegativeExamples(30)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	cache.AttachMetadata(negativeExamples, false, true)

	unlabeledExamples, err := db.ReadUnlabeledExamples(30)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	cache.AttachMetadata(unlabeledExamples, false, true)
	unlabeledExamples = util.FilterStatusCodeOkExamples(unlabeledExamples)

	var examples example.Examples
	examples = append(examples, positiveExamples...)
	examples = append(examples, negativeExamples...)
	examples = append(examples, unlabeledExamples...)
	lightenExamples(examples)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(examples)
}

func getExamplesFromList(w http.ResponseWriter, r *http.Request) {
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

	lightenExamples(examples)
	json.NewEncoder(w).Encode(examples)
}

func doServe(c *cli.Context) error {
	addr := c.String("addr")
	if addr == "" {
		addr = ":7778"
	}

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

	http.HandleFunc("/api/register_training_data", registerTrainingData)
	http.HandleFunc("/api/recent_added_examples", recentAddedExamples)
	http.HandleFunc("/api/examples", getExamplesFromList)
	return http.ListenAndServe(addr, nil)
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
