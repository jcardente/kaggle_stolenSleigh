package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/jcardente/santaStolen/types/gift"
	"github.com/jcardente/santaStolen/types/submission"	
	
)

var gifts map[int]gift.Gift
var sub submission.Submission

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

	// LOAD FILES -----------------------------------------------------------
        gifts, err := gift.LoadGifts(*giftFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
 
	sub = submission.NewSubmission()
	err = sub.LoadFile(*subFile)
	 if err != nil {
	 	fmt.Println(err)
	 	os.Exit(1)
	 }

	
	if !*quiet {
		fmt.Println("Read ", len(gifts) ," gifts")		
		fmt.Println("Read ", sub.NumTrips() ," trips")
	}


	// VALIDATE ------------------------------------------------------------
	if !sub.Validate(&gifts) {
		fmt.Println("Solution is NOT valid")
	}
	
	// SCORE ---------------------------------------------------------------
	totalWRW := sub.Score(&gifts)
	if !*quiet {
		fmt.Printf("Total WRW ")
	}
	fmt.Printf("%f\n",totalWRW)
	
}





