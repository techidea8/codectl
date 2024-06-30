package record

func Keys[K comparable, V any](input map[K]V) []K {
	result := []K{}
	for k := range input {
		result = append(result, k)
	}
	return result
}
