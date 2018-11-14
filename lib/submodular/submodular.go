package submodular

import (
	"math"

	"fmt"

	"github.com/pbnjay/clustering"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/feature"
)

func extractFeature(e example.Example) feature.FeatureVector {
	result := feature.FeatureVector{}
	result = append(result, feature.ExtractNounFeaturesWithoutPrefix(e.Title)...)
	result = append(result, feature.ExtractHostFeature(e.FinalUrl))
	return result
}

func SelectSubExamplesBySubModular(whole example.Examples, sizeConstraint int, alpha float64, r float64, lambda float64) example.Examples {
	selected := example.Examples{}
	remainings := whole
	simMat := GetSimilarityMatrixByTFIDF(whole)
	clusters := GetClusters(simMat, whole)

	clusters.EachCluster(-1, func(cluster int) {
		clusters.EachItem(cluster, func(e clustering.ClusterItem) {
			switch x := e.(type) {
			case *example.Example:
				fmt.Printf("%d %s %s\n", cluster, x.Title, x.Url) // for debuging
			}
		})
	})

	for {
		if len(selected) >= sizeConstraint || len(remainings) == 0 {
			break
		}
		argmax := SelectBestExample(simMat, clusters, remainings, selected, whole, alpha, r, lambda)
		selected = append(selected, remainings[argmax])
		remainings = append(remainings[:argmax], remainings[argmax+1:]...)
	}
	// (1 - 1/e)/2の保証を与えるためにはもうちょっと頑張る必要があるが、省略している
	// http://www.anthology.aclweb.org/E/E09/E09-1089.pdf
	return selected
}

func GetClusters(simMat SimilarityMatrix, whole example.Examples) clustering.ClusterSet {
	distMap := clustering.DistanceMap{}
	for _, e1 := range whole {
		m := make(map[clustering.ClusterItem]float64)
		for _, e2 := range whole {
			similarity := GetCosineSimilarity(simMat, e1, e2)
			// convert similarity to distance
			m[e2] = 1.0 / math.Max(0.00001, similarity)
		}
		distMap[e1] = m
	}
	clusters := clustering.NewDistanceMapClusterSet(distMap)

	clustering.Cluster(clusters, clustering.Threshold(100.0), clustering.CompleteLinkage())
	return clusters
}

// ref: http://www.lr.pi.titech.ac.jp/~morita/YANS.pdf
func DiversityFunction(mat SimilarityMatrix, clusters clustering.ClusterSet, subset example.Examples, whole example.Examples) float64 {
	sum := 0.0
	// Enumerate clusters and print members
	clusters.EachCluster(-1, func(cluster int) {
		tmp := 0.0
		clusters.EachItem(cluster, func(e clustering.ClusterItem) {
			switch x := e.(type) {
			case *example.Example:
				for _, s := range subset {
					if s.Url == x.Url {
						tmp += coverageFunction(mat, x, whole) / float64(len(whole))
						continue
					}
				}
			}
		})
		sum += math.Sqrt(tmp)
	})
	return sum
}

// ref: http://aclweb.org/anthology/N10-1134
// Algorithm 1 line4
func SelectBestExample(mat SimilarityMatrix, clusters clustering.ClusterSet, remainings example.Examples, selected example.Examples, whole example.Examples, alpha float64, r float64, lambda float64) int {
	maxScore := math.Inf(-1)
	argmax := 0
	c2 := CoverageFunction(mat, selected, whole, alpha) + lambda*DiversityFunction(mat, clusters, selected, whole)
	for idx, remaining := range remainings {
		subset := example.Examples{}
		for _, e := range selected {
			subset = append(subset, e)
		}
		subset = append(subset, remaining)
		c1 := CoverageFunction(mat, subset, whole, alpha) + lambda*DiversityFunction(mat, clusters, subset, whole)
		fv := extractFeature(*remaining)
		if len(fv) == 0 {
			continue
		}
		score := (c1 - c2) / math.Pow(float64(len(fv)), r)
		if score >= maxScore {
			argmax = idx
			maxScore = score
		}
	}
	return argmax
}

func CoverageFunction(mat SimilarityMatrix, subset example.Examples, whole example.Examples, alpha float64) float64 {
	sum := 0.0
	for _, e := range whole {
		sum += math.Min(
			coverageFunction(mat, e, subset),
			alpha*coverageFunction(mat, e, whole),
		)
	}
	return sum
}

func coverageFunction(mat SimilarityMatrix, example *example.Example, examples example.Examples) float64 {
	sum := 0.0
	for _, e := range examples {
		sum += GetCosineSimilarity(mat, e, example)
	}
	return sum
}

type SimilarityMatrix map[string]float64

func GetSimilarityMatrixByTFIDF(examples example.Examples) SimilarityMatrix {
	idf := GetIDF(examples)

	dfByURL := make(map[string]map[string]float64)
	sumByUrl := make(map[string]float64)
	for _, e := range examples {
		df := GetDF(*e)
		dfByURL[e.Url] = df

		sum := 0.0
		for k, v := range df {
			sum += v * v * idf[k] * idf[k]
		}
		sumByUrl[e.Url] = sum
	}

	mat := SimilarityMatrix{}
	for _, e1 := range examples {
		df1 := dfByURL[e1.Url]
		s1 := math.Sqrt(sumByUrl[e1.Url])
		if s1 == 0.0 {
			continue
		}

		for _, e2 := range examples {
			df2 := dfByURL[e2.Url]
			s2 := math.Sqrt(sumByUrl[e2.Url])
			if s2 == 0.0 {
				continue
			}

			s := 0.0
			for k, v := range df2 {
				s += v * df1[k] * idf[k] * idf[k]
			}
			mat[e1.Url+"+"+e2.Url] = s / (s1 * s2)
		}
	}
	return mat
}

func GetCosineSimilarity(mat SimilarityMatrix, e1 *example.Example, e2 *example.Example) float64 {
	return mat[e1.Url+"+"+e2.Url]
}

func GetDF(example example.Example) map[string]float64 {
	df := make(map[string]float64)
	n := 0.0
	fv := extractFeature(example)

	for _, f := range fv {
		df[f]++
		n++
	}

	for k, v := range df {
		df[k] = v / n
	}
	return df
}

func GetIDF(examples example.Examples) map[string]float64 {
	idf := make(map[string]float64)
	cnt := make(map[string]float64)
	n := float64(len(examples))

	for _, e := range examples {
		fv := extractFeature(*e)
		for _, f := range fv {
			cnt[f]++
		}
	}

	for k, v := range cnt {
		idf[k] = math.Log(n/v) + 1
	}
	return idf
}
