package haversine


import (
	"testing"
	"github.com/jcardente/santaStolen/types"
)


func TestDist(t *testing.T) {
	d := Dist(types.LocNew(36.12, -86.67), types.LocNew(33.94,-118.40))
	if d != 2886.4444428379825 {
		t.Error("Expected 2887.25, got ", d)
	}
}
