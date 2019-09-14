package util

import "github.com/syou6162/go-active-learning/lib/model"

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func LightenExamples(examples model.Examples) {
	for _, example := range examples {
		example.Fv = make([]string, 0)
		r := []rune(example.Body)
		example.Body = string(r[0:Min(500, len(r))])
	}
}
