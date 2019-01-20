package update_model

import (
	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/feature/tweet"
)

func doUpdateTweetModel(c *cli.Context) error {
	app, err := service.NewDefaultApp()
	if err != nil {
		return err
	}
	defer app.Close()

	tweets, err := app.SearchReferringTweets(30000)
	if err != nil {
		return err
	}
	exampleIds := make([]int, 0)
	for _, t := range tweets {
		exampleIds = append(exampleIds, t.ExampleId)
	}

	examples, err := app.SearchExamplesByIds(exampleIds)
	if err != nil {
		return err
	}
	exampleById := make(map[int]*model.Example)
	for _, e := range examples {
		exampleById[e.Id] = e
	}

	instances := classifier.LearningInstances{}
	for _, t := range tweets {
		e := exampleById[t.ExampleId]
		et := tweet_feature.GetExampleAndTweet(e, t)
		instances = append(instances, &et)
	}
	m, err := classifier.NewMIRAClassifierByCrossValidation(classifier.TWITTER, instances)
	if err != nil {
		return err
	}

	if err := app.InsertMIRAModel(*m); err != nil {
		return err
	}
	if err := postEvaluatedMetricsToMackerel("tweet_evaluation", m.Accuracy, m.Precision, m.Recall, m.Fvalue); err != nil {
		return err
	}
	return nil
}

var CommandUpdateTweetModel = cli.Command{
	Name:  "update-tweet-model",
	Usage: "update tweet model",
	Description: `
Update tweet model.
`,
	Action: doUpdateTweetModel,
}
