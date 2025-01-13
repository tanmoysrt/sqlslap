package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB(databaseName string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.ConnectionInfo.User, config.ConnectionInfo.Password, config.ConnectionInfo.Host, config.ConnectionInfo.Port, databaseName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Verify the connection to the database is still alive.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewDBIgnoreError(databaseName string) *sql.DB {
	db, err := NewDB(databaseName)
	if err != nil {
		log.Fatalf("Error while connecting to database %s: %v", databaseName, err)
	}
	return db
}

func FetchDatabases() ([]string, error) {
	db, err := NewDB("mysql")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	data, err := db.Query("show databases")
	if err != nil {
		return nil, err
	}
	var databases []string
	for data.Next() {
		var database string
		err = data.Scan(&database)
		if err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	return databases, nil
}

func FetchTables(database string) ([]string, error) {
	db, err := NewDB(database)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	data, err := db.Query("show tables")
	if err != nil {
		return nil, err
	}
	var tables []string
	for data.Next() {
		var table string
		err = data.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func CreateDatabases() error {
	db, err := NewDB("mysql")
	if err != nil {
		return err
	}
	defer db.Close()
	createdDatabases, err := FetchDatabases()
	if err != nil {
		return err
	}
	for _, database := range config.Databases {

		if !contains(createdDatabases, database.Name) {
			log.Printf("Creating database %s\n", database.Name)
			_, err = db.Exec(fmt.Sprintf("create database %s", database.Name))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
