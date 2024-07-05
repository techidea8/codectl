package slicekit

func Contains[T comparable](arr []T, ele T) bool {
	result := false
	for _, v := range arr {
		if v == ele {
			result = true
			break
		}
	}
	return result
}
func HasElement[T comparable](arr []T, str T) bool {
	return Contains(arr, str)
}
func HasSubStr(strarr []string, str string) bool {
	return Contains(strarr, str)
}
