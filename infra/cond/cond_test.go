package cond

import (
	"fmt"
	"testing"
)

func TestBuild(t *testing.T) {
	conds := []Cond{
		Cond{
			Op:    OPEQ,
			Field: "testeqint",
			Value: 10086,
		},
		Cond{
			Op:    OPEQ,
			Field: "testeqstr",
			Value: "10086",
		},
		Cond{
			Op:    OPLIKE,
			Field: "testlike",
			Value: "10086",
		},
		Cond{
			Op:    OPIN,
			Field: "testinstr",
			Value: []string{"s10086"},
		},
		Cond{
			Op:    OPIN,
			Field: "testinint",
			Value: []int{10086},
		},
		Cond{
			Op:    OPBETWEEN,
			Field: "testibetwee1",
			Value: []int{10086, 10087},
		},
		Cond{
			Op:    OPBETWEEN,
			Field: "testibetwee2",
			Value: []string{"s10086", "s10087"},
		},
	}
	for _, v := range conds {
		a, b, e := v.Build()
		if e != nil {
			t.Fatal(e)
		}
		fmt.Println(a, b)

	}
}
