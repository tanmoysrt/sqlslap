package main

import (
	"github.com/spf13/cobra"
)

var selectedTable string
var clients int

func init() {
	randomCmd.PersistentFlags().StringVar(&selectedTable, "table", "", "Table name to generate fake data for")
	randomCmd.MarkPersistentFlagRequired("table")
	randomCmd.PersistentFlags().IntVar(&clients, "clients", 5, "Number of concurrent clients")

	randomCmd.AddCommand(randomInsertCmd)
	randomCmd.AddCommand(randomDeleteCmd)
	randomCmd.AddCommand(randomUpdateCmd)
}

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Do random operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
