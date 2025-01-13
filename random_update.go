package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"golang.org/x/exp/rand"
)

var randomUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update random data",
	Run: func(cmd *cobra.Command, args []string) {
		op := NewJob("Update Job", clients, runUpdateJob)
		op.Start()
		op.WaitAndLog()
	},
}

func runUpdateJob(runnerId int, op *Job) {
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

	var selectedColumn *Column
	for _, col := range tableConfig.Columns {
		if col.PrimaryKey {
			selectedColumn = &col
		}
	}

	if selectedColumn == nil {
		// pick a random column
		selectedColumn = &tableConfig.Columns[rand.Intn(len(tableConfig.Columns))]
	}

	updatableColumn := &tableConfig.Columns[rand.Intn(len(tableConfig.Columns))]

	for {
		select {
		case <-op.StopChannel:
			fmt.Println("Stopping Update job")
			op.WaitGroup.Done()
			return
		default:
			// Select a random row to update
			selectQuery := fmt.Sprintf("SELECT %s FROM %s ORDER BY RAND() LIMIT 1", selectedColumn.Name, selectedTable)
			var value interface{}
			err := db.QueryRow(selectQuery).Scan(&value)
			if err != nil {
				op.AddToFailed(runnerId, 1)
				continue
			}

			// Generate new data for the selected column
			newValue := generate(updatableColumn.Generator) // Assuming we update the first column for simplicity

			// Update the selected row
			updateQuery := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s = ?", selectedTable, selectedColumn.Name, updatableColumn.Name)
			_, err = db.Exec(updateQuery, newValue, value)
			op.AddToTotal(runnerId, 1)
			if err != nil {
				op.AddToFailed(runnerId, 1)
			}
		}
	}
}
