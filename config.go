package main

import (
	"errors"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ConnectionInfo ConnectionInfo `yaml:"connection_info"`
	Databases      []Database     `yaml:"databases"`
}

type ConnectionInfo struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Database struct {
	Name   string           `yaml:"name"`
	Schema string           `yaml:"schema"`
	Tables map[string]Table `yaml:"tables"`
}

type Table struct {
	Engine  string   `yaml:"engine"`
	Columns []Column `yaml:"columns"`
}

type Column struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	PrimaryKey bool   `yaml:"primary_key,omitempty"`
	Nullable   bool   `yaml:"nullable,omitempty"`
	Generator  string `yaml:"generator,omitempty"`
}

var config *Config

func init() {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error while reading file: %v", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func FetchDatabaseConfig(dbName string) (*Database, error) {
	for _, db := range config.Databases {
		if db.Name == dbName {
			return &db, nil
		}
	}
	return nil, errors.New("database config not found")
}
