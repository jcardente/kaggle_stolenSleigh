package gift

import (
        "strconv"
        "os"
        "encoding/csv"
        "github.com/jcardente/santaStolen/types/location"
)

type Gift struct {
        Id        int
        Location  location.Loc
        Weight    float64
}

func GiftNew(id string, lat string, lon string, weight string) Gift {

        idnew, _  := strconv.Atoi(id)
        latnew, _ := strconv.ParseFloat(lat, 64)
        lonnew, _ := strconv.ParseFloat(lon, 64)
        wnew, _   := strconv.ParseFloat(weight, 64)

        return Gift{ idnew, location.LocNew(latnew, lonnew), wnew}
}


func LoadGifts(giftFile string)  (map[int]Gift, error) {
        gifts := map[int]Gift{}
        csvfile, err := os.Open(giftFile)
        if err == nil {

                defer csvfile.Close()

                reader   := csv.NewReader(csvfile)
                rec, err := reader.Read()
                reader.FieldsPerRecord = len(rec);
                for  true {
                        rec, err = reader.Read()
                        if err != nil {
                                break
                        }
                        g := rec2Gift(rec)
                        gifts[g.Id] = g
                }

        }

        return gifts, err
}

func rec2Gift(rec []string) Gift {
  return GiftNew(rec[0],rec[1],rec[2],rec[3])
}
