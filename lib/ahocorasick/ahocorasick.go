package ahocorasick

import (
	"sync"

	"github.com/anknown/ahocorasick"
	"github.com/syou6162/go-active-learning/lib/util"
)

var (
	once    sync.Once
	machine *goahocorasick.Machine
)

var keywords = []string{
	"機械学習",
	"machine learning",
	"end-to-end",

	"強化学習",
	"深層強化学習",
	"reinforcement learning",

	"自然言語処理",
	"言語処理",
	"natural language processing",
	"nlp",

	"言語モデル",
	"language model",

	"機械翻訳",
	"machine translation",

	"要約",
	"summarization",

	"ner",
	"parsing",

	"テキストマイニング",
	"text mining",

	"コンピュータビジョン",
	"computer vision",

	"異常検知",
	"anomaly detection",
	"外れ値検出",
	"outlier detection",

	"変化検知",
	"change point detection",

	"物体検出",
	"semantic segmentation",

	"cifar",
	"imagenet",

	"レコメンド",
	"recommender",
	"recommendation",

	"bandit",

	"データ基盤",
	"データ分析",
	"redash",
	"athena",
	"bigquery",

	"データセット",
	"dataset",
	"コーパス",
	"corpus",

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
	"pca",
	"svd",

	"ベイズ",
	"bayesian",

	"ベイズ最適化",
	"bayesian optimization",

	"転移学習",
	"transfer learning",
	"ドメイン適用",
	"domain adaptation",

	"圧縮",
	"高速化",

	"python",
	"julia",
	"c++",
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
	"nnabla",
	"edward",
	"jupyter",

	"cnn",
	"rnn",
	"lstm",
	"bert",
	"elmo",
	"gan",
	"u-net",
	"transformer",

	"埋め込み",
	"embedding",
	"glove",
	"word2vec",

	"gke",
	"gcp",
	"aws",
	"docker",
	"kubernetes",
	"kubeflow",

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
	return util.RemoveDuplicate(foundKeywords)
}
