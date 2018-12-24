package update_model

import (
	"os"
	"time"

	"github.com/codegangsta/cli"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
)

func doUpdateModel(c *cli.Context) error {
	app, err := service.NewDefaultApp()
	if err != nil {
		return err
	}
	defer app.Close()

	examples, err := app.SearchLabeledExamples(100000)
	if err != nil {
		return err
	}
	err = postNumOfPositiveAndNegativeExamplesToMackerel(examples)
	if err != nil {
		return err
	}

	okExamples := util.FilterStatusCodeOkExamples(examples)
	if err = app.AttachMetadata(okExamples); err != nil {
		return err
	}

	notOkExamples := util.FilterStatusCodeNotOkExamples(examples)
	app.Fetch(notOkExamples)
	for _, e := range util.FilterStatusCodeOkExamples(notOkExamples) {
		if err := app.UpdateFeatureVector(e); err != nil {
			return err
		}
	}
	examples = util.FilterStatusCodeOkExamples(examples)
	m, err := classifier.NewMIRAClassifierByCrossValidation(examples)
	if err != nil {
		return err
	}

	return app.InsertMIRAModel(*m)
}

func postNumOfPositiveAndNegativeExamplesToMackerel(examples model.Examples) error {
	apiKey := os.Getenv("MACKEREL_API_KEY")
	serviceName := os.Getenv("MACKEREL_SERVICE_NAME")
	if apiKey == "" || serviceName == "" {
		return nil
	}

	numPos := 0
	numNeg := 0
	for _, e := range examples {
		if e.Label == model.POSITIVE {
			numPos++
		} else if e.Label == model.NEGATIVE {
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

var CommandUpdateModel = cli.Command{
	Name:  "update-model",
	Usage: "update model",
	Description: `
Update model.
`,
	Action: doUpdateModel,
}
