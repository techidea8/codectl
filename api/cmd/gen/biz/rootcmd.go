package biz

import (
	"os"

	"github.com/spf13/cobra" // 安装依赖 go get -u github.com/spf13/cobra/cobra
	"github.com/techidea8/codectl/app/gen/conf"
	"github.com/techidea8/codectl/app/gen/logic"
	"github.com/techidea8/codectl/infra/logger"
)

// 这个是根命令定义
var rootCmd = &cobra.Command{
	Use:   "codectl", // 这个就是你的自己定义的根命令
	Short: "代码生成工具",
	Long:  `通过数据库生成代码`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// 命令执行方法
func Execute() {
	cobra.OnInitialize(func() {
		c := conf.DefaultAppConf
		c.Env = conf.ENVDEF(env)
		if env == string(conf.DEV) {
			logger.SetLevel(logger.DebugLevel)
		} else {
			logger.SetLevel(logger.InfoLevel)
		}
		logic.InitApp(c)
	})

	if len(os.Args) <= 1 {
		rootCmd.SetArgs([]string{"-h"})
	}
	// 执行命令 如果异常会返回错误信息
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&env, "env", "e", string(conf.PROD), "ENV OF APP: prod/dev")
}
