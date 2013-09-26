package trees

import (
	"../../Ch02/knn"
	"fmt"
	"math"
	"testing"
)

func createDataSet() (*knn.Group, []string) {
	group := knn.NewGroup(
		knn.NewPoint(1, 1, "yes"),
		knn.NewPoint(1, 1, "yes"),
		knn.NewPoint(1, 0, "no"),
		knn.NewPoint(0, 1, "no"),
		knn.NewPoint(0, 1, "no"))
	labels := []string{"no surfacing", "flippers"}

	return group, labels
}

func TestCalcShannonEnt(t *testing.T) {
	dataSet, _ := createDataSet()
	if math.Abs(CalcShannonEnt(dataSet)-0.9709505944546686) > 1e-10 {
		t.FailNow()
	}

	dataSet.Points[0].Label = "maybe"
	if math.Abs(CalcShannonEnt(dataSet)-1.3709505944546687) > 1e-10 {
		t.FailNow()
	}
}

func TestSplitDataSet(t *testing.T) {
	dataSet, _ := createDataSet()
	newDataSet := SplitDataSet(dataSet, 0, 1)
	if len(newDataSet.Points) != 3 {
		t.FailNow()
	}
	for _, point := range newDataSet.Points {
		if len(point.Positions) != 1 {
			t.FailNow()
		}
	}
	newDataSet = SplitDataSet(dataSet, 0, 0)
	if len(newDataSet.Points) != 2 {
		t.FailNow()
	}

	for _, point := range newDataSet.Points {
		if len(point.Positions) != 1 {
			t.FailNow()
		}
	}
}

func TestChooseBestFeatureToSplit(t *testing.T) {
	dataSet, _ := createDataSet()
	if ChooseBestFeatureToSplit(dataSet) != 0 {
		t.FailNow()
	}
}

func TestClassify(t *testing.T) {
	dataSet, labels := createDataSet()
	myTree := CreateTree(dataSet, labels)
	if myTree.Classify(labels, knn.NewPoint(1, 0)) != "no" {
		t.FailNow()
	}
	if myTree.Classify(labels, knn.NewPoint(1, 1)) != "yes" {
		t.FailNow()
	}

}

func TestTreeStoreRetrice(t *testing.T) {
	dataSet, labels := createDataSet()
	myTree := CreateTree(dataSet, labels)
	StoreTree(myTree, "/tmp/myTree.dat")
	node, _ := GrabTree("/tmp/myTree.dat")
	if fmt.Sprint(myTree) != fmt.Sprint(node) {
		t.FailNow()
	}
}
