/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"techidea8.com/codectl/internal/logic"
	"techidea8.com/codectl/internal/model"
)

var (
	table    string
	module   string
	function string
	savedir  string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export sechema to code",
	Long:  `export sechema to code`,
	Run: func(cmd *cobra.Command, args []string) {
		if table == "" || module == "" {
			cmd.Usage()
			return
		}

		// 如果当前数据库没指定,则去获取默认的db
		if len(database) == 0 || len(pkg) == 0 {
			if ds, err := logic.ListDs(); err != nil {
				fmt.Println(err.Error())
				return
			} else {
				for _, v := range ds {
					if v.Active == 1 {
						database = v.Database
						pkg = v.Package
					}
				}
				if len(database) == 0 {
					fmt.Println("please active current database use ds active -i xx")
					return
				} else {
					fmt.Println("default database ", database)
				}
			}
		}

		t := &model.Table{
			Module:    module,
			TableName: table,
			Method:    model.BuildMethod(function),
			Title:     comment,
			Database:  database,
			Package:   pkg,
			Savedir:   savedir,
			Tpldir:    tpldir,
		}
		if t.Title == "" {
			t.Title = t.Module
		}
		if err := logic.Export(t); err != nil {
			fmt.Println(err.Error())
		}
	},
}

//

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&table, "table", "t", "", "table for export ,eg: kf_test")
	exportCmd.Flags().StringVarP(&module, "module", "m", "", "module for export,eg: test")
	exportCmd.Flags().StringVarP(&function, "function", "f", "search,take,update,create,delete", "function create")
	exportCmd.Flags().StringVarP(&comment, "comment", "c", "", "comment for table")
	exportCmd.Flags().StringVarP(&savedir, "savedir", "s", "", "default dir for save")
}
