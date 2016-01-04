package submission

import (
        "strconv"
        "os"
        "fmt"
        "encoding/csv"
//      "github.com/jcardente/santaStolen/types/location"
        "github.com/jcardente/santaStolen/types/gift"
        "github.com/jcardente/santaStolen/types/trip"
)



type Submission struct {
        Trips   map[int]*trip.Trip
}

func NewSubmission() Submission {
        return Submission{map[int]*trip.Trip{}}
}

func (s Submission) AddTrip(t *trip.Trip) {
        s.Trips[(*t).Id] = t
}

func (s Submission) NumTrips() int {
        return len(s.Trips)
}


func (s Submission) Validate(gifts *map[int]gift.Gift) bool {

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

func (s Submission) Score(gifts *map[int]gift.Gift) float64 {

        totalWRW := 0.0
        for  _,t := range s.Trips {
                totalWRW += t.Score(gifts)
        }

	return totalWRW
}



func (s Submission) SaveFile(subFile string) error {
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

func (s Submission) LoadFile(subFile string) error {

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
