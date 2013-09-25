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

	newDataSet := splitDataSet(dataSet, 0, 1)
	for _, point := range newDataSet.Points {
		fmt.Println(point)
	}

	newDataSet = splitDataSet(dataSet, 0, 0)
	for _, point := range newDataSet.Points {
		fmt.Println(point)
	}

	fmt.Println(chooseBestFeatureToSplit(dataSet))

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
			positions := make([]float64, axis)
			copy(positions, point.Positions[:axis])
			positions = append(positions, point.Positions[axis+1:]...)
			retPoint := &knn.Point{Positions: positions, Label: point.Label}
			retDataSet.Append(retPoint)
		}
	}
	return retDataSet
}

func chooseBestFeatureToSplit(dataSet *knn.Group) int {
	numFeatures := len(dataSet.Points[0].Positions)
	baseEntropy := calcShannonEnt(dataSet)
	bestInfoGain := 0.0
	bestFeature := -1
	for i := 0; i < numFeatures; i++ {
		uniqueVals := make(map[float64]bool)
		for _, point := range dataSet.Points {
			if _, ok := uniqueVals[point.Positions[i]]; !ok {
				uniqueVals[point.Positions[i]] = true
			}
		}
		newEntropy := 0.0
		for value := range uniqueVals {
			subDataSet := splitDataSet(dataSet, i, value)
			prob := float64(len(subDataSet.Points)) / float64(len(dataSet.Points))
			newEntropy += prob * calcShannonEnt(subDataSet)
		}
		infoGain := baseEntropy - newEntropy
		if infoGain > bestInfoGain {
			bestInfoGain = infoGain
			bestFeature = i
		}
	}
	return bestFeature
}
