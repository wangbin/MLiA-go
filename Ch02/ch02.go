package main

import (
	"./knn"
	"fmt"
)

func main() {
	classify0()
}

func classify0() {
	dataSet := new(knn.Group)
	dataSet.Append(&knn.Point{Positions: []float64{1.0, 1.1}, Label: "A"})
	dataSet.Append(&knn.Point{Positions: []float64{1.0, 1.0}, Label: "A"})
	dataSet.Append(&knn.Point{Positions: []float64{0, 0}, Label: "B"})
	dataSet.Append(&knn.Point{Positions: []float64{0, 0.1}, Label: "B"})

	inXs := new(knn.Group)
	inXs.Append(&knn.Point{Positions: []float64{0, 0}})
	inXs.Append(&knn.Point{Positions: []float64{0, 0.5}})
	inXs.Append(&knn.Point{Positions: []float64{0, 1}})
	inXs.Append(&knn.Point{Positions: []float64{0.5, 0.5}})
	inXs.Append(&knn.Point{Positions: []float64{1, 0}})
	inXs.Append(&knn.Point{Positions: []float64{1, 0.5}})
	inXs.Append(&knn.Point{Positions: []float64{1, 1}})

	var inX *knn.Point

	for index := range inXs.Points {
		inX = inXs.Points[index]
		inX.Classify(dataSet, 3)
		fmt.Println(inX)
	}

}
