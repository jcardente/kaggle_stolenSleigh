package algs

import (
	"fmt"
	"math"
//	"os"
//	"github.com/jcardente/santaStolen/types/location"	
	"github.com/jcardente/santaStolen/types/gift"
	"github.com/jcardente/santaStolen/types/trip"	
	"github.com/jcardente/santaStolen/types/submission"
	"github.com/jcardente/santaStolen/sqt"
	"github.com/jcardente/santaStolen/kmeans"
	
)

func init () {
	Algs["cluster"] = Cluster
}


func Cluster(gifts *map[int]gift.Gift) submission.Submission {

	fmt.Println("Called Cluster algorithm")


	// Create SQT
	fmt.Println("Creating spherical quadtree")
	s := sqt.NewSQT()
	for _, g := range (*gifts) {
		s.AddNode(g.Id, g.Weight, g.Location.Lat, g.Location.Lon)
	}

	
	// Split the SQT
	fmt.Print("Splitting the SQT....")
	s.Split(func (tri *sqt.Triangle) bool {
		retval := false
		if (math.Ceil(tri.Weight / trip.WeightLimit) > 5) || (len(tri.Nodes) > 50){
		 	retval = true
		 }
		return retval
	})	
	fmt.Println(" count: ",len(s.Triangles))
	
	// Iterate over rectangles
	fmt.Println("Creating trips...")
	sub:= submission.NewSubmission()
	tid := 1	
	for _,tri := range s.Triangles {

		// Skip empty triangles
		if (tri.NumNodes() == 0) {
			continue
		}

		if (len(tri.Nodes) >0 ) {
			// Cluster the nodes
			k := int(math.Ceil(tri.Weight/trip.WeightLimit))

			//fmt.Println(" N", len(tri.Nodes), " K", k, " TW", tri.Weight)
			clusts := kmeans.Cluster(tri.Nodes, k, trip.WeightLimit)
		
			// For each cluster, create a trip
			foo := 0
			for _, clust := range clusts {
				foo += len(clust)
				t := trip.TripNew(tid)
				for _, n := range clust {
					t.AddGift(n.Id)
				}
				sub.AddTrip(t)
				tid++
			}
		}
	}

	return sub
}
