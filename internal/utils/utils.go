package utils

// StringInSlice checks if a string is present in slice
func SliceContains[T comparable](a T, list []T) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
