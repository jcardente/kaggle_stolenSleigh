package submission

import (
        "strconv"
        "os"
        "fmt"
	"math"
	"sort"
        "encoding/csv"
//      "github.com/jcardente/santaStolen/types/location"
        "github.com/jcardente/santaStolen/types/gift"
        "github.com/jcardente/santaStolen/types/trip"
)



type Submission struct {
        Trips   map[int]*trip.Trip
}

func NewSubmission() *Submission {
        return &Submission{map[int]*trip.Trip{}}
}

func (s *Submission) AddTrip(t *trip.Trip) {
        s.Trips[(*t).Id] = t
}

func (s *Submission) NumTrips() int {
        return len(s.Trips)
}



type TripOpt struct {
	T *trip.Trip
	Lat float64
}

type TripOptList []TripOpt

func NewTripOpt(t *trip.Trip, l float64) TripOpt {
	return TripOpt{t, l}
}

func (to TripOptList) Len() int {
	return len(to)
}

func (to TripOptList) Swap(i, j int) {
	to[i], to[j] = to[j], to[i]
}

func (to TripOptList) Less(i, j int) bool {
   return to[i].Lat < to[j].Lat
}



type TripList []*trip.Trip

func (tl TripList) Len() int {
	return len(tl)
}

func (tl TripList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (tl TripList) Less(i, j int) bool {
	// Sort from lighter to heavier
	return tl[i].Weight > tl[j].Weight
}


func (s *Submission) OptimizeTrips(gifts *map[int]gift.Gift) {

	// OPTIMIZE BY LATITUDE
	tripsByLon := map[int]TripOptList{}
	res := 0.5
        for _, t := range s.Trips {		
		t.OptimizeOrder(gifts)		
		t.CalcWeight(gifts)
		t.CacheScore(gifts)
		l := t.GetLongitude(gifts)
		b := int(math.Floor(l / res))
		tripsByLon[b] = append(tripsByLon[b], NewTripOpt(t, t.GetLatitude(gifts)))
	}

	MergedTrips := map[int]*trip.Trip{}	
	for _, tol := range tripsByLon {
	 	sort.Sort(tol)		
		skipTrip := map[*trip.Trip]bool{}
		for i := 0; i < len(tol); i++ {
			t1 := tol[i].T
			if (skipTrip[t1] == true) {
				continue
			}
			
			tnew := &trip.Trip{t1.Id, t1.Gifts, t1.Weight, t1.WRW}
			for j:= i+1; j < len(tol); j++ {
				t2 := tol[j].T
				if (skipTrip[t2] == true) {
					continue
				}

				if (tnew.Weight + t2.Weight) > trip.WeightLimit {
					continue
				}
			
				// Make sure to complete closest trip first
				ttest := trip.TripNew(-1)				
				ttest.Gifts = append(ttest.Gifts, tnew.Gifts...)
				ttest.Gifts = append(ttest.Gifts, t1.Gifts...)				
				
				newscore := ttest.Score(gifts)
				if (newscore < (tnew.WRW + t2.WRW)) {
					tnew.Gifts = append(t2.Gifts, tnew.Gifts...)
					skipTrip[t2] = true
					tnew.Score(gifts)
					tnew.CalcWeight(gifts)
				}				
			}
			MergedTrips[tnew.Id] = tnew
		}
	}

	fmt.Println(" LATOPT: ", len(s.Trips), " --> ", len(MergedTrips))
	s.Trips = MergedTrips


	
	// OPTIMIZE BY SPACE AVAILABLE
	MergedTrips = map[int]*trip.Trip{}
	unfilledTrips := TripList{}
	for _, t := range s.Trips {
		t.CalcWeight(gifts)
		t.CacheScore(gifts)		
		if ((trip.WeightLimit - t.Weight) > 100) {
			unfilledTrips = append(unfilledTrips, t)
		} else {
			MergedTrips[t.Id] = t
		}		
	}

	fmt.Println("Unfilled: ", len(unfilledTrips))
	sort.Sort(unfilledTrips)
	skipTrip := map[*trip.Trip]bool{}	
	for i := 0; i < len(unfilledTrips); i++ {
		t1 := unfilledTrips[i]
		if (skipTrip[t1] == true) {
			continue
		}
		
		tnew := &trip.Trip{t1.Id, t1.Gifts, t1.Weight, t1.WRW}
		for j:=i+1; j < len(unfilledTrips); j++ {
			t2 := unfilledTrips[j]
			if (skipTrip[t2] == true) {
				continue
			}

			if (tnew.Weight + t2.Weight) > trip.WeightLimit {
				continue
			}

			ttest := trip.TripNew(-1)					
			ttest.Gifts = append(ttest.Gifts, t2.Gifts...)
			ttest.Gifts = append(ttest.Gifts, tnew.Gifts...)				
			newscore := ttest.Score(gifts)
			if (newscore <= (tnew.WRW + t2.WRW)) {
				tnew.Gifts = append(t2.Gifts, tnew.Gifts...)
				skipTrip[t2] = true
				tnew.Score(gifts)
				tnew.CalcWeight(gifts)
			}
		}
		MergedTrips[tnew.Id] = tnew
	}
	fmt.Println(" SPACEOPT: ", len(s.Trips), " --> ", len(MergedTrips))
	s.Trips = MergedTrips


}

func (s *Submission) CountUndersize(gifts *map[int]gift.Gift) (int, float64) {
	count := 0
	weight:= 0.0
        for _, t := range s.Trips {
		t.CalcWeight(gifts)
		if t.Weight < trip.WeightLimit {
			count++
			weight += (trip.WeightLimit - t.Weight)
		}
	}

	weight = weight/float64(count)
	
	return count, weight
}

func (s *Submission) Validate(gifts *map[int]gift.Gift) bool {

        valid := true
        seen := make(map[int]bool)
        for _, t := range s.Trips {
                for _, gid := range (*t).Gifts {
                        if seen[gid] {
                                fmt.Println("Warning, gift seen twice ", gid)
                                valid = false
                        }
                        seen[gid] = true
                }
        }

        for gid, _ := range *gifts {
                if !seen[gid] {
                        fmt.Println("Warning gift not seen ", gid)
                        valid = false
                }
        }

        return valid
}

func (s *Submission) Score(gifts *map[int]gift.Gift) float64 {

        totalWRW := 0.0
        for  _,t := range s.Trips {
                totalWRW += t.Score(gifts)
        }

	return totalWRW
}



func (s *Submission) SaveFile(subFile string) error {
        subfile, err := os.Create(subFile)
        if err == nil {
                defer subfile.Close()

                writer := csv.NewWriter(subfile)
                writer.Write([]string{"GiftId","TripId"})
                for _, t := range s.Trips {
                        for _, gid := range (*t).Gifts {
                                writer.Write([]string{strconv.Itoa(gid), strconv.Itoa((*t).Id)})
                        }
                }
                writer.Flush()
        } else {
                fmt.Println(err)
        }
        return err

}

func (s *Submission) LoadFile(subFile string) error {

        subfile, err := os.Open(subFile)
        if err == nil {

                defer subfile.Close()

                reader   := csv.NewReader(subfile)
                rec, err := reader.Read()
                reader.FieldsPerRecord = len(rec)

                for true {
                        rec, err = reader.Read()
                        if err != nil {
                                break
                        }

                        gid, _ := strconv.Atoi(rec[0])
                        tid, _ := strconv.Atoi(rec[1])

                        t, exists := s.Trips[tid]
                        if !exists {
                                s.Trips[tid] = trip.TripNew(tid)
                                t            = s.Trips[tid]
                        }

                        t.AddGift(gid)
                }
        }
        return err
}
