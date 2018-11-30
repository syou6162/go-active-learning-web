package search

import (
	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
)

func Search(app service.GoActiveLearningApp, query string) (model.Examples, error) {
	keywords := feature.ExtractNounFeaturesWithoutPrefix(query)
	examples, err := app.SearchExamplesByKeywords(keywords, "ALL", 100)
	if err != nil {
		return nil, err
	}
	app.AttachLightMetadata(examples)
	return examples, nil
}

func removeOneCharKeywords(keywords []string) []string {
	result := make([]string, 0)
	for _, k := range keywords {
		if len([]rune(k)) > 1 {
			result = append(result, k)
		}
	}
	return result
}

func getUniqueWords(s string) []string {
	return util.RemoveDuplicate(removeOneCharKeywords(feature.ExtractNounFeaturesWithoutPrefix(s)))
}

func SearchSimilarExamples(app service.GoActiveLearningApp, query string, maxOutputs int) (model.Examples, []string, error) {
	keywords := getUniqueWords(query)
	examples, err := app.SearchExamplesByKeywords(keywords, "ANY", maxOutputs)
	if err != nil {
		return nil, make([]string, 0), err
	}
	app.AttachLightMetadata(examples)
	return examples, keywords, nil
}
