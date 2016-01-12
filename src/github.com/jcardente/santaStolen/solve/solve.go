package main

import (
//	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"math/rand"
//	"strconv"
//	"github.com/jcardente/santaStolen/types/location"
	"github.com/jcardente/santaStolen/types/gift"
//	"github.com/jcardente/santaStolen/types/trip"
//	"github.com/jcardente/santaStolen/types/submission"
	"github.com/jcardente/santaStolen/solve/algs"
)


var gifts map[int]gift.Gift

func main () {
	
	giftFile := flag.String("g","","Gift file")
        subFile  := flag.String("s","","Submission file")
        alg      := flag.String("a","","Algorithm")
	seed     := flag.Int64("r",1,"Random Seed")
	quiet    := flag.Bool("q",false,"Only print result")
	flag.Parse()


	if *giftFile == "" || *subFile == "" || *alg == "" {
		fmt.Println("Error: missing gift, submission, or alg argument")
		os.Exit(1)		
	}

	if !*quiet {
		fmt.Println("Using gift file: ", *giftFile)
		fmt.Println(" Using sub file: ",  *subFile)
		fmt.Println("Using algorithm: ",  *alg)
		fmt.Println("     Using seed: ", *seed)		
	}

	rand.Seed(*seed)
	
	// LOAD GIFTS ------------------------------------------------------------
        gifts, err := gift.LoadGifts(*giftFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	if !*quiet {
		fmt.Println("Number of Gifts: ", len(gifts) ," gifts")		
	}

	
	// RUN SELECTED ALGORITHM -----------------------------------------------
	//fmt.Println("Number of algs", len(algs.Algs))
	s := (algs.Algs[*alg])(&gifts)


	// OPTIMIZE -------------------------------------------------------------
	s.OptimizeTrips(&gifts)
	//c, w :=  s.CountUndersize(&gifts)
        //fmt.Println("     Undersized: ", c, " AvgW:", w)

	
        // SCORE AND SAVE -------------------------------------------------------
	
	fmt.Println("Number of trips: ", s.NumTrips())
	fmt.Println(" Solution Score: ", s.Score(&gifts))
	
	s.SaveFile(*subFile)
}
