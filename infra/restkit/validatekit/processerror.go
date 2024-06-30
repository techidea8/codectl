package validatekit

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// 获取错误提示
func ProcessError(u interface{}, err error) (result error) {
	if err == nil { //如果为nil 说明校验通过
		return nil
	}
	pv1 := reflect.ValueOf(u)
	tpstru := reflect.TypeOf(u)
	if pv1.Kind() == reflect.Ptr {
		tpstru = reflect.TypeOf(u).Elem()
	}
	invalid, ok := err.(*validator.InvalidValidationError) //如果是输入参数无效，则直接返回输入参数错误
	if ok {
		return fmt.Errorf("输入参数错误: %s" + invalid.Error())
	}
	validationErrs := err.(validator.ValidationErrors) //断言是ValidationErrors
	for _, validationErr := range validationErrs {
		fieldName := validationErr.Field()         //获取是哪个字段不符合格式
		field, ok := tpstru.FieldByName(fieldName) //通过反射获取filed
		if ok {
			errorInfo := field.Tag.Get("errmsg") //获取field对应的reg_error_info tag值
			return fmt.Errorf("输入参数错误：%s:%s", fieldName, errorInfo)
		} else {
			return err
		}
	}
	return nil
}
