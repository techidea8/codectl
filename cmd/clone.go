/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"techidea8.com/codectl/internal/logic"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clone a template from git repo",
	Long: `clone a template from git repo For example:
codectl clone http://github.com/techidea8/codetpl/golang-vue3.git`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		} else {
			result, err := logic.Clone(args[0])
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(result)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
