package update_model

import (
	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/classifier"
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

var CommandUpdateModel = cli.Command{
	Name:  "update-model",
	Usage: "update model",
	Description: `
Update model.
`,
	Action: doUpdateModel,
}
