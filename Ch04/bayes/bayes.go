package bayes

import (
	"math"
	"sort"
)

func CreateVocabList(dataSet [][]string) []string {
	vocabMap := make(map[string]int)
	for _, document := range dataSet {
		for _, vocab := range document {
			vocabMap[vocab] = 1
		}
	}
	vocabSet := make([]string, 0)
	for vocab := range vocabMap {
		vocabSet = append(vocabSet, vocab)
	}
	sort.Strings(vocabSet)
	return vocabSet
}

func SetOfWords2Vec(vocabList []string, input []string) []int {
	returnVec := make([]int, len(vocabList))
	for _, word := range input {
		for index, vocab := range vocabList {
			if vocab == word {
				returnVec[index] = 1
				break
			}
		}
	}
	return returnVec
}

func TrainNB0(trainMatrix [][]int, trainCategory []int) ([]float64, []float64,
	float64) {
	numTrainDocs := len(trainMatrix)
	numWords := len(trainMatrix[0])
	var pAbusive float64
	for _, v := range trainCategory {
		pAbusive += float64(v) / float64(numTrainDocs)
	}
	p0Num := make([]int, numWords)
	p1Num := make([]int, numWords)
	for i := 0; i < numWords; i++ {
		p0Num[i] = 1
		p1Num[i] = 1
	}
	p0Denom := 2
	p1Denom := 2
	for i := 0; i < numTrainDocs; i++ {
		if trainCategory[i] == 1 {
			for j := 0; j < numWords; j++ {
				p1Num[j] += trainMatrix[i][j]
			}
			for _, v := range trainMatrix[i] {
				p1Denom += v
			}
		} else {
			for j := 0; j < numWords; j++ {
				p0Num[j] += trainMatrix[i][j]
			}
			for _, v := range trainMatrix[i] {
				p0Denom += v
			}
		}
	}
	p0Vec := make([]float64, numWords)
	p1Vec := make([]float64, numWords)
	for i := 0; i < numWords; i++ {
		p0Vec[i] = math.Log(float64(p0Num[i]) / float64(p0Denom))
		p1Vec[i] = math.Log(float64(p1Num[i]) / float64(p1Denom))
	}
	return p0Vec, p1Vec, pAbusive
}

func ClassifyNB(vec2Classify []int, p0V []float64, p1V []float64, pClass1 float64) int {
	p1 := math.Log(pClass1)
	p0 := math.Log(1.0 - pClass1)
	for index := range vec2Classify {
		p1 += float64(vec2Classify[index]) * p1V[index]
		p0 += float64(vec2Classify[index]) * p0V[index]
	}
	if p1 > p0 {
		return 1
	}
	return 0
}
