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
	FlyerMiles     float64
	GameTimes      float64
	IcecreamLiters float64
	DatingLabel    string
}

type Daters struct {
	daters            []*Dater
	minFlyerMiles     float64
	maxFlyerMiles     float64
	minGameTimes      float64
	maxGameTimes      float64
	minIcecreamLiters float64
	maxIcecreamLiters float64
	ranges            []float64
	minVals           []float64
}

func main() {
	//classify0()
	daters := file2matrix("datingTestSet.txt")
	daters.autoNorm()
	// for _, dater := range daters[:20] {
	// 	fmt.Println(dater)
	// }
	// daters := new(Daters)
	// fmt.Println(len(daters.daters))
	// dd1 := &Dater{12.0, 11.0, 10.0, "A"}
	// dd2 := &Dater{22.0, 21.0, 20.0, "C"}
	// dd3 := &Dater{17.0, 16.0, 15.0, "C"}
	// daters.Append(dd1)
	// daters.Append(dd2)
	// daters.Append(dd3)
	// daters.autoNorm()
	fmt.Printf("ranges: %v\n", daters.ranges)
	fmt.Printf("minVals: %v\n", daters.minVals)
	// for _, dater := range daters.daters {
	// 	fmt.Println(dater)
	// }
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
	return fmt.Sprintf("[%f - %f - %f - %s]", dater.FlyerMiles, dater.GameTimes,
		dater.IcecreamLiters, dater.DatingLabel)
}

func classify0() {
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

func file2matrix(filename string) *Daters {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	daters := new(Daters)
	for {
		var v1, v2, v3, datingLabel string
		_, err := fmt.Fscanln(file, &v1, &v2, &v3, &datingLabel)
		if err != nil {
			break
		}

		flyerMiles, err := strconv.ParseFloat(v1, 64)
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
		daters.Append(dater)
	}
	return daters
}

func (daters *Daters) Append(dater *Dater) {
	if len(daters.daters) > 0 {
		daters.minFlyerMiles = math.Min(daters.minFlyerMiles, dater.FlyerMiles)
		daters.maxFlyerMiles = math.Max(daters.maxFlyerMiles, dater.FlyerMiles)
		daters.minGameTimes = math.Min(daters.minGameTimes, dater.GameTimes)
		daters.maxGameTimes = math.Max(daters.maxGameTimes, dater.GameTimes)
		daters.minIcecreamLiters = math.Min(daters.minIcecreamLiters,
			dater.IcecreamLiters)
		daters.maxIcecreamLiters = math.Max(daters.maxIcecreamLiters,
			dater.IcecreamLiters)
	} else {
		daters.minFlyerMiles = dater.FlyerMiles
		daters.maxFlyerMiles = dater.FlyerMiles
		daters.minGameTimes = dater.GameTimes
		daters.maxGameTimes = dater.GameTimes
		daters.minIcecreamLiters = dater.IcecreamLiters
		daters.maxIcecreamLiters = dater.IcecreamLiters
	}
	daters.daters = append(daters.daters, dater)
}

func (daters *Daters) autoNorm() {
	var dater *Dater
	flyerMilesRange := daters.maxFlyerMiles - daters.minFlyerMiles
	daters.ranges = append(daters.ranges, flyerMilesRange)
	daters.minVals = append(daters.minVals, daters.minFlyerMiles)
	gameTimesRange := daters.maxGameTimes - daters.minGameTimes
	daters.ranges = append(daters.ranges, gameTimesRange)
	daters.minVals = append(daters.minVals, daters.minGameTimes)
	icecreamLitersRange := daters.maxIcecreamLiters - daters.minIcecreamLiters
	daters.ranges = append(daters.ranges, icecreamLitersRange)
	daters.minVals = append(daters.minVals, daters.minIcecreamLiters)
	for index := range daters.daters {
		dater = daters.daters[index]
		dater.FlyerMiles = (dater.FlyerMiles - daters.minFlyerMiles) / flyerMilesRange
		dater.GameTimes = (dater.GameTimes - daters.minGameTimes) / gameTimesRange
		dater.IcecreamLiters = (dater.IcecreamLiters - daters.minIcecreamLiters) / icecreamLitersRange
	}
}
