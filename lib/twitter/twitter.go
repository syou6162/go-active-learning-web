package twitter

import (
	"errors"
	"fmt"
	"net/http"

	"time"

	"bufio"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/syou6162/go-active-learning/lib/classifier"
	tweet_feature "github.com/syou6162/go-active-learning/lib/feature/tweet"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/urfave/cli"
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

func GetScreenNameList(path string) ([]string, error) {
	screenNameList := make([]string, 0)

	file, err := os.Open(path)
	if err != nil {
		return screenNameList, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		screenNameList = append(screenNameList, scanner.Text())
	}
	return screenNameList, nil
}

func GetReferringTweets(url string, blacklist []string) (*model.ReferringTweets, error) {
	client := getClient()
	search, resp, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query:     fmt.Sprintf("%s -filter:retweets", url),
		Count:     100,
		TweetMode: "extended",
	})

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	tweets := make([]*model.Tweet, 0)
	for _, t := range search.Statuses {
		belongsToBlacklist := false
		for _, name := range blacklist {
			if name == t.User.ScreenName {
				belongsToBlacklist = true
				break
			}
		}

		if !belongsToBlacklist {
			createdAt, err := t.CreatedAtTime()
			if err != nil {
				createdAt = time.Now()
			}
			tweet := model.Tweet{
				CreatedAt:     createdAt,
				IdStr:         t.IDStr,
				FullText:      t.FullText,
				FavoriteCount: t.FavoriteCount,
				RetweetCount:  t.RetweetCount,
				Lang:          t.Lang,

				ScreenName:      t.User.ScreenName,
				Name:            t.User.Name,
				ProfileImageUrl: t.User.ProfileImageURLHttps,
				Label:           model.UNLABELED,
			}
			tweets = append(tweets, &tweet)
		}
	}
	return &model.ReferringTweets{Tweets: tweets, Count: len(tweets)}, nil
}

func setReferringTweets(app service.GoActiveLearningApp, listName string, blacklistFilename string) error {
	blacklist, err := GetScreenNameList(blacklistFilename)
	if err != nil {
		fmt.Println(err.Error())
	}

	examples, err := app.GetRecommendation(listName)
	if err != nil {
		return err
	}
	app.AttachMetadataIncludingFeatureVector(examples, 0, 10000)

	m, err := app.FindLatestMIRAModel(classifier.TWITTER)
	if err != nil {
		return err
	}

	for _, e := range examples {
		if e.UpdatedAt.Add(time.Hour * 240).Before(time.Now()) {
			continue
		}
		fmt.Println(e.FinalUrl)
		tweets, err := GetReferringTweets(e.FinalUrl, blacklist)
		if err != nil {
			u := e.FinalUrl
			if e.FinalUrl == "" {
				u = e.Url
			}
			fmt.Printf("cannot retrieve %s: %s", u, err.Error())
			continue
		}
		for _, t := range tweets.Tweets {
			et := tweet_feature.GetExampleAndTweet(e, t)
			t.Score = m.PredictScore(et.GetFeatureVector())
		}
		e.ReferringTweets = tweets
		if err = app.UpdateOrCreateReferringTweets(e); err != nil {
			fmt.Printf("cannot update %s: %s", e.Url, err.Error())
		}
	}
	return nil
}

func doSetReferringTweets(c *cli.Context) error {
	listName := c.String("listname")
	blacklistFilename := c.String("blacklist-filename")

	app, err := service.NewDefaultApp()
	if err != nil {
		return err
	}
	defer app.Close()
	return setReferringTweets(app, listName, blacklistFilename)
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
		cli.StringFlag{Name: "blacklist-filename", Usage: "Filename of blacklist usernames"},
	},
}
