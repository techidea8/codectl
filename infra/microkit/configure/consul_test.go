package configure_test

import (
	"fmt"
	"testing"

	"github.com/techidea8/codectl/infra/microkit/configure"
)

func Test01(t *testing.T) {
	configure := configure.NewConsulConfigure("http://127.0.0.1:8500", "/turingengin/dev/test.yaml")
	var c ConsulConfig
	err := configure.Parse(&c)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(c)
	}
}
