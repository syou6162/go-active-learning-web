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
	"データサイエンス",
	"データサイエンティスト",
	"データエンジニアリング",

	"強化学習",
	"深層強化学習",
	"reinforcement learning",
	"ベルマン方程式",
	"active learning",
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

	"言語モデル",
	"language model",
	"language generation",
	"text generation",

	"機械翻訳",
	"machine translation",

	"要約",
	"summarization",
	"attention",

	"parsing",

	"structure learning",
	"structured learning",
	"構造学習",

	"固有表現抽出",
	"固有表現認識",
	"seq2seq",
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

	"テキストマイニング",
	"text mining",
	"差分プライバシー",
	"differential privacy",
	"暗号化",
	"anonymization",
	"匿名化",

	"対話",
	"dialogue",
	"dialog",
	"マルチモーダル",
	"multimodal",

	"アノテーション",
	"annotation",

	"コンピュータビジョン",
	"computer vision",
	"画像生成",
	"image generation",
	"fashion",

	"信号処理",
	"音声認識",
	"音声合成",

	"異常検知",
	"anomaly detection",
	"外れ値検出",
	"outlier detection",
	"特異スペクトル変換法",

	"変化検知",
	"change point detection",

	"community detection",

	"metric learning",
	"距離学習",
	"双曲空間",

	"物体検出",
	"semantic segmentation",
	"segmentation",
	"セグメンテーション",
	"姿勢推定",
	"yolo",

	"cifar",
	"imagenet",

	"レコメンド",
	"レコメンデーション",
	"recommender",
	"recommendation",
	"協調フィルタリング",
	"factorization machine",
	"collaborative filtering",
	"learning to rank",
	"ランキング学習",
	"lambdamart",
	"approximate nearest neighbor",

	"bandit",

	"時系列",
	"time series",

	"医療",
	"medical",
	"fairness",

	"線形代数",
	"linear algebra",

	"データ基盤",
	"データ分析",
	"データウェアハウス",
	"データレイク",
	"data lake",
	"data warehouse",
	"データマート",
	"data mart",
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
	"conll",
	"semeval",

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
	"clustering",
	"クラスタリング",

	"ベイズ",
	"bayesian",
	"変分ベイズ",
	"生成モデル",
	"fisher information",

	"ベイズ最適化",
	"bayesian optimization",
	"凸最適化",
	"convex optimization",

	"半教師あり学習",
	"転移学習",
	"transfer learning",
	"ドメイン適用",
	"domain adaptation",
	"data augmentation",
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
	"julia",
	"c++",
	"tensorflow",
	"tensorboard",
	"keras",
	"chainer",
	"onnx",
	"chainerrl",
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
	"ngboost",
	"決定木",
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
	"fairseq",
	"tensor2tensor",
	"stacknet",
	"fsgan",
	"deepfake",
	"molgan",
	"scibert",
	"fastspeech",
	"ncrfpp",
	"blazeface",
	"mecab",
	"neo4j",
	"leam",
	"spacy",
	"flair",
	"continuous delivery",

	"cnn",
	"rnn",
	"lstm",
	"autoencoder",
	"オートエンコーダー",
	"vae",
	"bert",
	"albert",
	"elmo",
	"electra",
	"gan",
	"generative adversarial network",
	"敵対的生成ネットワーク",
	"generative adversarial training",
	"adversarial example",
	"u-net",
	"transformer",
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
	"vgg",
	"adam",
	"early stopping",
	"分散学習",
	"勾配法",
	"hsic",
	"biggan",
	"bertsum",
	"irgan",
	"dialogpt",
	"tinyvideonet",
	"blazeface",
	"singan",

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
	"kubeflow",
	"sagemaker",
	"s3",
	"kafka",
	"kinesis",
	"elasticsearch",
	"amazon forecast",
	"splunk",
	"mlops",
	"automl",
	"機械学習工学",
	"reproducibility",
	"cookiecutter",
	"datarobot",

	"kaggle",
	"kaggle days",
	"kaggle days tokyo",
	"前処理",
	"signate",
	"コンペ",
	"competition",
	"bagging",
	"abテスト",
	"検定",
	"statistics",
	"統計学",

	"self driving",
	"自動運転",
	"robotics",

	"音声変換",
	"link prediction",
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
