package knn

import "testing"

func TestPointdistance(t *testing.T) {
	p1 := &Point{Positions: []float64{0.5, 6}}
	p2 := &Point{Positions: []float64{4.5, 3}}
	p1.calculateDistance(p2)
	if p1.distance != 5 {
		t.Errorf("Distance between Point(%v) & Point(%v) should be 3, not %f", p1, p2, p1.distance)
	}
}
