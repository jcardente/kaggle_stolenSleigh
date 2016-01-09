package trip

import (
	"fmt"
	"os"
	"math"
	"github.com/jcardente/santaStolen/types/location"
	"github.com/jcardente/santaStolen/types/gift"
)



const (
	sleighWeight = 10.0
	WeightLimit  = 1000.0
)


type Trip struct {
	Id int
	Gifts []int
	Weight float64
	WRW    float64
}



func TripNew(tid int) *Trip {
	return &Trip{tid, []int{}, 0.0, 0.0}
}

func (t *Trip) AddGift(gid int) {
	t.Gifts = append(t.Gifts, gid)
}

func (t *Trip) CalcWeight(gifts *map[int]gift.Gift) {
	gids     := t.Gifts
	t.Weight  = 0.0
	for _, gid := range gids  {
		t.Weight += (*gifts)[gid].Weight
	}
}


func (t *Trip) OptimizeOrder(gifts *map[int]gift.Gift) {

	// Find the gift with the highest cost
	heaviestGid := -1
	maxCost     := 0.0
	gids        := t.Gifts
	tripWeight  := 0.0
	for _, gid := range gids  {
		g := (*gifts)[gid]
		tripWeight += g.Weight
		cost := location.Dist(location.NorthPole, g.Location) * g.Weight
		if (cost > maxCost) {
			heaviestGid = gid
			maxCost = cost
		}
	}

	// Start new list and identify remaining gifts
	newGids  := []int{heaviestGid}
	remGifts := make(map[int]bool)
	for _, gid := range gids {
		if (gid != heaviestGid) {
			remGifts[gid] = true
		}
	}

	// Loop until no remaining gifts
	lastGid := heaviestGid
	tripWeight -= (*gifts)[lastGid].Weight
	for len(remGifts) > 0 {
		bestGid  := -1
		bestCost := math.Inf(0)
		lastG := (*gifts)[lastGid]
		for gid, _ := range remGifts {
			g := (*gifts)[gid]
			cost := location.Dist(lastG.Location, g.Location) * tripWeight
			if (cost < bestCost) {
				bestGid = gid
				bestCost = cost
			}			
		}

		// Add to new list and remove from map
		newGids = append(newGids, bestGid)
		tripWeight -= (*gifts)[bestGid].Weight
		delete(remGifts, bestGid)		
	}

	t.Gifts = newGids	
}

func (t *Trip) GetLongitude(gifts *map[int]gift.Gift) float64 {

	// Return the longitude of the first gift in the trip list.
	// Used during optimization to combine trips along similar
	// longitudes
	lon := (*gifts)[t.Gifts[0]].Location.Lon
	lon = 180 * lon / (math.Pi)
	if (lon < 0) {
		lon = 360 - lon
	}

	return lon
}

func (t *Trip) GetLatitude(gifts *map[int]gift.Gift) float64 {
	lat := (*gifts)[t.Gifts[0]].Location.Lat
	//lat = 90 - lat

	lat = 180 * lat / (math.Pi)
	return lat
}


func (t *Trip) Score(gifts *map[int]gift.Gift) float64 {
	// NB - compute cost in reverse for efficiency
	lastPos    := location.NorthPole
	tripWeight := sleighWeight
 	tripWRW    := 0.0
	giftWeight := 0.0
	
	if len(t.Gifts) == 0 {
		fmt.Println("Warning, empty trip")
		os.Exit(1)
	}
	gids := t.Gifts
	for i:= len(gids) -1; i >=0; i-- {
		g := (*gifts)[gids[i]]
		tripWRW    += location.Dist(g.Location, lastPos) * tripWeight
		lastPos     = g.Location
		tripWeight += g.Weight
		giftWeight += g.Weight
	}

	if (t.Id >= 0) && (giftWeight > WeightLimit) {
		fmt.Println("Warning: Trip ", t.Id, " over weight limit")
	}
	tripWRW  += location.Dist(location.NorthPole,
		(*gifts)[gids[0]].Location) * tripWeight

 	return tripWRW
}

func (t *Trip) CacheScore(gifts *map[int]gift.Gift)  {
	t.WRW = t.Score(gifts)
}
