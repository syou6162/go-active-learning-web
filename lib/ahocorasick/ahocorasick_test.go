package ahocorasick_test

import (
	"testing"

	"github.com/syou6162/go-active-learning-web/lib/ahocorasick"
)

func TestSearchKeywords(t *testing.T) {
	if err := ahocorasick.Init(); err != nil {
		t.Error(err)
	}
	keywords := ahocorasick.SearchKeywords("yoloを使った機械学習")
	if len(keywords) != 1 {
		t.Errorf("wrong len(keywords): got %v want %v", len(keywords), 1)
	}
}
