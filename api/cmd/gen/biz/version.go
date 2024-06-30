package biz

import (
	"fmt"

	"github.com/spf13/cobra" // 安装依赖 go get -u github.com/spf13/cobra/cobra
)

// 定义一个参数
var Verbose bool

// 子命令定义 运行方法 go run main.go version 编译后 ./hugo version
var versionCmd = &cobra.Command{
	Use:   "version", // Use这里定义的就是命令的名称
	Short: "Print the version number ",
	Long:  `All software has versions. This is codectl`,
	Run: func(cmd *cobra.Command, args []string) { //这里是命令的执行方法
		fmt.Println("version is 1.0")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		//这个在命令执行前执行
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		//这个在命令执行后执行
	},
	// 还有其他钩子函数
}

func init() {
	//给我们定义的命令绑定参数 可以给我们定义的任何命令绑定参数
	rootCmd.AddCommand(versionCmd)
}
