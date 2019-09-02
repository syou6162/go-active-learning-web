package ahocorasick

import (
	"sync"

	goahocorasick "github.com/anknown/ahocorasick"
	"github.com/syou6162/go-active-learning/lib/util"
)

var (
	once    sync.Once
	machine *goahocorasick.Machine
)

var keywords = []string{
	"機械学習",
	"machine learning",
	"deep learning",
	"深層学習",
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
	"language generation",

	"機械翻訳",
	"machine translation",

	"要約",
	"summarization",

	"parsing",

	"固有表現抽出",
	"固有表現認識",
	"named entity recognition",

	"テキストマイニング",
	"text mining",
	"差分プライバシー",
	"differential privacy",

	"対話",
	"dialogue",

	"アノテーション",
	"annotation",

	"コンピュータビジョン",
	"computer vision",
	"画像生成",
	"image generation",
	"fashion",

	"異常検知",
	"anomaly detection",
	"外れ値検出",
	"outlier detection",

	"変化検知",
	"change point detection",

	"metric learning",
	"距離学習",

	"物体検出",
	"semantic segmentation",
	"segmentation",
	"セグメンテーション",
	"姿勢推定",
	"yolo",

	"cifar",
	"imagenet",

	"レコメンド",
	"recommender",
	"recommendation",
	"協調フィルタリング",
	"factorization machine",
	"collaborative filtering",
	"learning to rank",
	"ランキング学習",
	"lambdamart",

	"bandit",

	"時系列",

	"医療",
	"medical",

	"線形代数",
	"linear algebra",

	"データ基盤",
	"データ分析",
	"redash",
	"re:dash",
	"data studio",
	"google analytics",
	"athena",
	"bigquery",
	"embulk",
	"spark",
	"fluentd",
	"digdag",
	"airflow",
	"argo",
	"データパイプライン",

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
	"neurips",
	"cvpr",
	"icml",
	"recsys",
	"icdm",
	"kdd",
	"iclr",
	"ijcai",
	"aaai",
	"wsdm",
	"siggraph",
	"icra",

	"ロジステック回帰",
	"logistic regression",
	"ガウシアンプロセス",
	"ガウス過程",
	"gaussian process",
	"svm",

	"factorization machine",
	"matrix factorization",
	"行列分解",
	"pca",
	"svd",
	"lda",
	"トピックモデル",
	"topic model",
	"混合ガウス分布",
	"gaussian mixture model",

	"ベイズ",
	"bayesian",
	"変分ベイズ",

	"ベイズ最適化",
	"bayesian optimization",

	"半教師あり学習",
	"転移学習",
	"transfer learning",
	"ドメイン適用",
	"domain adaptation",
	"data augmentation",

	"圧縮",
	"高速化",

	"open source",
	"オープンソース",
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
	"numpy",
	"scikit-learn",
	"sklearn",
	"pandas",
	"scipy",
	"hivemall",
	"dask",
	"nnabla",
	"edward",
	"jupyter",
	"lightgbm",
	"xgboost",
	"random forest",
	"ランダムフォレスト",
	"optuna",
	"sentencepiece",
	"lime",
	"shap",
	"openpose",
	"pytext",
	"scdv",
	"tslearn",
	"econml",
	"deepgbm",
	"sudachi",
	"mecab",

	"cnn",
	"rnn",
	"lstm",
	"autoencoder",
	"オートエンコーダー",
	"vae",
	"bert",
	"elmo",
	"gan",
	"generative adversarial network",
	"u-net",
	"transformer",
	"xlnet",
	"stylegan",
	"octave",
	"gcn",
	"noise2void",
	"adam",
	"early stopping",

	"埋め込み",
	"分散表現",
	"embedding",
	"glove",
	"word2vec",
	"ナレッジグラフ",
	"knowledge graph",

	"gke",
	"gae",
	"gcp",
	"aws",
	"docker",
	"コンテナ",
	"kubernetes",
	"kubeflow",
	"sagemaker",
	"s3",
	"elasticsearch",
	"splunk",
	"mlops",
	"automl",
	"機械学習工学",

	"kaggle",
	"signate",
	"コンペ",
	"competition",
	"bagging",
	"abテスト",

	"self driving",
	"自動運転",
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
