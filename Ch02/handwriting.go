package main

import (
	"./knn"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	TrainingDigitsDir = "trainingDigits"
	TestDigitsDir     = "testDigits"
)

func main() {
	handwritingClassTest()
}

func img2vector(filename string) []float64 {
	returnVect := make([]float64, 1024)

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	index := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		for i := range line {
			number, err := strconv.ParseFloat(line[i:i+1], 64)
			if err != nil {
				continue
			}
			returnVect[index] = number
			index += 1
		}
	}

	return returnVect
}

func createDataSet() *knn.Group {
	trainingFileList, err := ioutil.ReadDir(TrainingDigitsDir)
	if err != nil {
		panic(err)
	}

	dataSet := knn.NewGroup()
	var fileName, label string
	var vector []float64
	var point *knn.Point
	for _, trainingFile := range trainingFileList {
		fileName = trainingFile.Name()
		label = strings.Split(fileName, "_")[0]
		vector = img2vector(fmt.Sprintf("%s/%s", TrainingDigitsDir, fileName))
		point = &knn.Point{Positions: vector, Label: label}
		dataSet.Append(point)
	}
	return dataSet
}

func handwritingClassTest() {
	dataSet := createDataSet()
	testFileList, err := ioutil.ReadDir(TestDigitsDir)
	if err != nil {
		panic(err)
	}
	mTest := len(testFileList)
	errorCount := 0.0
	var fileName, label string
	var vector []float64
	var point *knn.Point
	for _, testFile := range testFileList {
		fileName = testFile.Name()
		label = strings.Split(fileName, "_")[0]
		vector = img2vector(fmt.Sprintf("%s/%s", TestDigitsDir, fileName))
		point = &knn.Point{Positions: vector, Label: label}
		result, classCount := dataSet.Classify(point, 3)
		if result != label {
			errorCount += 1.0
			fmt.Printf("the classifier came back with: %s, the real answer is: %s.\n",
				result, label)
			fmt.Printf("file name: %s, class count: %v\n", fileName, classCount)
		}
	}
	fmt.Printf("the total number of errors is: %f\n", errorCount)
	fmt.Printf("the total error rate is: %f\n", errorCount/float64(mTest))
}
