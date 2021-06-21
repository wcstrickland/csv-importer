package main

import (
	//"context"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // this is done to make use of the drivers only
	_ "github.com/lib/pq"              // the underscore allows for import without explicit refrence
	"io"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"strings"
	"time"
)

var (
	sslMode, dbType, host, port, user, password, dbname string
	db                                                  *sql.DB
	err                                                 error
	csvLines                                            int
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
	cmdLineDB := flag.String("t", "", "selected database type")
	dbConnString := flag.String("c", "", "URI/DSN")
	quietFlag := flag.Bool("quiet", false, "suppress confirmation messages")
	flag.Parse()
	// validate that db type specified via command line is valid or blank
	isValidDBType := false
	for _, v := range validDBChoices {
		if *cmdLineDB == v || *cmdLineDB == "" {
			isValidDBType = true
			break
		}
	}
	if !isValidDBType {
		fmt.Println("invalid database type. The valid types are", validDBChoices)
		log.Fatalln()
	}

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
		if *cmdLineDB == "" {
			dbType = getUserChoice("database", validDBChoices)
			if !*quietFlag {
				fmt.Println("\nyour database type is:", dbType)
			}
		} else {
			dbType = fmt.Sprint(*cmdLineDB)
		}

		// CONNECTTODBTYPE HANDLES THE CONNECTION VIA SWITCH CASES FOR DIFFERENT DB TYPES
		if *dbConnString == "" {
			db, err = connectToDBtype(dbType)
		} else {
			db, err = sql.Open(dbType, *dbConnString)
		}
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

		// SANITIZE FIELD NAMES
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

		start := time.Now()

		// CREATE THE TABLE
		createTableString := createQueryString(tableName, fieldTypes, newFirstLine)
		if err := createTable(db, createTableString); err != nil {
			fmt.Println("error", err)
		}

		// PREPARE INSERT STATEMENT
		insertString := insertQueryString(tableName, newFirstLine)
		insertStmt, err := db.Prepare(insertString)
		if err != nil {
			log.Fatalln("ERROR MALFORMED INSERT STATEMENT:", err)
		}
		defer insertStmt.Close()

		// READ THE LINES OF THE CSV
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("error reading csv file:", err)
			}
			_, err = insertRow(insertStmt, record)
			if err != nil {
				fmt.Println("ERROR SOMETHING WENT WRONG INSERTING ROW:", err)
			}
		}
		stop := time.Since(start)
		fmt.Println(stop)
	}
}
