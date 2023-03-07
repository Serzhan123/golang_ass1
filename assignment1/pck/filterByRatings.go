package pck

func (il *DatabaseItems) FilterByRatings(ascending bool) {
	if ascending {
		for i := 0; i < len(il.Items); i++ {
			for j := i + 1; j < len(il.Items); j++ {
				if il.Items[i].Rating > il.Items[j].Rating {
					il.Items[i], il.Items[j] = il.Items[j], il.Items[i]
				}
			}
		}
	} else {
		for i := 0; i < len(il.Items); i++ {
			for j := i + 1; j < len(il.Items); j++ {
				if il.Items[i].Rating < il.Items[j].Rating {
					il.Items[i], il.Items[j] = il.Items[j], il.Items[i]
				}
			}
		}
	}
}
