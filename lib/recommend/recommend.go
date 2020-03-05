package recommend

import (
	"fmt"
	"log"

	"time"

	"regexp"

	"net/url"
	"sort"

	"github.com/syou6162/go-active-learning-web/lib/submodular"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/hatena_bookmark"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/urfave/cli"
)

var listName2Rule = map[string]*regexp.Regexp{
	"general": regexp.MustCompile(`.+`),
	"article": regexp.MustCompile(`.+`),                              // あとでog:typeで絞り込む
	"github":  regexp.MustCompile(`https://github.com/[^/]+/[^/]+$`), // リポジトリのトップのみ
	"slide":   regexp.MustCompile(`https://(www.slideshare.net|speakerdeck.com|docs.google.com/presentation/d)/.+$`),
	"arxiv":   regexp.MustCompile(`https://(arxiv.org/abs/.+$|openreview.net/forum\?id=.+$)`),
	"video":   regexp.MustCompile(`https?://(www.youtube.com/watch\?v=.+$|videolectures.net/.+$|vimeo.com/.+$)`),
	"event":   regexp.MustCompile(`https://(.*?\.?connpass.com/event/.+/$|techplay.jp/event/.+)$`),
}

var listName2Hosts = map[string][]string{
	"general": {"http"},
	"article": {"http"},
	"github":  {"https://github.com"},
	"slide":   {"https://www.slideshare.net", "https://speakerdeck.com"},
	"arxiv":   {"https://arxiv.org", "https://openreview.net"},
	"video":   {"https://www.youtube.com", "http://videolectures.net", "https://vimeo.com"},
	"event":   {"https://"},
}

func UniqByHost(examples model.Examples) model.Examples {
	result := model.Examples{}

	examplesByHost := map[string]model.Examples{}
	for _, e := range examples {
		if u, err := url.Parse(e.FinalUrl); err == nil {
			examplesByHost[u.Host] = append(examplesByHost[u.Host], e)
		}
	}
	for _, arry := range examplesByHost {
		sort.Sort(sort.Reverse(arry))
		result = append(result, arry[0])
	}
	return result
}

func doRecommend(c *cli.Context) error {
	subsetSelection := c.Bool("subset-selection")
	uniqueByHost := c.Bool("unique-by-host")
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

	app, err := service.NewDefaultApp()
	if err != nil {
		return err
	}
	defer app.Close()

	targetExamples := model.Examples{}
	hosts, ok := listName2Hosts[listName]
	if ok == false {
		return cli.NewExitError("No matched rule", 1)
	}

	for _, h := range hosts {
		tmp, err := app.SearchRecentExamplesByHost(h, time.Now().Add(-time.Duration(24*durationDay)*time.Hour), 10000)
		if err != nil {
			return err
		}
		targetExamples = append(targetExamples, tmp...)
	}

	targetExamples = util.RemoveNegativeExamples(targetExamples)
	log.Println("Started to attach metadata to positive or unlabeled...")
	if err = app.AttachMetadata(targetExamples, 0, 0); err != nil {
		return err
	}

	okExamples := util.FilterStatusCodeOkExamples(targetExamples)
	if err = app.AttachMetadataIncludingFeatureVector(okExamples, 0, 0); err != nil {
		return err
	}

	notOkExamples := util.FilterStatusCodeNotOkExamples(targetExamples)
	app.Fetch(notOkExamples)
	for _, e := range util.FilterStatusCodeOkExamples(notOkExamples) {
		app.UpdateOrCreateExample(e)
		app.UpdateFeatureVector(e)
	}

	targetExamples = util.FilterStatusCodeOkExamples(targetExamples)
	targetExamples = util.UniqueByFinalUrl(targetExamples)
	targetExamples = util.UniqueByTitle(targetExamples)
	log.Println(fmt.Sprintf("target size: %d", len(targetExamples)))

	m, err := app.FindLatestMIRAModel(classifier.EXAMPLE)
	if err != nil {
		return err
	}

	log.Println("Started to predict scores...")
	result := model.Examples{}
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
		e.Score = m.PredictScore(e.Fv)
		if e.Score > scoreThreshold {
			result = append(result, e)
		}
	}

	log.Println(fmt.Sprintf("Original result size: %d", len(result)))

	if uniqueByHost {
		result = UniqByHost(result)
		log.Println(fmt.Sprintf("Filtered by host: %d", len(result)))
	}

	log.Println("Started to filter by submodular...")
	if subsetSelection {
		result = submodular.SelectSubExamplesBySubModular(result, sizeConstraint, alpha, r, lambda)
	}

	log.Println("Started to write result...")
	err = app.UpdateRecommendation(listName, result)
	if err != nil {
		return err
	}

	for _, e := range result {
		if bookmark, err := hatena_bookmark.GetHatenaBookmark(e.FinalUrl); err == nil {
			if e.HatenaBookmark.Count > bookmark.Count {
				continue
			}
			e.HatenaBookmark = bookmark
			if err = app.UpdateHatenaBookmark(e); err != nil {
				log.Println(fmt.Sprintf("Error to update bookmark info %s %s", e.Url, err.Error()))
			}
		}
		fmt.Println(fmt.Sprintf("%0.03f\t%s", e.Score, e.Url))
	}

	return nil
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
		cli.BoolFlag{Name: "unique-by-host", Usage: "Filter entries to be unique by host"},
		cli.Int64Flag{Name: "size-constraint", Value: 10, Usage: "Budget constraint. Max number of entries to be contained"},
		cli.Float64Flag{Name: "alpha", Value: 1.0},
		cli.Float64Flag{Name: "r", Value: 1.0, Usage: "Scaling factor for number of words"},
		cli.Float64Flag{Name: "lambda", Value: 1.0, Usage: "Diversity parameter"},
		cli.Float64Flag{Name: "score-threshold", Value: 0.0},
		cli.StringFlag{Name: "listname", Usage: "List name for cache"},
		cli.Int64Flag{Name: "duration-day", Usage: "Time span for fetching prediction target", Value: 2},
	},
}
