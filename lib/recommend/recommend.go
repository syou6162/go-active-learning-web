package recommend

import (
	"fmt"
	"log"
	"strings"

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
	"github.com/syou6162/go-active-learning/lib/hatena_bookmark"
	"github.com/syou6162/go-active-learning/lib/util"
)

var listName2Rule = map[string]*regexp.Regexp{
	"general": regexp.MustCompile(`.+`),
	"article": regexp.MustCompile(`.+`),                              // あとでog:typeで絞り込む
	"github":  regexp.MustCompile(`https://github.com/[^/]+/[^/]+$`), // リポジトリのトップのみ
	"slide":   regexp.MustCompile(`https://(www.slideshare.net|speakerdeck.com)/.+`),
	"arxiv":   regexp.MustCompile(`https://arxiv.org/abs/.+`),
}

func doRecommend(c *cli.Context) error {
	subsetSelection := c.Bool("subset-selection")
	sizeConstraint := c.Int("size-constraint")
	alpha := c.Float64("alpha")
	r := c.Float64("r")
	lambda := c.Float64("lambda")
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
	examples = util.FilterStatusCodeOkExamples(examples)
	model := classifier.NewBinaryClassifier(examples)

	targetExamples, err := db.ReadRecentExamples(time.Now().Add(-time.Duration(24*durationDay) * time.Hour))
	if err != nil {
		return err
	}

	targetExamples = util.RemoveNegativeExamples(targetExamples)
	log.Println("Started to attach metadata to positive or unlabeled...")
	cache.AttachMetadata(targetExamples, true, false)
	targetExamples = util.FilterStatusCodeOkExamples(targetExamples)
	targetExamples = util.UniqueByFinalUrl(targetExamples)
	targetExamples = util.UniqueByTitle(targetExamples)
	log.Println(fmt.Sprintf("target size: %d", len(targetExamples)))

	log.Println("Started to predict scores...")
	result := example.Examples{}
	for _, e := range targetExamples {
		if !rule.MatchString(e.FinalUrl) {
			continue
		}
		if listName == "general" && e.IsTwitterUrl() {
			continue
		}
		if listName == "article" && !e.IsArticle() {
			continue
		}
		e.Score = model.PredictScore(e.Fv)
		e.Title = strings.Replace(e.Title, "\n", " ", -1)
		if err := cache.SetExample(*e); err != nil {
			return err
		}
		if e.Score > scoreThreshold {
			hour := 24 * 31 * 6 // 6 months
			cache.SetExampleExpire(*e, time.Hour*time.Duration(hour))
			result = append(result, e)
		}
	}

	log.Println("Started to filter by submodular...")
	log.Println(fmt.Sprintf("Original result size: %d", len(result)))
	if subsetSelection {
		result = submodular.SelectSubExamplesBySubModular(result, sizeConstraint, alpha, r, lambda)
	}

	log.Println("Started to write result...")
	err = cache.AddExamplesToList(listName, result)
	if err != nil {
		return err
	}

	for _, e := range result {
		if bookmarks, err := hatena_bookmark.GetHatenaBookmark(e.FinalUrl); err == nil {
			e.HatenaBookmarks = *bookmarks
			cache.SetExample(*e)
		}
		fmt.Println(fmt.Sprintf("%0.03f\t%s", e.Score, e.Url))
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

var CommandRecommend = cli.Command{
	Name:  "recommend",
	Usage: "Get recommendation list and store them",
	Description: `
Get recommendation list and store them.
`,
	Action: doRecommend,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "subset-selection", Usage: "Use subset selection algorithm (maximizing submodular function) to filter entries"},
		cli.Int64Flag{Name: "size-constraint", Value: 10, Usage: "Budget constraint. Max number of entries to be contained"},
		cli.Float64Flag{Name: "alpha", Value: 1.0},
		cli.Float64Flag{Name: "r", Value: 1.0, Usage: "Scaling factor for number of words"},
		cli.Float64Flag{Name: "lambda", Value: 1.0, Usage: "Diversity parameter"},
		cli.Float64Flag{Name: "score-threshold", Value: 0.0},
		cli.StringFlag{Name: "listname", Usage: "List name for cache"},
		cli.Int64Flag{Name: "duration-day", Usage: "Time span for fetching prediction target", Value: 2},
	},
}
