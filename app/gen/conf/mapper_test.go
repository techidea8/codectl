package conf

import (
	"testing"
)

func TestFile(t *testing.T) {
	ResetMapperRule("mysql-golang.yml")
}
func TestMap(t *testing.T) {
	ResetMapperRule(map[string]map[string]string{
		"mysql-golang": map[string]string{
			"datetime": "types.DateTime",
		},
	})
}
