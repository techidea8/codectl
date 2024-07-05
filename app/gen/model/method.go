package model

type Method struct {
	Name   string `json:"name" form:"name"`
	Title  string `json:"title" form:"title"`
	Enable bool   `json:"enable" form:"enable"`
}

var SimpleMethods []string = []string{
	"create", "create", "update", "delete", "getOne", "export", "meta", "deleteIts", "print",
}
var AllSuportMethods []Method = []Method{
	{Name: "search", Enable: true, Title: "搜索"},
	{Name: "create", Enable: true, Title: "创建"},
	{Name: "update", Enable: true, Title: "更新"},
	{Name: "delete", Enable: true, Title: "删除"},
	{Name: "getOne", Enable: true, Title: "查询"},
	{Name: "export", Enable: true, Title: "下载"},
	{Name: "meta", Enable: true, Title: "元数据"},
	{Name: "deleteIts", Enable: true, Title: "批量删除"},
	{Name: "print", Enable: false, Title: "打印"},
}
