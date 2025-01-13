package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"golang.org/x/exp/rand"
)

var randomDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete random data",
	Run: func(cmd *cobra.Command, args []string) {
		op := NewJob("Delete Job", clients, runDeleteJob)
		op.Start()
		op.WaitAndLog()
	},
}

func runDeleteJob(runnerId int, op *Job) {
	op.WaitGroup.Add(1)
	databaseConfig, err := FetchDatabaseConfig(selectedDB)
	if err != nil {
		log.Fatalf("%v", err)
	}

	tableConfig, ok := databaseConfig.Tables[selectedTable]
	if !ok {
		log.Fatalf("Table %s config not found in database %s config", selectedTable, selectedDB)
	}

	db, err := NewDB(selectedDB)
	defer func() {
		_ = db.Close()
	}()

	if err != nil {
		log.Fatalf("Error while connecting to database: %v", err)
	}

	selectedColumn := ""
	columns := []string{}
	for _, col := range tableConfig.Columns {
		columns = append(columns, col.Name)
		if col.PrimaryKey {
			selectedColumn = col.Name
		}
	}

	if selectedColumn == "" {
		// pick a random column
		selectedColumn = columns[rand.Intn(len(columns))]
	}

	for {
		select {
		case <-op.StopChannel:
			fmt.Println("Stopping Delete job")
			op.WaitGroup.Done()
			return
		default:
			// Select first row and delete it
			selectQuery := fmt.Sprintf("SELECT %s FROM %s LIMIT 1", selectedColumn, selectedTable)
			var value interface{}
			err := db.QueryRow(selectQuery).Scan(&value)
			if err != nil {
				op.AddToFailed(runnerId, 1)
				continue
			}

			// Delete the selected row
			deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", selectedTable, selectedColumn)
			_, err = db.Exec(deleteQuery, value)
			op.AddToTotal(runnerId, 1)
			if err != nil {
				op.AddToFailed(runnerId, 1)
			}
		}
	}
}
