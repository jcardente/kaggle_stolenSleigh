package kmeans

import (
	"fmt"
	"os"
	"math"
	"math/rand"
	"github.com/jcardente/santaStolen/sqt"
)

type Centroid struct {
	Id     int
	X,Y    float64
	Weight float64
	Nodes  []*sqt.Node
}


func distance(c Centroid, n *sqt.Node) float64 {
	return math.Sqrt(math.Pow(c.X - n.X, 2) + math.Pow(c.Y - n.Y,2))
}

func initCentroids(nodes []*sqt.Node, k int) *[]Centroid {
	centroids := make([]Centroid, k)
	for c, i := range (rand.Perm(len(nodes)))[:k] {
		centroids[c] = Centroid{c, nodes[i].X, nodes[i].Y,0, []*sqt.Node{}}
	}

	return &centroids
}


func centroidWss(c *Centroid) float64 {
	wss := 0.0
	for _, n := range c.Nodes {
		wss += distance(*c, n)
	}

	return wss/float64(len(c.Nodes))
}

func recomputeCentroid(centroid *Centroid) {

	newX := 0.0
	newY := 0.0
	for _, n := range centroid.Nodes {
		newX += n.X
		newY += n.Y
	}

	count := float64(len(centroid.Nodes))
	newX = newX/count
	newY = newY/count

	centroid.X = newX
	centroid.Y = newY

	centroid.Weight = 0
	centroid.Nodes  = []*sqt.Node{}
}

func closestCentroid(n *sqt.Node, maxWeight float64, centroids *[]Centroid) int {
	closest := -1
	dist    := math.Inf(0)
	for _, c := range *centroids {
		d:= distance(c, n)
		if (d < dist) && ((n.Weight + c.Weight) <= maxWeight) {
			dist = d
			closest = c.Id
		}
	}

	return closest
}


func Cluster(nodes []*sqt.Node, k int, maxWeight float64) [][]*sqt.Node {

	centroids := initCentroids(nodes, k)

	members := map[*sqt.Node]int{}
	for _, n := range nodes {
		members[n] = -1
	}
	
	// Loop until cluster centroids converge
	converged := false
	tryCount  := 0
	deltaLast := math.Inf(0)
	for (converged == false) {
		tryCount++
		changeCount := 0.0
		for _, n := range nodes {
			cid := closestCentroid(n, maxWeight, centroids)
			
			if (cid < 0 ) {
				// Didn't find a cluster to fit into,
				// make a new one				
				cid = len(*centroids)
				_c := append(*centroids, Centroid{cid,n.X, n.Y, 0.0, []*sqt.Node{}})
				centroids = &_c

			}

			(*centroids)[cid].Weight += n.Weight
			(*centroids)[cid].Nodes = append((*centroids)[cid].Nodes, n)
			
			if (cid != members[n]) {
				members[n]  = cid
				changeCount++

			}		
		}


		// Check for convergence
		delta := 0.0
		for cid, _ := range *centroids {
			delta += centroidWss(&(*centroids)[cid])
		}
		dchange := math.Abs(deltaLast - delta)/delta

		
		if (delta == 0) || (dchange <= 0.1) {
			converged = true
		} else {			
			for cid, _ := range *centroids {
				recomputeCentroid(&(*centroids)[cid])
			}
		}
		
		deltaLast = delta
	}

	clusts:= make([][]*sqt.Node, len(*centroids))
	for cid, _ := range *centroids {
		clusts[cid] = (*centroids)[cid].Nodes
		www := (*centroids)[cid].Weight

		if (www > maxWeight) {
			fmt.Print("Centroid overweight")
			os.Exit(1)
		}

		if len(clusts[cid]) != len((*centroids)[cid].Nodes) {
			fmt.Println("Error making cluster list")
			os.Exit(1)
		}
	}
	
	return clusts
}
