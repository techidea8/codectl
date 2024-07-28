package biz

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
)

var rule1 *regexp.Regexp = regexp.MustCompile(`//\s+(post|get|put|delete|options\S*)\s+([\/\w]+)`)
var regmap map[string]*regexp.Regexp = map[string]*regexp.Regexp{
	"// post /acc/create": rule1,
	"// gen by codectl ,donot modify ,https://github.com/techidea8/codectl.git": rule1,
	"post /acc/create":            rule1,
	"// post,get /acc/create":     rule1,
	"// post,get,put /acc/create": rule1,
}

func String(a any) string {
	s, _ := json.Marshal(a)
	return string(s)
}
func TestReg(t *testing.T) {
	for k, v := range regmap {
		arr := v.FindStringSubmatch(k)
		fmt.Println(k, String(arr))
		//t.Log(arr)
	}
}

func Test02(t *testing.T) {
	str1 := "// post,get,put /acc/create"
	patern := regexp.MustCompile(`//\s+(post|get|put[\,post|\,get|\,put]*)\s+[\/\w]+`)
	result := patern.FindStringSubmatch(str1)
	fmt.Println(str1, len(result), result)
}
