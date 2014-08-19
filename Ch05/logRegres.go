package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"strconv"
)

type Array []float64

func (this Array) Add(a Array) Array {
	length := len(this)
	result := make(Array, length)
	for i, _ := range this {
		result[i] = this[i] + a[i]
	}
	return result
}

func (this Array) Minus(a Array) Array {
	length := len(this)
	result := make(Array, length)
	for i, _ := range this {
		result[i] = this[i] - a[i]
	}
	return result
}

func (this Array) Sum() float64 {
	sum := 0.0
	for _, x := range this {
		sum += x
	}
	return sum
}

func (this Array) Multiple(a float64) Array {
	length := len(this)
	result := make(Array, length)
	for i, _ := range this {
		result[i] = this[i] * a
	}
	return result
}

func (this Array) MultipleArray(a Array) Array {
	length := len(this)
	result := make(Array, length)
	for i, _ := range this {
		result[i] = this[i] * a[i]
	}
	return result
}

func (this Array) Sigmoid() Array {
	s := make(Array, 0)
	for _, x := range this {
		s = append(s, sigmoid(x))
	}
	return s
}

func sigmoid(x float64) float64 {
	return 1.0 / (1 + math.Exp(0-x))
}

type Matrix []Array

func (this Matrix) Shape() (int, int) {
	return len(this), len(this[0])
}

func (this Matrix) Multiple(weights Array) (result Array) {
	m, n := this.Shape()
	for i := 0; i < m; i++ {
		sum := 0.0
		for j := 0; j < n; j++ {
			sum += this[i][j] * weights[j]
		}
		result = append(result, sum)
	}
	return
}

func (this Matrix) Transpose() Matrix {
	m, n := this.Shape()
	result := make(Matrix, n)
	for i := 0; i < n; i++ {
		result[i] = make(Array, m)
		for j := 0; j < m; j++ {
			result[i][j] = this[j][i]
		}
	}
	return result
}

func loadDataSet() (dataMat Matrix, labelMat Array) {
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

		dataMat = append(dataMat, Array{1.0, x1, x2})
		labelMat = append(labelMat, label)
	}
	return
}

func gradAscent(dataMat Matrix, labelMat Array) (weights Array) {
	_, n := dataMat.Shape()
	t := dataMat.Transpose()
	alpha := 0.001
	maxCycle := 500
	for i := 0; i < n; i++ {
		weights = append(weights, 1.0)
	}
	for i := 0; i < maxCycle; i++ {
		h := dataMat.Multiple(weights).Sigmoid()
		error := labelMat.Minus(h)
		weights = weights.Add(t.Multiple(error).Multiple(alpha))
	}
	return
}

func stocGradAscent0(dataMat Matrix, labelMat Array) (weights Array) {
	m, n := dataMat.Shape()
	alpha := 0.01
	for i := 0; i < n; i++ {
		weights = append(weights, 1.0)
	}
	for i := 0; i < m; i++ {
		h := sigmoid(dataMat[i].MultipleArray(weights).Sum())
		error := labelMat[i] - h
		weights = weights.Add(dataMat[i].Multiple(error).Multiple(alpha))
	}
	return
}

func stocGradAscent1(dataMat Matrix, labelMat Array, numIter int) (weights Array) {
	m, n := dataMat.Shape()
	for i := 0; i < n; i++ {
		weights = append(weights, 1.0)
	}
	for j := 0; j < numIter; j++ {
		dataIndex := make([]int, m)
		for i := 0; i < m; i++ {
			dataIndex[i] = i
		}
		for i := 0; i < m; i++ {
			alpha := 4.0/(1.0+float64(j)+float64(i)) + 0.001
			randIndex64, _ := rand.Int(rand.Reader, big.NewInt(int64(len(dataIndex))))
			randIndex, _ := strconv.Atoi(randIndex64.String())
			h := sigmoid(dataMat[randIndex].MultipleArray(weights).Sum())
			error := labelMat[randIndex] - h
			weights = weights.Add(dataMat[randIndex].Multiple(error).Multiple(alpha))
			dataIndex = append(dataIndex[:randIndex], dataIndex[randIndex+1:]...)
		}
	}
	return
}

func main() {
	dataMat, labelMat := loadDataSet()
	fmt.Println(gradAscent(dataMat, labelMat))
	fmt.Println(stocGradAscent0(dataMat, labelMat))
	fmt.Println(stocGradAscent1(dataMat, labelMat, 150))
	fmt.Println(stocGradAscent1(dataMat, labelMat, 500))
}
