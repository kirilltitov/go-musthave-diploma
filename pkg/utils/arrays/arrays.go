package arrays

// IsInArray Returns true if given comparable element exists in given slice.
// Caution: this function may have the complexity of O(n) at worst
func IsInArray[E comparable](s []E, v E) bool {
	for _, vs := range s {
		if v == vs {
			return true
		}
	}
	return false
}
