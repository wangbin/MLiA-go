package main

import (
	"../Ch02/knn"
	"fmt"
	"math"
)

func main() {
	dataSet := createDataSet()
	fmt.Println(calcShannonEnt(dataSet))
	// dataSet.Points[0].Label = "maybe"
	// fmt.Println(calcShannonEnt(dataSet))

	splitDataSet(dataSet, 0, 1)
}

func createDataSet() *knn.Group {
	group := knn.NewGroup(
		knn.NewPoint(1, 1, "yes"),
		knn.NewPoint(1, 1, "yes"),
		knn.NewPoint(1, 0, "no"),
		knn.NewPoint(0, 1, "no"),
		knn.NewPoint(0, 1, "no"))

	return group
}

func calcShannonEnt(dataSet *knn.Group) (shannonEnt float64) {
	numEntries := len(dataSet.Points)
	labelCounts := make(map[string]int)
	for _, featVec := range dataSet.Points {
		currentLabel := featVec.Label
		if _, ok := labelCounts[currentLabel]; ok {
			labelCounts[currentLabel] += 1
		} else {
			labelCounts[currentLabel] = 1
		}
	}

	for key := range labelCounts {
		prob := float64(labelCounts[key]) / float64(numEntries)
		shannonEnt -= prob * math.Log2(prob)
	}
	return
}

func splitDataSet(dataSet *knn.Group, axis int, value float64) *knn.Group {
	retDataSet := knn.NewGroup()
	for _, point := range dataSet.Points {
		if point.Positions[axis] == value {
			fmt.Println(point)
		}
	}
	return retDataSet
}
