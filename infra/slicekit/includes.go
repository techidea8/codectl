package slicekit

// Includes items in slice
func Includes[T comparable](slice []T, item T) bool {
	for _, _item := range slice {
		if _item == item {
			return true
		}
	}
	return false
}
