/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"techidea8.com/codectl/internal/model"
)

var (
	lang     string
	tpldir   string
	comment  string
	database string
	prefix   string
	id       int
	pkg      string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "codectl",
	Short: "代码生吃器",
	Long:  `代码生成器,支持golang,java,只需要定义相关模板即可`,
	PreRun: func(cmd *cobra.Command, args []string) {
		model.InitConfig()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&prefix, "prefix", "x", "", "prefix will be removed: kf_test->test")
	rootCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "go", "code lang")
	rootCmd.PersistentFlags().StringVarP(&tpldir, "tpldir", "p", "", "dir lang")
	rootCmd.PersistentFlags().StringVarP(&database, "database", "d", "", "current database for table")
	rootCmd.PersistentFlags().StringVarP(&pkg, "pkg", "g", "turinglet", "pkg of appname")
	rootCmd.PersistentFlags().IntVarP(&id, "id", "i", 0, "set id for record")
}
