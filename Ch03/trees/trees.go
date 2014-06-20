package trees

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"os"
)

type Point struct {
	Positions []interface{}
	Label     string
}

func NewPoint(params ...interface{}) *Point {
	point := new(Point)
	length := len(params)
	var positionsLength int
	switch params[length-1].(type) {
	case string:
		point.Label = params[length-1].(string)
		positionsLength = length - 1
	default:
		positionsLength = length
	}
	point.Positions = make([]interface{}, positionsLength)
	for i, param := range params[:positionsLength] {
		point.Positions[i] = param
	}
	return point
}

func CalcShannonEnt(dataSet []*Point) (shannonEnt float64) {
	numEntries := len(dataSet)
	labelCounts := make(map[string]int)
	for _, featVec := range dataSet {
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

func SplitDataSet(dataSet []*Point, axis int, value interface{}) []*Point {
	retDataSet := make([]*Point, 0)
	for _, point := range dataSet {
		if point.Positions[axis] == value {
			retPositions := make([]interface{}, axis)
			copy(retPositions, point.Positions[:axis])
			retPositions = append(retPositions, point.Positions[axis+1:]...)
			retPoint := &Point{Positions: retPositions, Label: point.Label}
			retDataSet = append(retDataSet, retPoint)
		}
	}
	return retDataSet
}

func ChooseBestFeatureToSplit(dataSet []*Point) int {
	numFeatures := len(dataSet[0].Positions)
	baseEntropy := CalcShannonEnt(dataSet)
	bestInfoGain := 0.0
	bestFeature := -1
	for i := 0; i < numFeatures; i++ {
		// since golang has no set type, use map instead
		uniqueVals := make(map[interface{}]bool)
		for _, point := range dataSet {
			uniqueVals[point.Positions[i]] = true
		}
		newEntropy := 0.0
		for value := range uniqueVals {
			subDataSet := SplitDataSet(dataSet, i, value)
			prob := float64(len(subDataSet)) / float64(len(dataSet))
			newEntropy += prob * CalcShannonEnt(subDataSet)
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
	SubNodes map[interface{}]*Node
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

func (node Node) Classify(featLabels []string, testVec *Point) string {
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

func majorityCnt(points []*Point) (result string) {
	classCount := make(map[string]int)
	for _, point := range points {
		label := point.Label
		if _, ok := classCount[label]; ok {
			classCount[label] += 1
		} else {
			classCount[label] = 1
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

func CreateTree(dataSet []*Point, labels []string) *Node {
	var myTree *Node

	if len(dataSet[0].Positions) == 0 {
		label := dataSet[0].Label
		myTree = &Node{Name: label}
		return myTree
	}

	label := dataSet[0].Label
	isClassEqual := true
	for _, point := range dataSet[1:] {
		if point.Label != label {
			isClassEqual = false
			break
		}
	}
	if isClassEqual {
		myTree = &Node{Name: majorityCnt(dataSet)}
		return myTree
	}

	newLabels := make([]string, len(labels))
	copy(newLabels, labels)
	bestFeat := ChooseBestFeatureToSplit(dataSet)
	bestFeatLabel := labels[bestFeat]
	newLabels = append(newLabels[:bestFeat], newLabels[bestFeat+1:]...)
	uniqueVals := make(map[interface{}]int)
	for _, point := range dataSet {
		uniqueVals[point.Positions[bestFeat]] = 1
	}
	myTree = &Node{Name: bestFeatLabel}
	myTree.SubNodes = make(map[interface{}]*Node)
	for value := range uniqueVals {
		myTree.SubNodes[value] = CreateTree(
			SplitDataSet(dataSet, bestFeat, value), newLabels)
	}
	return myTree
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
