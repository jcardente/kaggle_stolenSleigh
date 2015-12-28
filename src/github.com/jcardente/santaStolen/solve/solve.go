package main

import (
//	"encoding/csv"
	"flag"
	"fmt"
	"os"
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
	}

	// LOAD GIFTS ------------------------------------------------------------
        gifts, err := gift.LoadGifts(*giftFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	if !*quiet {
		fmt.Println("Read ", len(gifts) ," gifts")		
	}


	// RUN SELECTED ALGORITHM -----------------------------------------------
	fmt.Println("Number of algs", len(algs.Algs))
	s := (algs.Algs[*alg])(&gifts)


        // SCORE AND SAVE -------------------------------------------------------
	
	fmt.Println("Number of trips: ", s.NumTrips())
	fmt.Println(" Solution Score: ", s.Score(&gifts))
	s.SaveFile(*subFile)
}
