package ahocorasick

import (
	"sync"
	"github.com/anknown/ahocorasick"
)

var (
	once sync.Once
	machine *goahocorasick.Machine
)

var keywords = []string{
	"機械学習",
	"machine learning",

	"深層強化学習",

	"自然言語処理",
	"natural language processing",

	"言語モデル",
	"language model",

	"機械翻訳",
	"machine translation",

	"parsing",

	"異常検知",
	"anomaly detection",

	"変化検知",
	"change point detection",

	"物体検出",

	"データ基盤",

	"論文",
	"survey",
	"チュートリアル",
	"tutorial",

	"acl",
	"emnlp",
	"naacl",
	"coling",

	"nips",
	"cvpr",
	"icml",
	"recsys",
	"icdm",
	"kdd",
	"iclr",

	"ロジステック回帰",
	"logistic regression",
	"ガウシアンプロセス",
	"gaussian process",
	"svm",

	"factorization machine",
	"matrix factorization",
	"行列分解",

	"ベイズ",
	"bayesian",

	"ベイズ最適化",
	"bayesian optimization",

	"python",
	"julia",
	"tensorflow",
	"keras",
	"chainer",
	"pytorch",
	"allennlp",
	"coreml",
	"sagemaker",
	"colab",
	"scikit-learn",
	"pandas",
	"scipy",
	"hivemall",
	"dask",

	"cnn",
	"rnn",
	"lstm",
	"bert",
	"elmo",
	"gan",
	"u-net",
	"transformer",

	"kaggle",
	"コンペ",
}

func Init() error {
	var err error
	once.Do(func() {
		machine = new(goahocorasick.Machine)
		ks := make([][]rune, 0)
		for _, k := range keywords {
			ks = append(ks, []rune(k))
		}
		err = machine.Build(ks)
		if err != nil {
			return
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func SearchKeywords(content string) []string {
	terms := machine.MultiPatternSearch([]rune(content), false)
	foundKeywords := make([]string, 0)
	for _, t := range terms {
		foundKeywords = append(foundKeywords, string(t.Word))
	}
	return foundKeywords
}
