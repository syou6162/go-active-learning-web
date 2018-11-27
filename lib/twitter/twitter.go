package twitter

import (
	"errors"
	"fmt"
	"net/http"
	"sort"

	"strings"

	"time"

	"github.com/codegangsta/cli"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
)

type kv struct {
	Key   string
	Value int
}

func getClient() *twitter.Client {
	consumerKey := util.GetEnv("TWITTER_CONSUMER_KEY", "")
	consumerSecret := util.GetEnv("TWITTER_CONSUMER_SECRET", "")
	accessToken := util.GetEnv("TWITTER_ACCESS_TOKEN", "")
	accessSecret := util.GetEnv("TWITTER_ACCESS_SECRET", "")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}

func GetReferringTweets(url string) ([]string, error) {
	client := getClient()
	search, resp, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: url,
		Count: 100,
	})

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var result []kv
	for _, t := range search.Statuses {
		if t.Retweeted {
			continue
		}
		// twitterのcanonicalがlower caseになっているので、それに合わせる
		url := fmt.Sprintf("https://twitter.com/%s/status/%s", strings.ToLower(t.User.ScreenName), t.IDStr)
		cnt := t.FavoriteCount
		if cnt > 2 {
			result = append(result, kv{url, cnt})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Value > result[j].Value
	})
	tweets := make([]string, 0)
	for _, t := range result {
		tweets = append(tweets, t.Key)
	}
	return tweets, nil
}

func setReferringTweets(app service.GoActiveLearningApp, listName string) error {
	urls, err := app.GetUrlsFromList(listName, 0, 100)
	if err != nil {
		return err
	}

	examples, err := app.SearchExamplesByUlrs(urls)
	if err != nil {
		return err
	}
	app.AttachMetadata(examples)

	for _, e := range examples {
		if e.UpdatedAt.Add(time.Hour * 240).Before(time.Now()) {
			continue
		}
		fmt.Println(e.FinalUrl)
		tweets, err := GetReferringTweets(e.FinalUrl)
		if err != nil {
			u := e.FinalUrl
			if e.FinalUrl == "" {
				u = e.Url
			}
			fmt.Printf("cannot retrieve %s: %s", u, err.Error())
			continue
		}

		for _, t := range tweets {
			tweetExample := example.NewExample(t, model.UNLABELED)
			if err = app.InsertOrUpdateExample(tweetExample); err != nil {
				return err
			}
		}

		tweetExamples, err := app.SearchExamplesByUlrs(tweets)
		if err != nil {
			return err
		}
		app.Fetch(tweetExamples)
		app.UpdateExamplesMetadata(tweetExamples)
		tweetExamples = util.UniqueByFinalUrl(tweetExamples)

		tweets = []string{}
		for _, t := range tweetExamples {
			fmt.Printf("- %s\n", t.FinalUrl)
			tweets = append(tweets, t.FinalUrl)
		}
		tweets = append(tweets, e.ReferringTweets...)
		tweets = util.RemoveDuplicate(tweets)
		e.ReferringTweets = tweets
		app.UpdateExampleMetadata(*e)
	}
	return nil
}

func doSetReferringTweets(c *cli.Context) error {
	listName := c.String("listname")

	app, err := service.NewDefaultApp()
	if err != nil {
		return err
	}
	defer app.Close()
	return setReferringTweets(app, listName)
}

var CommandSetReferringTweets = cli.Command{
	Name:  "set-referring-tweets",
	Usage: "set referring tweets",
	Description: `
Set referring tweets.
`,
	Action: doSetReferringTweets,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "listname", Usage: "List name for cache"},
	},
}
