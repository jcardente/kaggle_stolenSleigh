package submission

import (
	"strconv"
	"os"
	"encoding/csv"
//	"github.com/jcardente/santaStolen/types/location"
	"github.com/jcardente/santaStolen/types/gift"
	"github.com/jcardente/santaStolen/types/trip"
)



type Submission struct {
  Trips  map[int]*trip.Trip

}

func NewSubmission() Submission {
	return Submission{map[int]*trip.Trip{}}
}

func (s Submission) LoadFile(subFile string) error {

	subfile, err := os.Open(subFile)
	if err == nil {

		defer subfile.Close()

		reader   := csv.NewReader(subfile)
		rec, err := reader.Read()
		reader.FieldsPerRecord = len(rec);		

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

func (s Submission) NumTrips() int {
	return len(s.Trips)
}


func (s Submission) Score(gifts *map[int]gift.Gift) float64 {

	totalWRW := 0.0
	for  _,t := range s.Trips {
		totalWRW += t.Score(gifts)		
	}
	
  return totalWRW
}

