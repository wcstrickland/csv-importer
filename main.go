package main

import (
	//"context"
	//"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // this is done to make use of the drivers only
	_ "github.com/lib/pq"              // the underscore allows for import without explicit refrence
	"log"
	_ "modernc.org/sqlite"
	"os"
	"time"
)

var host, user, password, dbname string
var port int

func main() {
	validDBChoices := map[int]string{
		1: "MySQL",
		2: "Postgres",
		4: "SQLite",
	}
	dbType := getUserChoice("database", validDBChoices)
	fmt.Println("\nyour database type is:", dbType)

	// connectToDBtype handles the connection via switch cases for different db types
	db, err := connectToDBtype(dbType)
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()

	// Ping the db and Fatal out if the connection is not successful
	err = db.Ping()
	if err != nil {
		log.Fatalln("\nYour database connection failed!\nThe program will now terminate")
	}
	fmt.Println("\nyou connection status is: connected \n")

	// get the table name
	tableName := getTableName()
	fmt.Println("\nYour table is named:", tableName)

	// create valid choices map
	// TODO create this dynamically by db type
	validTypeChoices := map[int]string{
		1: "string",
		2: "int",
		3: "float",
	}

	// open csv file
	// TODO get this dynamically either via command line or by user input at runtime
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

	// sanitize field names
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
		fmt.Println(row) // this `row` is []interface{} ready for insertion
	}
}
