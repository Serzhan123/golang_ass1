package ratings

import (
	"fmt"

	"github.com/Bektemis/golang_ass_1/pck"
)

func GiveRating(rating int, item string, db *pck.DatabaseItems) {
	for _, it := range db.Items {
		if it.Name == item {
			it.ChangeRating((it.Rating*it.HaveRated + rating) / (it.HaveRated + 1))
			it.HaveRated++
			fmt.Println("Have rated item, new rating:", it.Rating)
		}
	}
}
