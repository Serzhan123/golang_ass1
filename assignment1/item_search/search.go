package itemsearch

import (
	"fmt"
	"strings"

	"github.com/Bektemis/golang_ass_1/pck"
)

func ItemSearch(it string, items *pck.DatabaseItems) []string {
	var list []string
	for _, item := range items.Items {
		if strings.Contains(item.Name, it) {
			list = append(list, fmt.Sprintf("Name: %s, Price: %d, Rating: %d ", item.Name, item.Price, item.Rating))
		}
	}
	return list
}
