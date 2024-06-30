package model

type Method struct {
	Name   string `json:"name" form:"name"`
	Title  string `json:"title" form:"title"`
	Enable bool   `json:"enable" form:"enable"`
}

var SimpleMethods []string = []string{
	"create", "create", "update", "delete", "getOne", "export", "meta",
}
var AllSuportMethods []Method = []Method{
	{Name: "search", Enable: true, Title: "搜索"},
	{Name: "create", Enable: true, Title: "创建"},
	{Name: "update", Enable: true, Title: "更新"},
	{Name: "delete", Enable: true, Title: "删除"},
	{Name: "getOne", Enable: true, Title: "查询"},
	{Name: "export", Enable: true, Title: "下载"},
	{Name: "meta", Enable: true, Title: "元数据"},
}
