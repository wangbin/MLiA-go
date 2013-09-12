package main

import (
	"fmt"
	"math"
	"sort"
)

type Point struct {
	Position []float64
	Label    byte
	Distance float64
}

type Group []*Point

func main() {
	dataSet := Group{
		&Point{Position: []float64{1.0, 1.1}, Label: 'A'},
		&Point{Position: []float64{1.0, 1.0}, Label: 'A'},
		&Point{Position: []float64{0, 0}, Label: 'B'},
		&Point{Position: []float64{0, 0.1}, Label: 'B'},
	}

	inXs := Group{
		&Point{Position: []float64{0, 0}},
		&Point{Position: []float64{0, 0.5}},
		&Point{Position: []float64{0, 1}},
		&Point{Position: []float64{0.5, 0.5}},
		&Point{Position: []float64{1, 0}},
		&Point{Position: []float64{1, 0.5}},
		&Point{Position: []float64{1, 1}}}

	for _, inX := range inXs {
		inX.classify(dataSet, 3)
		fmt.Println(inX)
	}
}

func (point Point) String() string {
	return fmt.Sprintf("%v - %c", point.Position, point.Label)
}

func (p1 *Point) calculateDistance(p2 *Point) {
	distance := 0.0
	for index := range p1.Position {
		distance += math.Pow(p1.Position[index]-p2.Position[index], 2)
	}
	p1.Distance = math.Sqrt(distance)
}

func (inX *Point) classify(group Group, k int) {
	for _, p := range group {
		p.calculateDistance(inX)
	}
	sort.Sort(group)
	classCount := make(map[byte]int)
	for _, p := range group[:k] {
		if _, ok := classCount[p.Label]; ok {
			classCount[p.Label] += 1
		} else {
			classCount[p.Label] = 1
		}
	}
	var result byte
	maxCount := 0
	for label, count := range classCount {
		if count > maxCount {
			maxCount = count
			result = label
		}
	}
	inX.Label = result
}

func (group Group) Len() int {
	return len(group)
}

func (group Group) Less(i, j int) bool {
	return group[i].Distance < group[j].Distance
}

func (group Group) Swap(i, j int) {
	group[i], group[j] = group[j], group[i]
}
