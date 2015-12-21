package trip

import (
	"fmt"
	"github.com/jcardente/santaStolen/types/location"
	"github.com/jcardente/santaStolen/types/gift"
)



const (
	sleighWeight = 10.0
	weightLimit  = 1000.0
)


type Trip struct {
	Id int
	Gifts []int
}


func TripNew(tid int) *Trip {
	return &Trip{tid, []int{}}
}

func (t *Trip) AddGift(gid int) {
	t.Gifts = append(t.Gifts, gid)
}

func (t *Trip) Score(gifts *map[int]gift.Gift) float64 {
	// NB - compute cost in reverse for efficiency
	lastPos    := location.NorthPole
	tripWeight := sleighWeight
 	tripWRW    := 0.0

	gids := t.Gifts
	for i:= len(gids) -1; i >=0; i-- {
		g := (*gifts)[gids[i]]
		tripWRW    += location.Dist(g.Location, lastPos) * tripWeight
		lastPos     = g.Location
		tripWeight += g.Weight
	}

	if (tripWeight > weightLimit) {
		fmt.Println("Warning: Trip ", t.Id, " over weight limit")
	}
	tripWRW  += location.Dist(location.NorthPole,
		(*gifts)[gids[0]].Location) * tripWeight

 	return tripWRW
}
