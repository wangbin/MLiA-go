package main

import (
	"./bayes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func main() {
	spamTest()
}

func textParse(bigString string) (tokens []string) {
	re := regexp.MustCompile(`\W+`)

	for _, token := range re.Split(bigString, -1) {
		if len(token) > 2 {
			tokens = append(tokens, strings.ToLower(token))
		}
	}
	return
}

func spamTest() {
	var buf []byte
	var err error
	var wordList []string
	docList := make([][]string, 0)
	fullText := make([]string, 0)
	classList := make([]int, 0)
	for i := 1; i < 26; i++ {
		buf, err = ioutil.ReadFile(fmt.Sprintf("email/spam/%d.txt", i))
		if err != nil {
			panic(err)
		}
		wordList = textParse(string(buf))
		docList = append(docList, wordList)
		fullText = append(fullText, wordList...)
		classList = append(classList, 1)
		buf, err = ioutil.ReadFile(fmt.Sprintf("email/ham/%d.txt", i))
		if err != nil {
			panic(err)
		}
		wordList = textParse(string(buf))
		docList = append(docList, wordList)
		fullText = append(fullText, wordList...)
		classList = append(classList, 0)
	}
	vocabList := bayes.CreateVocabList(docList)
	trainSet := make([]int, 0)
	testSet := make([]int, 0)
	var randIndex int
	for i := 0; i < 50; i++ {
		trainSet = append(trainSet, i)
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		randIndex = rand.Intn(len(trainSet))
		testSet = append(testSet, trainSet[randIndex])
		trainSet = append(trainSet[:randIndex], trainSet[randIndex+1:]...)
	}
	trainMat := make([][]int, 0)
	trainClasses := make([]int, 0)
	for _, docIndex := range trainSet {
		trainMat = append(trainMat, bayes.SetOfWords2Vec(vocabList, docList[docIndex]))
		trainClasses = append(trainClasses, classList[docIndex])
	}
	p0V, p1V, pSpam := bayes.TrainNB0(trainMat, trainClasses)
	errorCount := 0
	var wordVector []int
	for _, docIndex := range testSet {
		wordVector = bayes.SetOfWords2Vec(vocabList, docList[docIndex])
		if bayes.ClassifyNB(wordVector, p0V, p1V, pSpam) != classList[docIndex] {
			errorCount += 1
			fmt.Printf("classification error %v\n", docList[docIndex])
		}
	}
	fmt.Printf("the error rate is %f\n", float64(errorCount)/float64(len(testSet)))
}
