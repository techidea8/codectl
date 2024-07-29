package stringx

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomNumber(min, max int) string {
	rand.Seed(time.Now().UnixNano()) // 初始化随机数生成器
	num := rand.Intn(max-min+1) + min
	return fmt.Sprintf("%d", num)
}
