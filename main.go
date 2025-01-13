package main

import (
	"log"
)

func main() {
	preVerifyDBConnections()
	rootCmd.Execute()
}

func preVerifyDBConnections() {
	// Ensure the credentials are correct
	db, err := NewDB("mysql")
	if err != nil {
		log.Fatalf("Error while connecting to database: %v", err)
	}
	// Create the databases if they don't exist
	err = CreateDatabases()
	if err != nil {
		log.Fatalf("Error while creating databases: %v", err)
	}
	db.Close()
}
