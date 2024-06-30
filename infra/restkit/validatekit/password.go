package validatekit

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func PasswordScore(f validator.FieldLevel) bool {
	str := f.Param()
	dstscore, _ := strconv.Atoi(str)
	value := f.Field().String()
	score := 0

	bts := []byte(value)
	// 包含英文
	reglower := regexp.MustCompile(`[a-z]+`)
	if reglower.Match(bts) {
		score += 1
	}
	// 大小写
	reglower = regexp.MustCompile(`[A-Z]+`)
	if reglower.Match(bts) {
		score += 1
	}
	// 数字
	reglower = regexp.MustCompile(`\d+`)
	if reglower.Match(bts) {
		score += 1
	}
	//特殊字符
	reglower = regexp.MustCompile(`[,!\?\:\'\"\(\)\$\%\^\*\@]+`)
	if reglower.Match(bts) {
		score += 1
	}
	// 大于5位,否则零分
	if len(value) > 5 {
		score *= 1
	} else {
		score *= 0
	}
	if dstscore > score {
		return false
	} else {
		return true
	}
}
