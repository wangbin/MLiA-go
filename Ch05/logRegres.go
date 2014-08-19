package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

func loadDataSet() (dataMat [][]float64, labelMat []float64) {
	file, err := os.Open("testSet.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for {
		var x1Str, x2Str, labelStr string
		_, err := fmt.Fscanln(file, &x1Str, &x2Str, &labelStr)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		x1, err := strconv.ParseFloat(x1Str, 64)
		if err != nil {
			panic(err)
		}

		x2, err := strconv.ParseFloat(x2Str, 64)
		if err != nil {
			panic(err)
		}

		label, err := strconv.ParseFloat(labelStr, 64)
		if err != nil {
			panic(err)
		}

		dataMat = append(dataMat, []float64{1.0, x1, x2})
		labelMat = append(labelMat, label)
	}
	return
}

func sigmoid(x float64) float64 {
	return 1.0 / (1 + math.Exp(0-x))
}

func gradAscent(dataMat [][]float64, labelMat []float64) (weights []float64) {
	m := len(dataMat)
	n := len(dataMat[0])
	alpha := 0.001
	maxCycle := 500
	for i := 0; i < n; i++ {
		weights = append(weights, 1.0)
	}
	for i := 0; i < maxCycle; i++ {
		for j := 0; j < m; j++ {
			data := dataMat[j]
			label := labelMat[j]
			for k := 0; k < n; k++ {
				h := sigmoid(data[k] * weights[k])
				error := label - h
				weights[k] += alpha * error
			}
		}
	}
	return
}

func main() {
	dataMat, labelMat := loadDataSet()
	fmt.Println(dataMat)
	fmt.Println(labelMat)
	fmt.Println(gradAscent(dataMat, labelMat))
}
