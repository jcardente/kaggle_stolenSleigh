package algs

import (
	"fmt"
	"github.com/jcardente/santaStolen/types/gift"
	"github.com/jcardente/santaStolen/types/trip"	
	"github.com/jcardente/santaStolen/types/submission"
	
)

func init () {
	Algs["naive"] = Naive
}


func Naive(gifts *map[int]gift.Gift) submission.Submission {

	fmt.Println("Called Naive function")
	
	s:= submission.NewSubmission()
	tid := 1
	for _,g  := range (*gifts) {
		t := trip.TripNew(tid)
		t.AddGift(g.Id)
		s.AddTrip(t)
		tid++
	}

	fmt.Println("Last tid: ", tid)
	return s
}
