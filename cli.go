package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var selectedDB string
var selectedTable string
var clients int

func init() {
	rootCmd.PersistentFlags().StringVar(&selectedDB, "db", "", "Database name to generate fake data for")
	rootCmd.MarkPersistentFlagRequired("db")

	randomCmd.PersistentFlags().StringVar(&selectedTable, "table", "", "Table name to generate fake data for")
	randomCmd.MarkPersistentFlagRequired("table")
	randomCmd.PersistentFlags().IntVar(&clients, "clients", 5, "Number of concurrent clients")

	randomCmd.AddCommand(randomInsertCmd)
	randomCmd.AddCommand(randomDeleteCmd)
	randomCmd.AddCommand(randomUpdateCmd)
	rootCmd.AddCommand(initDBCmd)
	rootCmd.AddCommand(randomCmd)
}

var rootCmd = &cobra.Command{
	Use:   "sqlslap",
	Short: "sqlslap is db real world load testing tool",
	Long:  `sqlslap is a tool to generate and populate databases with random data and load test them`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var initDBCmd = &cobra.Command{
	Use: "init-db",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("This action will drop all tables and recreate them. Do you want to continue? (y/n) : ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			log.Println("Exiting...")
			return
		}
		db, err := NewDB(selectedDB)
		if err != nil {
			log.Fatalf("Error while connecting to database: %v", err)
		}
		defer db.Close()
		existedTables, err := FetchTables(selectedDB)
		if err != nil {
			log.Fatalf("Error while fetching tables: %v", err)
		}
		for _, table := range existedTables {
			_, err := db.Exec(fmt.Sprintf("DROP TABLE %s", table))
			if err != nil {
				log.Fatalf("Error while dropping table: %v", err)
			}
			log.Printf("Table %s dropped successfully", table)
		}

		databaseConfig, err := FetchDatabaseConfig(selectedDB)
		if err != nil {
			log.Fatalf("%v", err)
		}
		queries := strings.Split(databaseConfig.Schema, ";")
		for _, query := range queries {
			if strings.TrimSpace(query) == "" {
				continue
			}
			_, err := db.Exec(query)
			if err != nil {
				log.Fatalf("Error while creating table: %v", err)
			}
		}

	},
}

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Do random operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
