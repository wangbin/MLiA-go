package knn

import (
	"fmt"
	"strings"
	"testing"
)

func TestPointStringer(t *testing.T) {
	p1 := NewPoint(12.0, "A", 345)
	p1Str := fmt.Sprint(p1)
	if !strings.HasSuffix(p1Str, " - A") {
		t.Errorf("Point(%v) Stringer should include Label A.", p1)
	}

	p2 := NewPoint(2, 10, 5)
	if strings.Contains(fmt.Sprint(p2), "-") {
		t.Errorf("Point(%v) Stringer should not include Label.", p2)
	}
}

func TestPointDistance(t *testing.T) {
	p1 := NewPoint(0.5, 6)
	p2 := NewPoint(4.5, 3)
	distance := p1.Distance(p2)
	if distance != 5 {
		t.Errorf("Distance between Point(%v) & Point(%v) should be 5, not %f",
			p1, p2, distance)
	}
}

func TestGroupautoNorm(t *testing.T) {
	p1 := NewPoint(1000, 150, 2.5)
	p2 := NewPoint(950, 100, 1)
	p3 := NewPoint(500, 50, 0.5)
	group := NewGroup(p1, p2, p3)
	group.AutoNorm()

	p1 = group.Points[0]
	for _, position := range p1.Positions {
		if position != 1 {
			t.Error(position)
		}
	}
	p2 = group.Points[1]
	p2Positions := []float64{0.9, 0.5, 0.25}
	for index, position := range p2.Positions {
		if p2Positions[index] != position {
			t.Error(position)
		}
	}
	p3 = group.Points[2]
	for _, position := range p3.Positions {
		if position != 0 {
			t.Error(position)
		}
	}
}

func TestClassify(t *testing.T) {
	p1 := NewPoint(1.0, 1.1, "A")
	p2 := NewPoint(1.0, 1.0, "A")
	p3 := NewPoint(0, 0, "B")
	p4 := NewPoint(0, 0.1, "B")
	dataSet := NewGroup(p1, p2, p3, p4)
	testData := NewGroup(NewPoint(0, 0), NewPoint(0, 0.5), NewPoint(0, 1),
		NewPoint(0.5, 0.5), NewPoint(1, 0), NewPoint(1, 0.5))
	labels := []string{"B", "B", "B", "B", "B", "A"}
	for index, point := range testData.Points {
		label, _ := dataSet.Classify(point, 3)
		if label != labels[index] {
			t.Errorf("Label of Point(%v) should be %s, not %s.", point, labels[index],
				label)
		}
	}
}
