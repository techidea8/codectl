package stringx

import "testing"

func TestCamel(t *testing.T) {
	r := CamelUcFirst("a_b_c_d")
	t.Log(r == "aBCD")
}
