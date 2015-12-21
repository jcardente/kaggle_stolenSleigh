package location

import (
	"math"
)
	
type Loc struct {
	Lat float64
	Lon float64
} 


func LocNew(Lat float64, Lon float64) Loc {
	return Loc{Lat * math.Pi / 180, Lon * math.Pi / 180}
}


var NorthPole = Loc{90 * math.Pi/ 180, 0}


const EarthRadius = 6371

func hav(theta float64) float64 {
	return 0.5 * (1 - math.Cos(theta))
}

func Dist(l1 Loc, l2 Loc) float64 {
	dLat := l2.Lat - l1.Lat
	dLon := l2.Lon - l1.Lon

	d := 2 * EarthRadius *
		math.Asin(math.Sqrt(hav(dLat) +
		                    math.Cos(l2.Lat) * math.Cos(l1.Lat) * hav(dLon)))

	return d
}

