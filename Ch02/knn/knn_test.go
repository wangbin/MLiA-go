package knn

import (
	"fmt"
	"strings"
	"testing"
)

func TestPointStringer(t *testing.T) {
	p1 := NewPoint(12.0, 345, "A")
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
