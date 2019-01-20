package update_model

import (
	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/syou6162/go-active-learning/lib/util/converter"
)

func doUpdateExampleModel(c *cli.Context) error {
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
	if err := postEvaluatedMetricsToMackerel("evaluation", m.Accuracy, m.Precision, m.Recall, m.Fvalue); err != nil {
		return err
	}
	return nil
}

var CommandUpdateExampleModel = cli.Command{
	Name:  "update-example-model",
	Usage: "update example model",
	Description: `
Update example model.
`,
	Action: doUpdateExampleModel,
}
