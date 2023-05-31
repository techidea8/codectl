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
	dbtype   string
	user     string
	password string
	port     int
	query    string
	host     string
	active   int
)

type Dnsfunc func(*model.Ds) error

var dnsfuncmap map[string]Dnsfunc = map[string]Dnsfunc{
	"list":   logic.ListDns,
	"add":    logic.CreateDns,
	"update": logic.UpdateDns,
	"active": logic.ActiveDns,
	"remove": logic.RemoveDns,
}

// dnsCmd represents the dns command
var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "operate datasource:dns add/update/list/active/remove",
	Long:  `operate datasource: dns list/add/update/active/remove`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			return
		}
		cmdstr := args[0]

		ds := &model.Ds{
			Database: database,
			User:     user,
			Port:     port,
			Pass:     password,
			Query:    query,
			Host:     host,
			Dbtype:   dbtype,
			Id:       uint(id),
			Active:   active,
			Tpl:      tpldir,
		}
		if fun, ok := dnsfuncmap[cmdstr]; ok {
			if err := fun(ds); err != nil {
				fmt.Println(err.Error())
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
	dnsCmd.Flags().StringVarP(&user, "user", "u", "root", "user for dns")
	dnsCmd.Flags().StringVarP(&password, "password", "w", "root", "password for dns")
	dnsCmd.Flags().IntVarP(&port, "port", "p", 3306, "port for dns")
	dnsCmd.Flags().StringVarP(&query, "query", "q", "charset=utf8mb4&parseTime=True&loc=Local", "other params name")
	dnsCmd.Flags().StringVarP(&dbtype, "dbtype", "t", "mysql", "data base type mysql|sqlite|sqlite3")
	dnsCmd.Flags().IntVarP(&active, "active", "a", 2, "active or not 1|2 1=active,2=deactive")
	dnsCmd.Flags().StringVarP(&host, "host", "s", "127.0.0.1", "db host")
}
