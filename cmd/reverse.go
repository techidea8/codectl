/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"techidea8.com/codectl/internal/logic"
)

var reverseCmd = &cobra.Command{
	Use:   "reverse",
	Short: "transfer all table to code",
	Long:  `transfer all table to code: reverse -x kf_ -l go`,
	Run: func(cmd *cobra.Command, args []string) {
		logic.Reverse(database)
	},
}

func init() {
	rootCmd.AddCommand(reverseCmd)

}
