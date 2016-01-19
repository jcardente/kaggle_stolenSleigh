package algs

import (
	"fmt"
	"math"
	"github.com/jcardente/santaStolen/types/gift"
	"github.com/jcardente/santaStolen/types/trip"	
	"github.com/jcardente/santaStolen/types/submission"
	"github.com/jcardente/santaStolen/sqt"
	"github.com/jcardente/santaStolen/kmeans"
	
)

func init () {
	Algs["cluster"] = Cluster
}


func Cluster(gifts *map[int]gift.Gift) *submission.Submission {
	
	// Create SQT
	fmt.Println("Creating SQT....")
	
	s := sqt.NewSQT()
	for _, g := range (*gifts) {
		s.AddNode(g.Id, g.Weight, g.Location.Lat, g.Location.Lon)
	}
	
	// Split the SQT
	fmt.Println("Splitting SQT...")	
	s.Split(func (tri *sqt.Triangle) bool {
		retval := false
		if (math.Ceil(tri.Weight / trip.WeightLimit) > 2) || (len(tri.Nodes) > 10){
		 	retval = true
		 }
		return retval
	})	
	
	// Iterate over triangles and cluster nodes
	fmt.Println("Clustering......")	
	sub:= submission.NewSubmission()
	tid := 1	
	for _,tri := range s.Triangles {

		if (tri.NumNodes() == 0) {
			continue
		}

		if (len(tri.Nodes) >0 ) {

			k := int(math.Ceil(tri.Weight/trip.WeightLimit))
			clusts := kmeans.Cluster(tri.Nodes, k, trip.WeightLimit)

			// Create a trip for each cluster
			for _, clust := range clusts {
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
