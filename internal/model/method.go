package model

import "strings"

type Method struct {
	Name   string `json:"name" form:"name"`
	Title  string `json:"title" form:"title"`
	Enable bool   `json:"enable" form:"enable"`
}

var DefaultMethods []Method = []Method{
	{Name: "search", Enable: true, Title: "搜索"},
	{Name: "create", Enable: true, Title: "创建"},
	{Name: "update", Enable: true, Title: "更新"},
	{Name: "delete", Enable: true, Title: "删除"},
	{Name: "take", Enable: true, Title: "查询"},
	{Name: "download", Enable: true, Title: "下载"},
}

func BuildMethod(input string) []Method {
	tmp := make([]Method, 0)
	for _, v := range DefaultMethods {

		tmp = append(tmp, Method{
			Name: v.Name, Enable: strings.Contains(input, v.Name), Title: v.Title,
		})
	}
	return tmp

}
