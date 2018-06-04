package web

import (
	"fmt"
	"net/http"

	"io/ioutil"
	"os"
	"strings"

	"html/template"

	"encoding/json"

	"sort"

	"github.com/codegangsta/cli"
	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning-web/lib/assets"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
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

func recentAddedExamples(w http.ResponseWriter, r *http.Request) {
	cache, err := cache.NewCache()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	defer cache.Close()

	conn, err := db.CreateDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	defer conn.Close()

	examples, err := db.ReadLabeledExamples(conn, 100)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	cache.AttachMetaData(examples)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// clear feature vector
	for _, example := range examples {
		example.Fv = make([]string, 0)
	}
	json.NewEncoder(w).Encode(examples)
}

func readAssetTemplate(p string) (string, error) {
	f, err := templates.Assets.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()

	templateByte, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(templateByte), nil
}

func index(w http.ResponseWriter, r *http.Request) {
	var t *template.Template

	indexTemplate, err := readAssetTemplate("/templates/index.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	headTemplate, err := readAssetTemplate("/templates/head.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	footerTemplate, err := readAssetTemplate("/templates/footer.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	t = template.Must(template.New("index").Parse(indexTemplate))
	t = template.Must(t.Parse(headTemplate))
	t = template.Must(t.Parse(footerTemplate))

	err = t.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
}

func getExamplesFromList(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	listName := queryValues.Get("listName")

	cache, err := cache.NewCache()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	defer cache.Close()

	conn, err := db.CreateDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	defer conn.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	getUrlsFromList := func(listName string) (example.Examples, error) {
		generalUrls, err := cache.GetUrlsFromList(listName, 0, 100)
		if err != nil {
			return nil, err
		}
		examples, err := db.SearchExamplesByUlrs(conn, generalUrls)
		if err != nil {
			return nil, err
		}
		cache.AttachMetaData(examples)
		sort.Sort(sort.Reverse(examples))
		return examples, nil
	}

	examples, err := getUrlsFromList(listName)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// clear feature vector
	for _, example := range examples {
		example.Fv = make([]string, 0)
	}

	json.NewEncoder(w).Encode(examples)
}

func doServe(c *cli.Context) error {
	addr := c.String("addr")
	if addr == "" {
		addr = ":7777"
	}

	http.HandleFunc("/", index) // ハンドラを登録してウェブページを表示させる
	http.HandleFunc("/register_training_data", registerTrainingData)
	http.HandleFunc("/api/recent_added_examples", recentAddedExamples)
	http.HandleFunc("/api/examples", getExamplesFromList)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(templates.Assets)))
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
