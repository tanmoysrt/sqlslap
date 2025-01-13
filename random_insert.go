package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var randomInsertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert random data",
	Run: func(cmd *cobra.Command, args []string) {
		op := NewJob("Insert Job", clients, runInsertJob)
		op.Start()
		op.WaitAndLog()
	},
}

func runInsertJob(runnerId int, op *Job) {
	op.WaitGroup.Add(1)
	databaseConfig, err := FetchDatabaseConfig(selectedDB)
	if err != nil {
		log.Fatalf("%v", err)

	}

	if _, ok := databaseConfig.Tables[selectedTable]; !ok {
		log.Fatalf("Table %s config not found in database %s config", selectedTable, selectedDB)
	}

	db, err := NewDB(selectedDB)
	defer func() {
		_ = db.Close()
	}()

	if err != nil {
		log.Fatalf("Error while connecting to database: %v", err)
	}

	columns := databaseConfig.Tables[selectedTable].Columns

	for {
		select {
		case <-op.StopChannel:
			fmt.Println("Stopping Insert job")
			op.WaitGroup.Done()
			return
		default:
			q := generateInsertQueryWithData(selectedTable, columns)
			_, err := db.Exec(q)
			op.AddToTotal(runnerId, 1)
			if err != nil {
				op.AddToFailed(runnerId, 1)
			}
		}
	}
}
