package ahocorasick

import (
	"sync"

	goahocorasick "github.com/anknown/ahocorasick"
	"github.com/syou6162/go-active-learning/lib/util"
)

var (
	once                  sync.Once
	machinePriorityHigh   *goahocorasick.Machine
	machinePriorityMedium *goahocorasick.Machine
	machinePriorityLow    *goahocorasick.Machine
)

var keywordsPriorityHigh = []string{
	"active learning",

	"差分プライバシー",
	"differential privacy",

	"アノテーション",
	"annotation",

	"変化検知",
	"change point detection",
	"外れ値検出",
	"outlier detection",
	"特異スペクトル変換法",

	"姿勢推定",
	"yolo",

	"lambdamart",
	"approximate nearest neighbor",

	"redash",
	"re:dash",
	"data studio",
	"google analytics",
	"athena",
	"bigquery",
	"embulk",
	"spark",
	"sql",
	"fluentd",
	"digdag",
	"airflow",
	"argo",
	"データパイプライン",
	"fisher information",

	"hivemall",
	"dask",
	"nnabla",
	"edward",

	"ngboost",
	"sentencepiece",
	"openpose",
	"scdv",
	"tslearn",
	"econml",
	"deepgbm",
	"sudachi",
	"ncrfpp",
	"hsic",
	"irgan",
	"dialogpt",
	"tinyvideonet",
	"blazeface",
	"singan",

	"cookiecutter",
	"datarobot",
	"dataflow",
	"mlflow",
}

var keywordsPriorityMedium = []string{
	"深層強化学習",
	"reinforcement learning",

	"言語モデル",
	"language model",
	"language generation",
	"text generation",

	"要約",
	"summarization",
	"attention",

	"構文解析",
	"parsing",

	"structure learning",
	"structured learning",
	"構造学習",

	"固有表現抽出",
	"固有表現認識",
	"conditional random field",
	"named entity recognition",
	"relation extraction",
	"関係抽出",
	"系列ラベリング",
	"distant supervision",
	"情報検索",
	"information retrieval",
	"画像検索",
	"質問応答",
	"question answering",

	"対話",
	"dialogue",
	"dialog",
	"マルチモーダル",
	"multimodal",

	"異常検知",
	"anomaly detection",

	"metric learning",
	"距離学習",
	"双曲空間",

	"learning to rank",
	"ランキング学習",

	"bandit",
	"時系列",
	"time series",

	"医療",
	"medical",
	"fairness",

	"データウェアハウス",
	"データレイク",
	"data lake",
	"data warehouse",
	"データマート",
	"data mart",

	"conll",
	"semeval",

	"recsys",
	"icdm",
	"wsdm",

	"factorization machine",
	"matrix factorization",

	"変分ベイズ",

	"ベイズ最適化",
	"bayesian optimization",
	"凸最適化",
	"convex optimization",

	"転移学習",
	"transfer learning",
	"ドメイン適用",
	"domain adaptation",
	"data augmentation",

	"tensorboard",
	"allennlp",
	"sagemaker",
	"colab",
	"scikit-learn",
	"sklearn",
	"pandas",
	"scipy",
	"julia",

	"lightgbm",
	"xgboost",
	"lime",
	"shap",
	"fairseq",
	"tensor2tensor",
	"stacknet",
	"fsgan",
	"deepfake",
	"molgan",
	"scibert",
	"fastspeech",
	"blazeface",
	"mecab",
	"neo4j",
	"leam",
	"spacy",
	"flair",
	"continuous delivery",

	"albert",
	"elmo",
	"electra",
	"u-net",
	"xlnet",
	"stylegan",
	"octave",
	"gcn",
	"cyclegan",
	"mobilenet",
	"noise2void",
	"snorkel",
	"yolact",
	"kelpnet",

	"kubeflow",
	"kafka",
	"kinesis",
	"elasticsearch",
	"amazon forecast",
	"splunk",
	"mlops",
	"automl",
	"機械学習工学",
	"reproducibility",

	"kaggle days",
	"kaggle days tokyo",

	"abテスト",
	"検定",
	"statistics",
	"統計学",

	"音声変換",
	"link prediction",
}

var keywordsPriorityLow = []string{
	"機械学習",
	"machine learning",
	"deep learning",
	"深層学習",
	"end-to-end",
	"データサイエンス",
	"データサイエンティスト",
	"データエンジニアリング",

	"強化学習",
	"ベルマン方程式",
	"neural net",
	"neural network",
	"ニューラルネット",
	"ディープラーニング",
	"tpu",
	"gpu",

	"自然言語処理",
	"言語処理",
	"natural language processing",
	"nlp",

	"機械翻訳",
	"machine translation",

	"seq2seq",

	"テキストマイニング",
	"text mining",
	"暗号化",
	"anonymization",
	"匿名化",

	"コンピュータビジョン",
	"computer vision",
	"画像生成",
	"image generation",
	"fashion",

	"信号処理",
	"音声認識",
	"音声合成",

	"community detection",

	"物体検出",
	"semantic segmentation",
	"segmentation",
	"セグメンテーション",

	"cifar",
	"imagenet",

	"レコメンド",
	"レコメンデーション",
	"recommender",
	"recommendation",
	"協調フィルタリング",
	"factorization machine",
	"collaborative filtering",

	"線形代数",
	"linear algebra",

	"データ基盤",
	"データ分析",
	"データセット",
	"dataset",
	"コーパス",
	"corpus",

	"論文",
	"survey",
	"サーベイ",
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
	"kdd",
	"iclr",
	"ijcai",
	"aaai",
	"siggraph",
	"icra",
	"eccv",
	"interspeech",

	"miru",
	"yans",
	"pycon",

	"ロジステック回帰",
	"logistic regression",
	"ガウシアンプロセス",
	"ガウス過程",
	"gaussian process",
	"svm",

	"行列分解",
	"pca",
	"svd",
	"lda",
	"トピックモデル",
	"topic model",
	"混合ガウス分布",
	"gaussian mixture model",
	"clustering",
	"クラスタリング",

	"ベイズ",
	"bayesian",
	"生成モデル",

	"半教師あり学習",
	"事前学習",
	"学習済みモデル",
	"pretrain",

	"pagerank",

	"圧縮",
	"高速化",
	"可視化",
	"visualizer",

	"open source",
	"オープンソース",
	"python",
	"c++",
	"tensorflow",
	"keras",
	"chainer",
	"onnx",
	"chainerrl",
	"pytorch",
	"coreml",
	"numpy",
	"jupyter",
	"決定木",
	"random forest",
	"ランダムフォレスト",
	"optuna",
	"pytext",

	"cnn",
	"rnn",
	"lstm",
	"autoencoder",
	"オートエンコーダー",
	"vae",
	"bert",
	"gan",
	"generative adversarial network",
	"敵対的生成ネットワーク",
	"generative adversarial training",
	"adversarial example",
	"transformer",
	"vgg",
	"adam",
	"early stopping",
	"分散学習",
	"勾配法",
	"biggan",
	"bertsum",

	"埋め込み",
	"分散表現",
	"embedding",
	"glove",
	"word2vec",
	"ナレッジグラフ",
	"knowledge graph",
	"embedrank",

	"gke",
	"gae",
	"gcp",
	"aws",
	"docker",
	"コンテナ",
	"kubernetes",

	"kaggle",
	"前処理",
	"signate",
	"コンペ",
	"competition",
	"bagging",

	"self driving",
	"自動運転",
	"robotics",
}

func build(machine *goahocorasick.Machine, keywords []string) error {
	machine = new(goahocorasick.Machine)
	ks := make([][]rune, 0)
	for _, k := range keywords {
		ks = append(ks, []rune(k))
	}
	err := machine.Build(ks)
	if err != nil {
		return err
	}
	return nil
}

func Init() error {
	var finalErr error
	once.Do(func() {
		if err := build(machinePriorityHigh, keywordsPriorityHigh); err != nil {
			finalErr = err
			return
		}
		if err := build(machinePriorityMedium, keywordsPriorityMedium); err != nil {
			finalErr = err
			return
		}
		if err := build(machinePriorityLow, keywordsPriorityLow); err != nil {
			finalErr = err
			return
		}
	})
	if finalErr != nil {
		return finalErr
	}
	return nil
}

func searchKeywords(machine *goahocorasick.Machine, content string) []string {
	terms := machine.MultiPatternSearch([]rune(content), false)
	foundKeywords := make([]string, 0)
	for _, t := range terms {
		foundKeywords = append(foundKeywords, string(t.Word))
	}
	return util.RemoveDuplicate(foundKeywords)
}

func SearchKeywords(content string) []string {
	if terms := searchKeywords(machinePriorityHigh, content); len(terms) != 0 {
		return terms
	}
	if terms := searchKeywords(machinePriorityMedium, content); len(terms) != 0 {
		return terms
	}
	return searchKeywords(machinePriorityLow, content)
}
