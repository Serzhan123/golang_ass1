package pck

import "fmt"

type Item struct {
	Name      string
	Price     int
	Rating    int
	HaveRated int
}
type DatabaseItems struct {
	Items []Item
}

type AccessToItems interface {
	GetListOfItems() []string
	FilterByPrice(bool)
	FilterByRatings(bool)
}

func (items *DatabaseItems) GetListOfItems() []string {
	var list []string
	for _, item := range items.Items {
		list = append(list, fmt.Sprintf("Name: %s, Price: %d, Rating: %d ", item.Name, item.Price, item.Rating))
	}
	return list
}

func (item *Item) ChangeRating(rating int) {
	item.Rating = rating
}
