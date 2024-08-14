package cryptor

import (
	"fmt"
	"testing"
)

func TestJust(t *testing.T) {
	a := ErrorTest()
	if a != nil {
		fmt.Printf("t: %v\n", a.Error())
	}
}
