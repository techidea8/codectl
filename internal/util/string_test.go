package util

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCalCase(t *testing.T) {

	// convey.Convey("Hello_word_mine_test", t, func() {
	// 	result := UnderlineToCamelCase("Hello_word_mine_test")
	// 	convey.So(result, convey.ShouldEqual, "helloWordMineTest")
	// })
	convey.Convey("_hello_word_mine_test", t, func() {
		result := UnderlineToCamelCase("_hello_word_mine_test")
		convey.So(result, convey.ShouldEqual, "helloWordMineTest")
	})
}
