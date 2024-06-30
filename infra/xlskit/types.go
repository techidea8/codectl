package xlskit

import "strings"

const chars string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var charfileds []string = strings.Split(chars, "")

type Meta struct {
	Index int    `json:"index"`
	Field string `json:"field"`
	Title string `json:"title"`
}
type DataItem interface {
	Keys() []string
	Value(key string) any
}
type Record map[string]interface{}

func Conv(input map[string]interface{}) Record {
	return Record(input)
}
func (m Record) Keys() []string {
	keys := []string{}

	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
func (m Record) Value(key string) any {
	return m[key]
}
