package update_model

import (
	"os"
	"time"

	"github.com/codegangsta/cli"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/syou6162/go-active-learning/lib/util/converter"
)

func postEvaluatedMetricsToMackerel(accuracy float64, precision float64, recall float64, fvalue float64) error {
	apiKey := os.Getenv("MACKEREL_API_KEY")
	serviceName := os.Getenv("MACKEREL_SERVICE_NAME")
	if apiKey == "" || serviceName == "" {
		return nil
	}

	client := mkr.NewClient(apiKey)
	now := time.Now().Unix()
	err := client.PostServiceMetricValues(serviceName, []*mkr.MetricValue{
		{
			Name:  "evaluation.accuracy",
			Time:  now,
			Value: accuracy,
		},
		{
			Name:  "evaluation.precision",
			Time:  now,
			Value: precision,
		},
		{
			Name:  "evaluation.recall",
			Time:  now,
			Value: recall,
		},
		{
			Name:  "evaluation.fvalue",
			Time:  now,
			Value: fvalue,
		},
	})
	return err
}

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
	m, err := classifier.NewMIRAClassifierByCrossValidation(classifier.EXAMPLE, converter.ConvertExamplesToLearningInstances(examples))
	if err != nil {
		return err
	}

	if err := app.InsertMIRAModel(*m); err != nil {
		return err
	}
	if err := postEvaluatedMetricsToMackerel(m.Accuracy, m.Precision, m.Recall, m.Fvalue); err != nil {
		return err
	}
	return nil
}

var CommandUpdateModel = cli.Command{
	Name:  "update-model",
	Usage: "update model",
	Description: `
Update model.
`,
	Action: doUpdateModel,
}
