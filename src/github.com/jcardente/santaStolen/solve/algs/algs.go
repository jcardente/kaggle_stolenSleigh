package algs


import (
	"github.com/jcardente/santaStolen/types/gift"
	"github.com/jcardente/santaStolen/types/submission"
)

type AlgFunc func(gifts *map[int]gift.Gift) *submission.Submission

var Algs = createAlgDispatch()

func createAlgDispatch() map[string]AlgFunc {
	return map[string]AlgFunc{}
}


