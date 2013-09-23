package knn

import (
	"fmt"
	"math"
	//	"os"
	"sort"

//	"strconv"
)

type Point struct {
	Positions []float64
	Label     string
	distance  float64
}

func (point Point) String() string {
	return fmt.Sprintf("%v - %s", point.Positions, point.Label)
}

func (point *Point) autoNorm(minVals []float64, ranges []float64) {
	var minVal, range_ float64
	for index, position := range point.Positions {
		minVal = minVals[index]
		range_ = ranges[index]
		point.Positions[index] = (position - minVal) / range_
	}
}

func (p1 *Point) calculateDistance(p2 *Point) {
	distance := 0.0
	for index := range p1.Positions {
		distance += math.Pow(p1.Positions[index]-p2.Positions[index], 2)
	}
	p1.distance = math.Sqrt(distance)
}

func (inX *Point) Classify(group *Group, k int) {
	for _, p := range group.Points {
		p.calculateDistance(inX)
	}
	sort.Sort(group)
	classCount := make(map[string]int)
	for _, p := range group.Points[:k] {
		if _, ok := classCount[p.Label]; ok {
			classCount[p.Label] += 1
		} else {
			classCount[p.Label] = 1
		}
	}
	var result string
	maxCount := 0
	for label, count := range classCount {
		if count > maxCount {
			maxCount = count
			result = label
		}
	}
	inX.Label = result
}

type Group struct {
	Points  []*Point
	maxVals []float64
	minVals []float64
	ranges  []float64
}

func (group Group) Len() int {
	return len(group.Points)
}

func (group Group) Less(i, j int) bool {
	return group.Points[i].distance < group.Points[j].distance
}

func (group Group) Swap(i, j int) {
	group.Points[i], group.Points[j] = group.Points[j], group.Points[i]
}

func (group *Group) Append(point *Point) {
	if len(group.Points) > 0 {
		for index, val := range point.Positions {
			if val > group.maxVals[index] {
				group.maxVals[index] = val
			}
			if val < group.minVals[index] {
				group.minVals[index] = val
			}
		}
	} else {
		//group.minVals = make([]float64, len(point.Positions))
		//group.maxVals = make([]float64, len(point.Positions))
		copy(group.minVals, point.Positions)
		copy(group.maxVals, point.Positions)
	}
	group.Points = append(group.Points, point)
}

func (group *Group) AutoNorm() {
	var maxVal float64
	for index, minVal := range group.minVals {
		maxVal = group.maxVals[index]
		group.ranges[index] = maxVal - minVal
	}
	var point *Point
	for index := range group.Points {
		point = group.Points[index]
		point.autoNorm(group.minVals, group.ranges)
	}
}
