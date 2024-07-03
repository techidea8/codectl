package main

import (
	// 安装依赖 go get -u github.com/spf13/cobra/cobra

	"github.com/techidea8/codectl/api/cmd/gen/biz"
	"github.com/techidea8/codectl/api/cmd/gen/tpl"
)

// 入口函数
func main() {
	tpl.Release()
	biz.Execute()
}
