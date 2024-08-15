package stringx

import "strings"

// 是否有前缀
// hasprefix("/a/hello",["h","/a"])
func HasPrefix(s string, prefixs ...string) bool {
	has := false
	for _, prefix := range prefixs {
		has = has || strings.HasPrefix(s, prefix)
		if has {
			break
		}
	}
	return has
}

// 是否有前缀
// hassuffix("/a/hello",["h","/a"])
func HasSuffix(s string, prefixs ...string) bool {
	has := false
	for _, prefix := range prefixs {
		has = has || strings.HasSuffix(s, prefix)
		if has {
			break
		}
	}
	return has
}
