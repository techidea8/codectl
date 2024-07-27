package biz

//环境变量
var env string
var configfile string
var rulefile string = ""
var dirsave string = ""

//-source .  -dstpkg turingdance.com/turing/devfast/api/rest/sys -dstfile routefunc.go
// 需要解析router 的
var source string = "."

//需要的包路径
var dstpkg string = "techidea8.com/turing/app"

// 生成的router 文件

var dstfile = "router_file.go"
