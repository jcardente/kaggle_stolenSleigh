package kmeans

import (
	"fmt"
	"math"
	"math/rand"
	"github.com/jcardente/santaStolen/sqt"
)

type Centroid struct {
	X,Y float64
	Weight float64
}


func distance(c Centroid, n *sqt.Node) float64 {
	return math.Sqrt(math.Pow(c.X - n.X, 2) + math.Pow(c.Y - n.Y,2))
}

func initCentroids(nodes []*sqt.Node, k int) *[]Centroid {

	// Pick initial centroids
	centroids := make([]Centroid, k)
	for c, i := range (rand.Perm(len(nodes)))[:k] {
		centroids[c] = Centroid{nodes[i].X, nodes[i].Y,0}
	}

	return &centroids
}

func recomputeCentroid(nodes []*sqt.Node) Centroid {
	centroid := Centroid{0,0,0}
	for _, n := range nodes {
		centroid.X += n.X
		centroid.Y += n.Y
	}

	count := float64(len(nodes))
	centroid.X = centroid.X/count
	centroid.Y = centroid.Y/count
	centroid.Weight = 0
	
	return centroid
}

func closestCentroid(n *sqt.Node, maxWeight float64, centroids *[]Centroid) int {
	closest := -1
	dist    := math.Inf(0)
	for i, c := range *centroids {
		d:= distance(c, n)
		if (d < dist) && ((n.Weight + c.Weight) <= maxWeight) {
			dist = d
			closest = i
		}
	}

	return closest
}


func Cluster(nodes []*sqt.Node, k int, maxWeight float64) map[int][]*sqt.Node {

	// Store each cluster as a list of node pointers
	clusts:= map[int][]*sqt.Node{}

	// Initialize centroids using randomly selected
	// nodes
	centroids := initCentroids(nodes, k)

	// Initialize memeberships clusters
	members := map[*sqt.Node]int{}
	for _, n := range nodes {
		members[n] = -1
	}
	
	// Loop until memberships don't change
	converged := false
	for (!converged) {
		converged = true

		for i:=0; i < len(clusts); i++ {
			clusts[i] = []*sqt.Node{}
		}

		for _, n := range nodes {
			cid := closestCentroid(n, maxWeight, centroids)
			if (cid >=0 ) {
				(*centroids)[cid].Weight += n.Weight
				clusts[cid] = append(clusts[cid], n)
			} else {
				// Uh oh, didn't find a cluster to fit into
				// Make a new one
				_c := append(*centroids, Centroid{n.X, n.Y, n.Weight})
				centroids = &_c
				cid = len(*centroids)-1
				clusts[cid] = []*sqt.Node{n}
			}
			
			if (cid != members[n]) {
				members[n]  = cid
				converged   = false
			}
		}

		fmt.Println("LCents: ", len(*centroids)," LClusts: ", len(clusts))
		if (!converged) {
			// Update cluster centroids
			for cid, clust := range clusts {
				(*centroids)[cid] = recomputeCentroid(clust)
			}
		}
	}
	
	return clusts
}
