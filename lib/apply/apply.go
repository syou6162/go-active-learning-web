package apply

import (
	"fmt"
	"strings"

	"encoding/json"

	"time"

	"regexp"

	"os"

	"github.com/codegangsta/cli"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/syou6162/go-active-learning-web/lib/submodular"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util"
)

var listName2Rule = map[string]*regexp.Regexp{
	"general": regexp.MustCompile(`.+`),
	"article": regexp.MustCompile(`.+`), // あとでog:typeで絞り込む
	"github":  regexp.MustCompile(`https://github.com/.+`),
	"slide":   regexp.MustCompile(`https://(www.slideshare.net|speakerdeck.com)/.+`),
	"twitter": regexp.MustCompile(`https://twitter.com/.+`),
	"arxiv":   regexp.MustCompile(`https://arxiv.org/abs/.+`),
}

func doApply(c *cli.Context) error {
	filterStatusCodeOk := c.Bool("filter-status-code-ok")
	jsonOutput := c.Bool("json-output")
	subsetSelection := c.Bool("subset-selection")
	sizeConstraint := c.Int("size-constraint")
	alpha := c.Float64("alpha")
	r := c.Float64("r")
	scoreThreshold := c.Float64("score-threshold")
	durationDay := c.Int64("duration-day")
	listName := c.String("listname")
	rule, ok := listName2Rule[listName]
	if ok == false {
		return cli.NewExitError("No matched rule", 1)
	}

	err := cache.Init()
	if err != nil {
		return err
	}
	defer cache.Close()

	err = db.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	examples, err := db.ReadLabeledExamples(100000)
	if err != nil {
		return err
	}
	err = postNumOfPositiveAndNegativeExamplesToMackerel(examples)
	if err != nil {
		return err
	}

	cache.AttachMetadata(examples, true, false)
	if filterStatusCodeOk {
		examples = util.FilterStatusCodeOkExamples(examples)
	}
	model := classifier.NewBinaryClassifier(examples)

	targetExamples, err := db.ReadRecentExamples(time.Now().Add(-time.Duration(24*durationDay) * time.Hour))
	if err != nil {
		return err
	}

	targetExamples = util.RemoveNegativeExamples(targetExamples)
	cache.AttachMetadata(targetExamples, true, false)
	if filterStatusCodeOk {
		targetExamples = util.FilterStatusCodeOkExamples(targetExamples)
	}

	result := example.Examples{}
	for _, e := range targetExamples {
		if !rule.MatchString(e.FinalUrl) {
			continue
		}
		if listName == "article" && !e.IsArticle() {
			continue
		}
		e.Score = model.PredictScore(e.Fv)
		e.Title = strings.Replace(e.Title, "\n", " ", -1)
		if e.Score > scoreThreshold {
			result = append(result, e)
		}
	}

	if subsetSelection {
		result = submodular.SelectSubExamplesBySubModular(result, sizeConstraint, alpha, r)
	}

	err = cache.AddExamplesToList(listName, result)
	if err != nil {
		return err
	}

	for _, e := range result {
		if err := cache.SetExample(*e); err != nil {
			return err
		}
		if jsonOutput {
			b, err := json.Marshal(e)
			if err != nil {
				return err
			}
			fmt.Println(string(b))
		} else {
			fmt.Println(fmt.Sprintf("%0.03f\t%s", e.Score, e.Url))
		}
	}

	return nil
}

func postNumOfPositiveAndNegativeExamplesToMackerel(examples example.Examples) error {
	apiKey := os.Getenv("MACKEREL_API_KEY")
	serviceName := os.Getenv("MACKEREL_SERVICE_NAME")
	if apiKey == "" || serviceName == "" {
		return nil
	}

	numPos := 0
	numNeg := 0
	for _, e := range examples {
		if e.Label == example.POSITIVE {
			numPos++
		} else if e.Label == example.NEGATIVE {
			numNeg++
		}
	}

	client := mkr.NewClient(apiKey)
	now := time.Now().Unix()
	err := client.PostServiceMetricValues(serviceName, []*mkr.MetricValue{
		{
			Name:  "count.positive",
			Time:  now,
			Value: numPos,
		},
		{
			Name:  "count.negative",
			Time:  now,
			Value: numNeg,
		},
	})
	return err
}

var CommandApply = cli.Command{
	Name:  "apply",
	Usage: "Apply classifier to unlabeled examples",
	Description: `
Apply classifier to unlabeled examples, and print a pair of score and url.
`,
	Action: doApply,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "filter-status-code-ok", Usage: "Use only examples with status code = 200"},
		cli.BoolFlag{Name: "json-output", Usage: "Make output with json format or not (tsv format)."},
		cli.BoolFlag{Name: "subset-selection", Usage: "Use subset selection algorithm (maximizing submodular function) to filter entries"},
		cli.Int64Flag{Name: "size-constraint", Value: 10, Usage: "Budget constraint. Max number of entries to be contained"},
		cli.Float64Flag{Name: "alpha", Value: 1.0},
		cli.Float64Flag{Name: "r", Value: 1.0, Usage: "Scaling factor for number of words"},
		cli.Float64Flag{Name: "score-threshold", Value: 0.0},
		cli.StringFlag{Name: "listname", Usage: "List name for cache"},
		cli.Int64Flag{Name: "duration-day", Usage: "Time span for fetching prediction target", Value: 2},
	},
}
