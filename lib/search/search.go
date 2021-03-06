package search

import (
	"strings"

	"github.com/syou6162/go-active-learning-web/lib/ahocorasick"
	web_util "github.com/syou6162/go-active-learning-web/lib/util"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
)

func Search(app service.GoActiveLearningApp, query string) (model.Examples, error) {
	keywords := make([]string, 0)
	for _, s := range strings.Split(query, " ") {
		if s != "" {
			keywords = append(keywords, s)
		}
	}

	examples, err := app.SearchExamplesByKeywords(keywords, "ALL", 10)
	if err != nil {
		return nil, err
	}
	app.AttachMetadata(examples, 0, 0)
	web_util.LightenExamples(examples)
	return examples, nil
}

func SearchSimilarExamples(app service.GoActiveLearningApp, example *model.Example, maxOutputs int) (model.Examples, []string, error) {
	query := example.Title
	keywords := ahocorasick.SearchKeywords(strings.ToLower(query))
	examples, err := app.SearchRelatedExamples(example)
	if err != nil {
		return nil, keywords, err
	}
	if len(examples) != 0 {
		app.AttachMetadata(examples, 0, 0)
		web_util.LightenExamples(examples)
		return examples, keywords, nil
	}

	examples, err = app.SearchTopAccessedExamples()
	if err != nil {
		return nil, keywords, err
	}
	app.AttachMetadata(examples, 0, 0)
	web_util.LightenExamples(examples)
	return examples, keywords, nil
}
