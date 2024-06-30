package logic

import (
	"github.com/techidea8/codectl/app/gen/model"
	"github.com/techidea8/codectl/infra/utils/slice"
)

func BuildMethod(methodArr []string) []model.Method {
	tmp := make([]model.Method, 0)
	for _, method := range model.AllSuportMethods {
		tmp = append(tmp, model.Method{
			Name:   method.Name,
			Title:  method.Title,
			Enable: slice.Contains(methodArr, method.Name),
		})
	}
	return tmp

}
