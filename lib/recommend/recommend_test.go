package recommend_test

import (
	"testing"

	"github.com/syou6162/go-active-learning-web/lib/command"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util/file"
	"github.com/urfave/cli"
)

func TestDoRecommend(t *testing.T) {
	inputFilename := "../../tech_input_example.txt"
	train, err := file.ReadExamples(inputFilename)
	if err != nil {
		t.Error(err)
	}

	a, err := service.NewDefaultApp()
	if err != nil {
		t.Error(err.Error())
	}
	defer a.Close()

	if err = a.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	for _, example := range train {
		if err = a.UpdateOrCreateExample(example); err != nil {
			t.Error(err)
		}
	}

	app := cli.NewApp()
	app.Commands = command.Commands
	args := []string{
		"go-active-learning-web",
		"update-example-model",
	}

	if err := app.Run(args); err != nil {
		t.Error(err)
	}

	args = []string{
		"go-active-learning-web",
		"recommend",
		"--subset-selection",
		"-r=0.75",
		"--size-constraint=20",
		"--score-threshold=-0.1",
		"--listname=general",
	}

	if err := app.Run(args); err != nil {
		t.Error(err)
	}
}
