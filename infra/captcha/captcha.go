package captcha

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/techidea8/codectl/infra/cachekit"
	"github.com/wenlng/go-captcha/captcha"
)

func GenerateCaptcha() (result map[int]captcha.CharDot, b64 string, tb64 string, key string, err error) {
	capt := captcha.GetCaptcha()
	dots, b64, tb64, key, err := capt.Generate()
	if err != nil {
		return
	}

	in, _ := json.Marshal(dots)
	err = cachekit.Set(key, string(in), time.Second*30)
	if err != nil {
		return
	}

	return dots, b64, tb64, key, err
}

// 这是router
func CheckCaptcha(dots, key string) (r bool, err error) {
	if dots == "" || key == "" {
		err = fmt.Errorf("验证码或验证码参数缺失")
		return
	}

	ret, err := cachekit.Get(key)
	if err != nil {
		//fmt.Println(err.Error())
		return false, err
	}
	cacheData := ret.(string)

	src := strings.Split(dots, ",")

	var dct map[int]captcha.CharDot
	err = json.Unmarshal([]byte(cacheData), &dct)
	if err != nil {
		err = fmt.Errorf("键值无效")
		return
	}

	chkRet := false
	if (len(dct) * 2) == len(src) {
		for i, dot := range dct {
			j := i * 2
			k := i*2 + 1
			sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[j]), 64)
			sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[k]), 64)

			// 检测点位置
			// chkRet = captcha.CheckPointDist(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height))

			// 校验点的位置,在原有的区域上添加额外边距进行扩张计算区域,不推荐设置过大的padding
			// 例如：文本的宽和高为30，校验范围x为10-40，y为15-45，此时扩充5像素后校验范围宽和高为40，则校验范围x为5-45，位置y为10-50
			chkRet = captcha.CheckPointDistWithPadding(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height), 5)
			if !chkRet {
				break
			}
		}
	}

	return chkRet, nil
}

/**
 * @Description: Calculate the distance between two points
 * @param sx
 * @param sy
 * @param dx
 * @param dy
 * @param width
 * @param height
 * @return bool
 */
func checkDist(sx, sy, dx, dy, width, height int64) bool {
	return sx >= dx &&
		sx <= dx+width &&
		sy <= dy &&
		sy >= dy-height
}
