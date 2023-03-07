package pck

func (il *DatabaseItems) FilterByPrice(ascending bool) {
	if ascending {
		for i := 0; i < len(il.Items); i++ {
			for j := i + 1; j < len(il.Items); j++ {
				if il.Items[i].Price > il.Items[j].Price {
					il.Items[i], il.Items[j] = il.Items[j], il.Items[i]
				}
			}
		}
	} else {
		for i := 0; i < len(il.Items); i++ {
			for j := i + 1; j < len(il.Items); j++ {
				if il.Items[i].Price < il.Items[j].Price {
					il.Items[i], il.Items[j] = il.Items[j], il.Items[i]
				}
			}
		}
	}
}
