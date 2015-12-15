package haversine

import (
	"math"
	"github.com/jcardente/santaStolen/types"
)


const earthRadius = 6371

func hav(theta float64) float64 {
	return 0.5 * (1 - math.Cos(theta))
}

func Dist(l1 types.Loc, l2 types.Loc) float64 {
	dLat := l2.Lat - l1.Lat
	dLon := l2.Lon - l1.Lon

	d := 2 * earthRadius *
		math.Asin(math.Sqrt(hav(dLat) +
		                    math.Cos(l2.Lat) * math.Cos(l1.Lat) * hav(dLon)))

	return d
}

