package main

// 144525525772.0

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"github.com/jcardente/santaStolen/types"
	"github.com/jcardente/santaStolen/haversine"		
)


const (
	sleighWeight = 10.0
	weightLimit  = 1000.0
)

var northPole types.Loc


func init() {
	northPole = types.LocNew(90,0)
}


var gifts map[int]types.Gift
var trips map[int]*types.Trip

func main () {
	
	giftFile := flag.String("g","","Gift file")
        subFile  := flag.String("s","","Submission file")
	quiet    := flag.Bool("q",false,"Only print result")
	flag.Parse()

	if *giftFile == "" || *subFile == "" {
		fmt.Println("Error: missing gift or submission file argument")
		os.Exit(1)
		
	}

	if !*quiet {
		fmt.Println("Using gift file ", *giftFile)
		fmt.Println("Using sub file ",  *subFile)
	}

	// GIFTS FILE ------------------------------------------------------------
	csvfile, err := os.Open(*giftFile)
	if err !=nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer csvfile.Close()

	
	reader   := csv.NewReader(csvfile)
	rec, err := reader.Read()
	reader.FieldsPerRecord = len(rec);


	gifts = map[int]types.Gift{}	
	for  true {
		rec, err = reader.Read()
		if err != nil {
			break
		}
		g := rec2Gift(rec)
		gifts[g.Id] = g
	}
	if !*quiet {
		fmt.Println("Read ", len(gifts) ," gifts")
	}

	// SUBMISSION FILE ------------------------------------------------------------
	subfile, err := os.Open(*subFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
        defer subfile.Close()

	reader   = csv.NewReader(subfile)
	rec, err = reader.Read()
	reader.FieldsPerRecord = len(rec);


	trips = make(map[int]*types.Trip)
        for true {
		rec, err = reader.Read()
		if err != nil {
			break
		}
		
		gid, _ := strconv.Atoi(rec[0])
		tid, _ := strconv.Atoi(rec[1])

		trip, exists := trips[tid]
		if !exists {
			trips[tid] = types.TripNew(tid)
			trip       = trips[tid]
		}

		trip.AddGift(gid)
	}
	if !*quiet {
		fmt.Println("Read ", len(trips) ," trips")
	}


	
	// SCORING ------------------------------------------------------------	
	count    := 0
	totalWRW := 0.0
	for  _,t := range trips {

		// NB - compute cost in reverse for efficiency
		lastPos    := northPole
		tripWeight := sleighWeight
		tripWRW    := 0.0

		gids := t.Gifts
		for i:= len(gids) -1; i >=0; i-- {
			g := gifts[gids[i]]
			tripWRW    += haversine.Dist(g.Location, lastPos) * tripWeight
			lastPos     = g.Location
			tripWeight += g.Weight
		}
		tripWRW  += haversine.Dist(northPole, gifts[gids[0]].Location) * tripWeight			
		totalWRW += tripWRW
		
		count++			
	}

	if !*quiet {
		fmt.Println("Read ",count," trips")
		fmt.Printf("Total WRW %f \n", totalWRW)
	} else {
		fmt.Printf("%f\n",totalWRW)
	}
}




func rec2Gift(rec []string) types.Gift {
  return types.GiftNew(rec[0],rec[1],rec[2],rec[3])
}

