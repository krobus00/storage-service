package util

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func ContainsMultiple[T comparable](elems []T, v []T) bool {
	for _, s := range v {
		if Contains(elems, s) {
			return true
		}
	}
	return false
}
