package main

import (
	//"context"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // this is done to make use of the drivers only
	_ "github.com/lib/pq"              // the underscore allows for import without explicit refrence
	"log"
	_ "modernc.org/sqlite"
	"os"
	"strings"
	"time"
)

var (
	dbType, host, port, user, password, dbname string
	db                                         *sql.DB
	err                                        error
)

var validDBChoices = map[int]string{
	1: "mysql",
	2: "postgres",
	3: "sqlite",
}

var validMysqlChoices = map[int]string{
	1: "VARCHAR(60)",
	2: "INT",
	3: "FLOAT",
	4: "DECIMAL",
	5: "DATE",
	6: "TIME",
}
var validPostgresChoices = map[int]string{
	1: "varchar(60)",
	2: "integer",
	3: "float(12)",
	4: "date",
	5: "time(6)",
}
var validSqliteChoices = map[int]string{
	1: "TEXT",
	2: "INTEGER",
	3: "REAL",
}
var dbTypeChoices = map[string]map[int]string{
	"mysql":    validMysqlChoices,
	"postgres": validPostgresChoices,
	"sqlite":   validSqliteChoices,
}

func main() {
	quietFlag := flag.Bool("quiet", false, "suppress confirmation messages")
	flag.Parse()

	// LOOP OVER ALL COMAND LINE ARGUMENTS AND PERFORM THE PROGRAM ON EACH CSV FILE
	for _, v := range os.Args[1:] {
		if strings.HasPrefix(v, "-") {
			continue
		}
		f, err := os.Open(v)
		if err != nil {
			fmt.Println("!!!!!!!!!!!!!!!!!")
			fmt.Println("\nerror:", err)
			fmt.Println("\nMoving to next file\n")
			fmt.Println("!!!!!!!!!!!!!!!!!")
			continue
		}
		defer f.Close()
		fmt.Println("\nThe currently selected file is:", v)
		// make a csv Reader from the file
		r := csv.NewReader(f)

		//GET DB TYPE
		dbType := getUserChoice("database", validDBChoices)
		if !*quietFlag {
			fmt.Println("\nyour database type is:", dbType)
		}

		// CONNECTTODBTYPE HANDLES THE CONNECTION VIA SWITCH CASES FOR DIFFERENT DB TYPES
		db, err := connectToDBtype(dbType)
		if err != nil {
			fmt.Println("error:", err)
			panic(err)
		}
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		defer db.Close()

		// PING THE DB AND FATAL OUT IF THE CONNECTION IS NOT SUCCESSFUL
		err = db.Ping()
		if err != nil {
			fmt.Println("error:", err)
			log.Fatalln("\nYour database connection failed!\nThe program will now terminate")
		}
		if !*quietFlag {
			fmt.Println("\nyou connection status is: connected \n")
		}

		// GET THE TABLE NAME
		tableName := getTableName()
		if !*quietFlag {
			fmt.Println("\nYour table is named:", tableName)
		}

		// READ FIRST LINE FOR HEADERS
		firstLine, err := r.Read()
		if err != nil {
			fmt.Println("error:", err)
		}

		// sANITIZE FIELD NAMES
		var newFirstLine []string
		for _, fd := range firstLine {
			newFirstLine = append(newFirstLine, sanitize(fd))
		}

		// GET TYPES OF HEADERS
		var fieldTypes []string
		for _, col := range newFirstLine {
			userChoice := getUserChoice(col, dbTypeChoices[dbType])
			fieldTypes = append(fieldTypes, userChoice)
		}
		if !*quietFlag {
			fmt.Println(fieldTypes)
		}

		// CREATE THE TABLE
		createTableString := createQueryString(tableName, fieldTypes, newFirstLine)
		if err := createTable(db, createTableString); err != nil {
			fmt.Println("error", err)
		}
		// read lines temporarily using a loop to work with smaller numbers of lines
		//		for i := 0; i < 1; i++ {
		//			record, err := r.Read()
		//			if err != nil {
		//				fmt.Println("error:", err)
		//			}
		// range over an record to access colums
		//			var row []interface{}
		//			for i, col := range record {
		//				c := parseValueByChoice(fieldTypes[i], col)
		//				row = append(row, c)
		//			}
		//			if !*quietFlag {
		//				fmt.Println(row) // this `row` is []interface{} ready for insertion
		//			}
		//		}
	}
}
