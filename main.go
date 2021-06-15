package main

import (
	"encoding/csv"
	"fmt"
	"github.com/jackc/pgx/4"
	"github.com/jackc/pgx/tree/master/pgxpool"
	"os"
)

func main() {

	tableName := getTableName()
	fmt.Println("\nYour table is named:", tableName)

	validDBChoices := map[int]string{
		1: "MySQL",
		2: "Postgress",
		3: "SQLServer",
		4: "SQLite",
	}

	dbType := getUserChoice("database", validDBChoices)
	fmt.Println("\nyour database type is:", dbType)

	// create valid choices map
	// TODO create this dynamically by db type
	validTypeChoices := map[int]string{
		1: "string",
		2: "int",
		3: "float",
	}

	// open csv file
	f, err := os.Open("sample.csv")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer f.Close()

	// make new csv.Reader from file
	r := csv.NewReader(f)

	// read first line for headers
	firstLine, err := r.Read()
	if err != nil {
		fmt.Println("error:", err)
	}

	var newFirstLine []string
	for _, fd := range firstLine {
		newFirstLine = append(newFirstLine, sanitize(fd))
	}

	// get types of headers
	var fieldTypes []string
	for _, col := range newFirstLine {
		userChoice := getUserChoice(col, validTypeChoices)
		fieldTypes = append(fieldTypes, userChoice)
	}
	fmt.Println(fieldTypes)

	// read lines temporarily using a loop to work with smaller numbers of lines
	for i := 0; i < 1; i++ {
		record, err := r.Read()
		if err != nil {
			fmt.Println("error:", err)
		}
		// range over an record to access colums
		var row []interface{}
		for i, col := range record {
			c := parseValueByChoice(fieldTypes[i], col)
			row = append(row, c)
		}
		fmt.Println(row) // this `row` is a slice with values of different types ready for insertion
	}
}
