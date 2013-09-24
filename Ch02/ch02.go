package main

import (
	"./knn"
	"fmt"
	"os"
	"strconv"
)

func main() {
	//classify0()
	//datingClassTest()
	classifyPerson()
}

func classify0() {
	dataSet := knn.NewGroup(knn.NewPoint(1.0, 1.1, "A"), knn.NewPoint(1.0, 1.0, "A"),
		knn.NewPoint(0, 0, "B"), knn.NewPoint(0, 0.1, "B"))

	var inX *knn.Point
	inX = knn.NewPoint(1, 0)
	fmt.Println(dataSet.Classify(inX, 3))
}

func file2matrix(filename string) *knn.Group {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var p1Str, p2Str, p3Str, label string
	group := new(knn.Group)
	for {
		_, err := fmt.Fscanln(file, &p1Str, &p2Str, &p3Str, &label)
		if err != nil {
			break
		}
		p1, err := strconv.ParseFloat(p1Str, 64)
		if err != nil {
			continue
		}

		p2, err := strconv.ParseFloat(p2Str, 64)
		if err != nil {
			continue
		}

		p3, err := strconv.ParseFloat(p3Str, 64)
		if err != nil {
			continue
		}
		point := knn.NewPoint(p1, p2, p3, label)
		group.Append(point)
	}
	return group
}

func datingClassTest() {
	hoRatio := 0.50
	errorCount := 0.0
	dataSet := file2matrix("datingTestSet2.txt")
	dataSet.AutoNorm()
	m := len(dataSet.Points)
	numTestVecs := int(float64(m) * hoRatio)
	testData := make([]*knn.Point, numTestVecs)
	copy(testData, dataSet.Points[:numTestVecs])
	dataSet.Points = dataSet.Points[numTestVecs:]
	var result string
	for index, point := range testData {
		result = dataSet.Classify(point, 3)
		if result != point.Label {
			fmt.Printf("the classifier came back with: %s, the real answer is: %s\n", result,
				point.Label)
			fmt.Println(index)
			errorCount += 1.0
		}
	}
	fmt.Printf("the total error rate is: %f\n", errorCount/float64(numTestVecs))
	fmt.Println(errorCount)
}

func classifyPerson() {
	dataSet := file2matrix("datingTestSet.txt")
	dataSet.AutoNorm()
	var percentTatsStr, ffMilesStr, iceCreamStr string

	fmt.Println("percentage of time spent playing video games?")
	fmt.Scanln(&percentTatsStr)
	percentTats, err := strconv.ParseFloat(percentTatsStr, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println("frequent flier miles earned per year?")
	fmt.Scanln(&ffMilesStr)
	ffMiles, err := strconv.ParseFloat(ffMilesStr, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println("liters of ice cream consumed per year?")
	fmt.Scanln(&iceCreamStr)
	iceCream, err := strconv.ParseFloat(iceCreamStr, 64)
	if err != nil {
		panic(err)
	}
	p := knn.NewPoint(ffMiles, percentTats, iceCream)
	result := dataSet.Classify(p, 3)
	fmt.Printf("You will probably like this person: %s\n", result)

}
