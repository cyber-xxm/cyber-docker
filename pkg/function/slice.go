package function

import "cmp"

func IsEmptyArray[T interface{}](v []T) bool {
	if v == nil {
		return true
	}
	if len(v) == 0 {
		return true
	}
	return false
}

func IsEmptySlice[T interface{}](v []T) bool {
	if v == nil || len(v) == 0 {
		return true
	}
	return false
}

func PluckArrayWalk[T interface{}, R interface{}](v []T, walk func(i T) (R, bool)) []R {
	result := make([]R, 0)
	for _, item := range v {
		newItem, ok := walk(item)
		if ok {
			result = append(result, newItem)
		}
	}
	return result
}

func InSlice[T cmp.Ordered](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
