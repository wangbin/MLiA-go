package main

import (
	"./bayes"
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"math/rand"
	"sort"
	"time"
)

const (
	SF_RSS_URL = "http://sfbay.craigslist.org/stp/index.rss"
	NY_RSS_URL = "http://newyork.craigslist.org/stp/index.rss"
	TIMEOUT    = 15
	Threshold  = -6.0
)

var (
	sf_feed []*rss.Item
	ny_feed []*rss.Item
)

func main() {
	getFeeds(SF_RSS_URL)
	getFeeds(NY_RSS_URL)
	getTopWords(ny_feed, sf_feed)
}

func getFeeds(url string) {
	feed := rss.New(TIMEOUT, true, nil, itemHandler)
	if err := feed.Fetch(url, nil); err != nil {
		panic(err)
	}
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	switch feed.Url {
	case SF_RSS_URL:
		sf_feed = newitems
	case NY_RSS_URL:
		ny_feed = newitems
	}
}

type tokenFreq struct {
	token string
	freq  int
}

func (this tokenFreq) String() string {
	return fmt.Sprintf("(%s, %d)", this.token, this.freq)
}

type tokenFreqs []*tokenFreq

func (this tokenFreqs) Len() int {
	return len(this)
}

func (this tokenFreqs) Less(i, j int) bool {
	return this[i].freq < this[j].freq
}

func (this tokenFreqs) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func calcMostFreq(vocabList []string, fullText []string) []*tokenFreq {
	freqDict := make(map[string]int)
	sortedFreq := make(tokenFreqs, 0)
	var tf *tokenFreq
	for _, token := range vocabList {
		if _, ok := freqDict[token]; !ok {
			freqDict[token] = 0
		}
		for _, t := range fullText {
			if t == token {
				freqDict[token] += 1
			}
		}
	}
	for token, freq := range freqDict {
		tf = &tokenFreq{token, freq}
		sortedFreq = append(sortedFreq, tf)
	}
	sort.Sort(sort.Reverse(sortedFreq))
	return sortedFreq[:30]
}

func localWords(feed1, feed0 []*rss.Item) ([]string, []float64, []float64) {
	var wordList []string
	docList := make([][]string, 0)
	fullText := make([]string, 0)
	classList := make([]int, 0)
	feed0Len := len(feed0)
	feed1Len := len(feed1)
	var minLen int
	if feed0Len <= feed1Len {
		minLen = feed0Len
	} else {
		minLen = feed1Len
	}
	for i := 0; i < minLen; i++ {
		wordList = bayes.TextParse(feed0[i].Description)
		docList = append(docList, wordList)
		fullText = append(fullText, wordList...)
		classList = append(classList, 0)
		wordList = bayes.TextParse(feed1[i].Description)
		docList = append(docList, wordList)
		fullText = append(fullText, wordList...)
		classList = append(classList, 1)
	}
	vocabList := bayes.CreateVocabList(docList)
	top30Words := calcMostFreq(vocabList, fullText)
	for _, tf := range top30Words {
		for index := range vocabList {
			if vocabList[index] == tf.token {
				vocabList = append(vocabList[:index], vocabList[index+1:]...)
				break
			}
		}
	}
	trainSet := make([]int, 0)
	testSet := make([]int, 0)
	var randIndex int
	for i := 0; i < minLen*2; i++ {
		trainSet = append(trainSet, i)
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20; i++ {
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
		}
	}
	fmt.Printf("the error rate is %f\n", float64(errorCount)/float64(len(testSet)))
	return vocabList, p0V, p1V
}

type TopWord struct {
	word string
	freq float64
}

type TopWords []*TopWord

func (this TopWords) Len() int {
	return len(this)
}

func (this TopWords) Less(i, j int) bool {
	return this[i].freq < this[j].freq
}

func (this TopWords) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func getTopWords(ny []*rss.Item, sf []*rss.Item) {
	topNY := make(TopWords, 0)
	topSF := make(TopWords, 0)
	vocabList, p0V, p1V := localWords(ny, sf)
	for i := 0; i < len(p0V); i++ {
		if p0V[i] > Threshold {
			topSF = append(topSF, &TopWord{vocabList[i], p0V[i]})
		}
		if p1V[i] > Threshold {
			topNY = append(topNY, &TopWord{vocabList[i], p1V[i]})
		}
	}
	sort.Sort(sort.Reverse(topNY))
	sort.Sort(sort.Reverse(topSF))
	fmt.Println("SF**SF**SF**SF**SF**SF**SF**SF**SF**SF**SF**SF**SF**SF**SF**SF**")
	for _, topWord := range topSF {
		fmt.Println(topWord.word)
	}
	fmt.Println("NY**NY**NY**NY**NY**NY**NY**NY**NY**NY**NY**NY**NY**NY**NY**NY**")
	for _, topWord := range topNY {
		fmt.Println(topWord.word)
	}
}
