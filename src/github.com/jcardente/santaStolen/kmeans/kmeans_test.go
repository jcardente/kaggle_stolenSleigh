package kmeans

import (
	"math"
	"testing"
	"github.com/jcardente/santaStolen/sqt"
)


func TestDistance(t *testing.T) {
	n:= sqt.Node{1.0,1.0,0,0.0}
	c:= Centroid{1,1.0,1.0,0.0,[]*sqt.Node{}}

	d:= distance(c, &n)
	if (d != 0.0) {
		t.Fatal("Incorrect distance, expected 0 but got ", d)
	}

	c.X = 0
	c.Y = 0

	d = distance(c, &n)
	if (d != math.Sqrt(2)) {
		t.Fatal("Incorrect distance, expected sqrt(2) but got ", d)
	}
}

func TestInitCent(t *testing.T) {
	nodes := []*sqt.Node{
		&sqt.Node{1,1,1,0},
		&sqt.Node{2,2,2,0},
		&sqt.Node{3,3,3,0},
		&sqt.Node{4,4,4,0},
		&sqt.Node{5,5,5,0},
		&sqt.Node{6,6,6,0},
		
	}

	centroids := initCentroids(nodes, 3)

	if (len(*centroids) != 3) {
		t.Fatal("Got incorrect number of centroids ", len(*centroids))
	}

}

func TestReComputeCentroid(t *testing.T) {

	nodes := []*sqt.Node{
		&sqt.Node{1,1,1,1},
		&sqt.Node{2,2,2,1},
		&sqt.Node{3,3,3,1},
		&sqt.Node{4,4,4,1},
		&sqt.Node{5,5,5,1},
	}

	c := Centroid{1,0,0,0,nodes}
	recomputeCentroid(&c)
	if (c.X != 3.0) || (c.Y != 3.0) || (c.Weight != 0) {
		t.Fatal("Incorrect centroid ", c)
	}

	nodes = []*sqt.Node{
		&sqt.Node{64.758426,1.607777,1,0},
		&sqt.Node{3.199829,36.810519,2,0},
		&sqt.Node{80.587278,4.723220,3,0},
		&sqt.Node{59.478304,32.165355,4,0},
		&sqt.Node{18.782064,62.800670,5,0},
	}

	c.Nodes = nodes
	recomputeCentroid(&c)
	if (math.Abs(c.X - 45.36) > 0.01) || (math.Abs(c.Y-27.62) > 0.01) {
		t.Fatal("Incorrect centroid ", c)
	}
}


func TestClosestCentroid(t *testing.T) {

	node := sqt.Node{1,1,1,1}
	centroids := []Centroid{
		Centroid{0,0.5,0.5,0,[]*sqt.Node{}},
		Centroid{1,2,2,0,[]*sqt.Node{}},
		Centroid{2,3,3,0,[]*sqt.Node{}},
	}

	cid := closestCentroid(&node, math.Inf(0), &centroids)
	if (cid != 0) {
		t.Fatal("Incorrect cloest centroid expected 0, got ", cid)
	}
	
	centroids = []Centroid{
		Centroid{0,0.5,0.5,0,[]*sqt.Node{}},
		Centroid{1,2,2,0,[]*sqt.Node{}},
		Centroid{2,3,3,0,[]*sqt.Node{}},
		Centroid{3,1.25,1.25,0,[]*sqt.Node{}},		
	}

	cid = closestCentroid(&node, math.Inf(0), &centroids)
	if (cid != 3) {
		t.Fatal("Incorrect cloest centroid ", cid)
	}

	centroids = []Centroid{
		Centroid{0,0.5,0.5,5,[]*sqt.Node{}},
		Centroid{1,2,2,0,[]*sqt.Node{}},
		Centroid{2,3,3,5,[]*sqt.Node{}},
		Centroid{3,1.25,1.25,5,[]*sqt.Node{}},		
	}
	
	cid = closestCentroid(&node, 5, &centroids)
	if (cid != 1) {
		t.Fatal("Incorrect cloest centroid ", cid)
	}


	centroids = []Centroid{
		Centroid{0,0.5,0.5,5,[]*sqt.Node{}},
		Centroid{1,2,2,5,[]*sqt.Node{}},
		Centroid{2,3,3,5,[]*sqt.Node{}},
		Centroid{3,1.25,1.25,5,[]*sqt.Node{}},		
	}
	
	cid = closestCentroid(&node, 5, &centroids)
	if (cid != -1) {
		t.Fatal("Incorrect cloest centroid ", cid)
	}	
}



func TestCluster(t *testing.T) {

	nodes := []*sqt.Node{
		&sqt.Node{1.0,1.0,1,1},
		&sqt.Node{1.1,1.1,2,1},
		&sqt.Node{4.1,4.0,3,1},
		&sqt.Node{4.5,4.3,4,1},
		&sqt.Node{4.3,4.5,5,1},
	}

	clusts := Cluster(nodes, 2, math.Inf(0))
	if len(clusts) != 2 {
		t.Fatal("Incorrect number of clusters")
	}

	nodes = []*sqt.Node{
		&sqt.Node{1.0,1.0,1,1},
		&sqt.Node{1.1,1.1,2,1},
		&sqt.Node{3.0,3.0,3,1},
		&sqt.Node{4.5,4.3,4,1},
		&sqt.Node{4.3,4.5,5,1},
	}

	clusts = Cluster(nodes, 3, 2)
	if len(clusts) != 3 {
		t.Fatal("Incorrect number of clusters")
	}

	 for cid, c := range clusts {
	 	for _, n := range c {
	 		t.Log(cid, ": ",n)
	 	}
	 }
	

	t.Log("-----")
	clusts = Cluster(nodes, 2, 2)
	 for cid, c := range clusts {
	 	for _, n := range c {
	 		t.Log(cid, ": ",n)
	 	}
	 }

	 if len(clusts) != 3 {
	 	t.Fatal("Incorrect number of clusters ", len(clusts))
	 }
	
}

