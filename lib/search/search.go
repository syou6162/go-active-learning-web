package search

import (
	"strings"

	"github.com/syou6162/go-active-learning-web/lib/ahocorasick"
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
	app.AttachLightMetadata(examples)
	return examples, nil
}

func SearchSimilarExamples(app service.GoActiveLearningApp, query string, maxOutputs int) (model.Examples, []string, error) {
	keywords := ahocorasick.SearchKeywords(strings.ToLower(query))
	examples, err := app.SearchExamplesByKeywords(keywords, "ANY", maxOutputs)
	if err != nil {
		return nil, make([]string, 0), err
	}
	app.AttachLightMetadata(examples)
	return examples, keywords, nil
}
