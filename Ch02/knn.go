package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

type Point struct {
	Position []float64
	Label    byte
	Distance float64
}

type Group []*Point

type Dater struct {
	FlyerMiles     int
	GameTimes      float64
	IcecreamLiters float64
	DatingLabel    string
}

type Daters []*Dater

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

	daters := file2matrix("datingTestSet.txt")
	for _, dater := range daters[:20] {
		fmt.Println(dater)
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

func (dater Dater) String() string {
	return fmt.Sprintf("[%d - %f - %f - %s]", dater.FlyerMiles, dater.GameTimes,
		dater.IcecreamLiters, dater.DatingLabel)
}

func file2matrix(filename string) (daters Daters) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for {
		var v1, v2, v3, datingLabel string
		_, err := fmt.Fscanln(file, &v1, &v2, &v3, &datingLabel)
		if err != nil {
			break
		}

		flyerMiles, err := strconv.Atoi(v1)
		if err != nil {
			continue
		}
		gameTimes, err := strconv.ParseFloat(v2, 64)
		if err != nil {
			continue
		}
		icecreamLiters, err := strconv.ParseFloat(v3, 64)
		if err != nil {
			continue
		}

		dater := &Dater{flyerMiles, gameTimes, icecreamLiters, datingLabel}
		daters = append(daters, dater)
	}
	return
}
