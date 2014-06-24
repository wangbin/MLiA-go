package main

import (
	"./bayes"
	"fmt"
)

func loadDataSet() (postingList [][]string, classVec []int) {
	postingList = [][]string{
		[]string{"my", "dog", "has", "flea", "problems", "help", "please"},
		[]string{"maybe", "not", "take", "him", "to", "dog", "park", "stupid"},
		[]string{"my", "dalmation", "is", "so", "cute", "I", "love", "him"},
		[]string{"stop", "posting", "stupid", "worthless", "garbage"},
		[]string{"mr", "licks", "ate", "my", "steak", "how", "to", "stop", "him"},
		[]string{"quit", "buying", "worthless", "dog", "food", "stupid"},
	}
	classVec = []int{0, 1, 0, 1, 0, 1}
	return
}

func main() {
	testingNB()
}

func testingNB() {
	listOPosts, listClasses := loadDataSet()
	myVocabList := bayes.CreateVocabList(listOPosts)
	trainMat := make([][]int, 0)
	for _, postinDoc := range listOPosts {
		trainMat = append(trainMat, bayes.SetOfWords2Vec(myVocabList, postinDoc))
	}
	p0V, p1V, pAb := bayes.TrainNB0(trainMat, listClasses)
	var testEntry []string
	testEntry = []string{"love", "my", "dalmation"}
	thisDoc := bayes.SetOfWords2Vec(myVocabList, testEntry)
	fmt.Printf("%v classified as: %d\n", testEntry, bayes.ClassifyNB(thisDoc, p0V, p1V, pAb))
	testEntry = []string{"stupid", "garbage"}
	thisDoc = bayes.SetOfWords2Vec(myVocabList, testEntry)
	fmt.Printf("%v classified as: %d\n", testEntry, bayes.ClassifyNB(thisDoc, p0V, p1V, pAb))
}
