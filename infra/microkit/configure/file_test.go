package configure_test

import (
	"fmt"
	"testing"

	"github.com/techidea8/codectl/infra/microkit/configure"
)

func Test02(t *testing.T) {
	configure := configure.NewFileConfigure("./test.yml")
	var c ConsulConfig
	err := configure.Parse(&c)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(c)
	}
}
