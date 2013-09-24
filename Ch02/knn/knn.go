package knn

import (
	"fmt"
	"math"
	"sort"
)

type Point struct {
	Positions []float64
	Label     string
	distance  float64
}

func (point Point) String() string {
	if len(point.Label) > 0 {
		return fmt.Sprintf("%v - %s", point.Positions, point.Label)
	} else {
		return fmt.Sprintf("%v", point.Positions)
	}
}

func (point *Point) autoNorm(minVals []float64, ranges []float64) {
	var minVal, range_ float64
	for index, position := range point.Positions {
		minVal = minVals[index]
		range_ = ranges[index]
		point.Positions[index] = (position - minVal) / range_
	}
}

func (p1 *Point) Distance(p2 *Point) (distance float64) {
	for index := range p1.Positions {
		distance += math.Pow(p1.Positions[index]-p2.Positions[index], 2)
	}
	distance = math.Sqrt(distance)
	return
}

func NewPoint(params ...interface{}) *Point {
	point := new(Point)
	for _, param := range params {
		switch param.(type) {
		case string:
			point.Label = param.(string)
		case int:
			point.Positions = append(point.Positions, float64(param.(int)))
		default:
			point.Positions = append(point.Positions, param.(float64))
		}
	}
	return point
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
		length := len(point.Positions)
		group.minVals = make([]float64, length)
		group.maxVals = make([]float64, length)
		group.ranges = make([]float64, length)
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

func (group *Group) Classify(point *Point, k int) string {
	needAutoNorm := false
	for index, position := range point.Positions {
		minVal, maxVal := group.minVals[index], group.maxVals[index]
		if position > maxVal || position < minVal {
			needAutoNorm = true
			break
		}
	}

	if needAutoNorm {
		point.autoNorm(group.minVals, group.ranges)
	}

	for _, p := range group.Points {
		p.distance = p.Distance(point)
	}
	sort.Sort(group)
	classCount := make(map[string]int)
	distanceCount := make(map[string]float64)
	for _, p := range group.Points[:k] {
		if _, ok := classCount[p.Label]; ok {
			classCount[p.Label] += 1
		} else {
			classCount[p.Label] = 1
		}
		if distance, ok := distanceCount[p.Label]; ok {
			if distance > p.distance {
				distanceCount[p.Label] = p.distance
			}
		} else {
			distanceCount[p.Label] = p.distance
		}
	}
	var result string
	maxCount := 0
	for label, count := range classCount {
		if count > maxCount {
			maxCount = count
			result = label
		}
		if count == maxCount {
			if distanceCount[result] > distanceCount[label] {
				result = label
			}
		}
	}
	return result
}

func NewGroup(points ...*Point) *Group {
	group := new(Group)
	for index := range points {
		group.Append(points[index])
	}
	return group
}
