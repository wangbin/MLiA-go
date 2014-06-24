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
	listOPosts, listClasses := loadDataSet()
	myVocabList := bayes.CreateVocabList(listOPosts)
	trainMat := make([][]int, 0)
	for _, postinDoc := range listOPosts {
		trainMat = append(trainMat, bayes.SetOfWords2Vec(myVocabList, postinDoc))
	}
	p0V, p1V, pAb := bayes.TrainNB0(trainMat, listClasses)
	fmt.Printf("p0V = %f\np1V = %f\npAb = %f\n", p0V, p1V, pAb)
}
