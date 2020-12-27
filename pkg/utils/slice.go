package utils

import "reflect"

func Contains(slice interface{}, item interface{}) bool {
	sl := reflect.ValueOf(slice)

	if sl.Kind() != reflect.Slice {
		return false
	}

	found := false
	for idx := 0; idx < sl.Len() && !found; idx++ {
		found = sl.Index(idx).Interface() == item
	}

	return found
}
