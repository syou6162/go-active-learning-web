package update_model

import (
	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/syou6162/go-active-learning/lib/classifier"
	mkr "github.com/mackerelio/mackerel-client-go"
	"os"
	"time"
)

func doUpdateModel(c *cli.Context) error {
	app, err := service.NewDefaultApp()
	if err != nil {
		return err
	}
	defer app.Close()

	examples, err := app.ReadLabeledExamples(100000)
	if err != nil {
		return err
	}
	err = postNumOfPositiveAndNegativeExamplesToMackerel(examples)
	if err != nil {
		return err
	}

	app.Fetch(examples)
	if err := app.UpdateExamplesMetadata(examples); err != nil {
		return err
	}
	examples = util.FilterStatusCodeOkExamples(examples)
	m := classifier.NewMIRAClassifierByCrossValidation(examples)

	if err := app.InsertMIRAModel(*m); err != nil {
		return err
	}

	for _, e := range examples {
		e.Score = m.PredictScore(e.Fv)
		if err := app.UpdateScore(e); err != nil {
			return err
		}
	}
	return nil
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
