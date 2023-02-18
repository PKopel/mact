package utils

func Unsnoc[T any](slice []T) ([]T, *T) {
	sliceLen := len(slice)
	if sliceLen != 0 {
		return slice[:sliceLen-1], &slice[sliceLen-1]
	}
	return slice, nil
}

func Contains[T comparable](slice []T, val T) bool {
	for _, a := range slice {
		if a == val {
			return true
		}
	}
	return false
}
