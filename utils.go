package main

import (
	"strconv"

	"github.com/go-faker/faker/v4"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func generate(generator string) string {
	switch generator {
	case "email":
		return faker.Email()
	case "name":
		return faker.Name()
	case "phone":
		return faker.Phonenumber()
	case "address":
		return faker.GetRealAddress().Address
	case "date":
		return faker.Date()
	case "time":
		return faker.TimeString()
	case "datetime":
		return faker.Timestamp()
	case "timestamp":
		return strconv.Itoa(int(faker.UnixTime()))
	case "uuid":
		return faker.UUIDHyphenated()
	case "word":
		return faker.Word()
	case "sentence":
		return faker.Sentence()
	case "lat":
		return strconv.Itoa(int(faker.Latitude()))
	case "lon":
		return strconv.Itoa(int(faker.Longitude()))
	default:
		return faker.Sentence()
	}
}

func generateData(columns []Column) []string {
	var data []string
	for _, column := range columns {
		data = append(data, generate(column.Generator))
	}
	return data
}

func generateInsertQueryWithData(tableName string, columns []Column) string {
	var query string
	query = "INSERT INTO " + tableName + " ("
	for i, column := range columns {
		query += column.Name
		if i != len(columns)-1 {
			query += ", "
		}
	}
	query += ") VALUES"

	query += " ("
	data := generateData(columns)
	for j, d := range data {
		query += strconv.Quote(d)
		if j != len(data)-1 {
			query += ", "
		}
	}
	query += ")"

	return query
}
