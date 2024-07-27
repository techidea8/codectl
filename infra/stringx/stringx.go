package stringx

import (
	"html/template"
	"strings"
)

func Ucfirst(input string) (result string) {
	if len(input) < 1 {
		return ""
	}
	return (strings.ToUpper(input[:1]) + input[1:])
}
func Upper(input string) (out string) {
	return strings.ToUpper(input)
}
func Lower(input string) (out string) {
	return strings.ToLower(input)
}
func Lcfirst(input string) (result string) {
	if len(input) < 1 {
		return ""
	}
	return (strings.ToLower(input[:1]) + input[1:])
}

func JSStr(input string, args ...string) (result template.JSStr) {
	return template.JSStr(strings.Join(append([]string{input}, args...), ""))
}

func JS(input string, args ...string) (result template.JS) {
	return template.JS(strings.Join(append([]string{input}, args...), ""))
}

// 下划线转驼峰命名
func UnderlineToCamelCase(input string) string {
	strArr := strings.Split(input, "_")
	result := make([]string, 0)
	for _, v := range strArr {
		result = append(result, Ucfirst(v))
	}
	return Lcfirst(strings.Join(result, ""))
}

// 下划线转驼峰命名
func UnderlineToUperCamelCase(input string) string {
	strArr := strings.Split(input, "_")
	result := make([]string, 0)
	for _, v := range strArr {
		result = append(result, Ucfirst(v))
	}
	return Ucfirst(strings.Join(result, ""))
}
