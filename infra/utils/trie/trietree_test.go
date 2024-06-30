package trie

import (
	"fmt"
	"testing"
)

func TestTrieTree(t *testing.T) {
	tree := NewTrieTree("/")
	routers := map[string]interface{}{
		"/acc/create":  "createfunc",
		"/acc/search":  "searchfunc",
		"/acc/:action": "action",
		"/acc/delete":  "delete"}
	for path, funcname := range routers {
		tree.Insert(path, funcname)
	}
	f, ok := tree.Find("/acc/create")
	if ok {
		fmt.Println(f)
	}

}
