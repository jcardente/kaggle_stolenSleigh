package types

import (
	"math"
	"strconv"
)

type Loc struct {
	Lat float64
	Lon float64
} 


func LocNew(Lat float64, Lon float64) Loc {
	return Loc{Lat * math.Pi / 180, Lon * math.Pi / 180}
}


type Gift struct {
	Id        int
        Location  Loc
	Weight    float64
}

func GiftNew(id string, lat string, lon string, weight string) Gift {

	idnew, _  := strconv.Atoi(id)
	latnew, _ := strconv.ParseFloat(lat, 64)
	lonnew, _ := strconv.ParseFloat(lon, 64)
	wnew, _   := strconv.ParseFloat(weight, 64)
		
	return Gift{ idnew, LocNew(latnew, lonnew), wnew}
}

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
