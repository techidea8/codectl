package stringx

import (
	"regexp"
	"strings"
)

// 驼峰命名
// A_b_c =>aBC
func CamelLcFirst(input string) string {
	return Lcfirst(CamelUcFirst(input))
}
func CamelUcFirst(input string) string {
	p := regexp.MustCompile(`_\w{1}`)
	return p.ReplaceAllStringFunc(input, func(s string) string {
		return strings.ToUpper(s[1:])
	})
}
