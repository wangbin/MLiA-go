package main

import (
	"../Ch02/knn"
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"os"
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

	labels := []string{"no surfacing", "flippers"}
	myTree := createTree(dataSet, labels)
	fmt.Println(myTree.Classify(labels, knn.NewPoint(1, 0)))
	fmt.Println(myTree.Classify(labels, knn.NewPoint(1, 1)))
	StoreTree(myTree, "myTree.dat")
	node, _ := GrabTree("myTree.dat")
	fmt.Println(node)
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

type Node struct {
	Name     string
	SubNodes map[float64]*Node
}

func (node Node) String() string {
	if len(node.SubNodes) == 0 {
		return fmt.Sprintf("'%v'", node.Name)
	}
	var buffer bytes.Buffer
	for key, subNode := range node.SubNodes {
		buffer.WriteString(fmt.Sprintf("%v: %v, ", key, subNode))
	}
	return fmt.Sprintf("{'%v': {%s}}", node.Name, buffer.String()[:buffer.Len()-2])
}

func (node Node) Classify(featLabels []string, testVec *knn.Point) string {
	var featIndex int
	for index, label := range featLabels {
		if label == node.Name {
			featIndex = index
			break
		}
	}
	var classLabel string
	for key, node := range node.SubNodes {
		if testVec.Positions[featIndex] == key {
			if len(node.SubNodes) == 0 {
				classLabel = node.Name
			} else {
				classLabel = node.Classify(featLabels, testVec)
			}
		}
	}
	return classLabel
}

func StoreTree(node *Node, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(node)
	return err
}

func GrabTree(filename string) (*Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	node := &Node{}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(node)
	return node, err
}

func majorityCnt(points []*knn.Point) (result string) {
	classCount := make(map[string]int)
	for _, point := range points {
		if _, ok := classCount[point.Label]; ok {
			classCount[point.Label] += 1
		} else {
			classCount[point.Label] = 1
		}
	}
	maxCount := 0
	for label, count := range classCount {
		if count > maxCount {
			maxCount = count
			result = label
		}
	}
	return
}

func createTree(dataSet *knn.Group, labels []string) *Node {
	var myTree *Node

	if len(dataSet.Points[0].Positions) == 0 {
		myTree = &Node{Name: dataSet.Points[0].Label}
		return myTree
	}

	label := dataSet.Points[0].Label
	isClassEqual := true
	for _, point := range dataSet.Points[1:] {
		if point.Label != label {
			isClassEqual = false
			break
		}
	}
	if isClassEqual {
		myTree = &Node{Name: majorityCnt(dataSet.Points)}
		return myTree
	}

	newLabels := make([]string, len(labels))
	copy(newLabels, labels)
	bestFeat := chooseBestFeatureToSplit(dataSet)
	bestFeatLabel := labels[bestFeat]
	newLabels = append(newLabels[:bestFeat], newLabels[bestFeat+1:]...)
	uniqueVals := make(map[float64]int)
	for _, point := range dataSet.Points {
		uniqueVals[point.Positions[bestFeat]] = 1
	}
	myTree = &Node{Name: bestFeatLabel}
	myTree.SubNodes = make(map[float64]*Node)
	for value := range uniqueVals {
		myTree.SubNodes[value] = createTree(
			splitDataSet(dataSet, bestFeat, value), newLabels)
	}
	return myTree
}
