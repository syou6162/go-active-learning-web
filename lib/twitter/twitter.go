package twitter

import (
	"errors"
	"fmt"
	"net/http"

	"time"

	"github.com/codegangsta/cli"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
)

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

func GetReferringTweets(url string) (model.ReferringTweets, error) {
	client := getClient()
	search, resp, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: fmt.Sprintf("%s -filter:retweets", url),
		Count: 100,
		TweetMode: "extended",
	})

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	tweets := model.ReferringTweets{}
	for _, t := range search.Statuses {
		createdAt, err := t.CreatedAtTime()
		if err != nil {
			createdAt = time.Now()
		}
		tweet := model.Tweet{
			CreatedAt: createdAt,
			IdStr: t.IDStr,
			FullText: t.FullText,
			FavoriteCount: t.FavoriteCount,
			RetweetCount: t.RetweetCount,
			Lang: t.Lang,

			ScreenName: t.User.ScreenName,
			Name: t.User.Name,
			ProfileImageUrl: t.User.ProfileImageURLHttps,
		}
		tweets = append(tweets, &tweet)
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
		e.ReferringTweets = &tweets
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
