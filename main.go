package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // this is done to make use of the drivers only
	_ "github.com/lib/pq"              // the underscore allows for import without explicit refrence
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	sslMode, dbType, host, port, user, password, dbname string
	db                                                  *sql.DB
	err                                                 error
	csvLines                                            int
	wg                                                  sync.WaitGroup
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
	3: "DECIMAL",
}
var dbTypeChoices = map[string]map[int]string{
	"mysql":    validMysqlChoices,
	"postgres": validPostgresChoices,
	"sqlite":   validSqliteChoices,
}

func main() {
	// flag stuff
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
	for _, filename := range os.Args[1:] {
		if strings.HasPrefix(filename, "-") {
			continue
		}
		f, err := os.Open(filename)
		if err != nil {
			fmt.Println("!!!!!!!!!!!!!!!!!")
			fmt.Println("\nerror:", err)
			fmt.Println("\nMoving to next file\n")
			fmt.Println("!!!!!!!!!!!!!!!!!")
			continue
		}
		fmt.Println("\nThe currently selected file is:", filename)
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
		db.SetMaxOpenConns(0)
		db.SetMaxIdleConns(30)
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
		lenRecord := len(firstLine)
		if err != nil {
			fmt.Println("error reading CSV:", err)
		}
		var newFirstLine []string
		for _, fd := range firstLine { // sanitize headers
			newFirstLine = append(newFirstLine, sanitize(fd))
		}

		// GET TYPES OF HEADERS
		var fieldTypes []string
		for _, col := range newFirstLine {
			userChoice := getUserChoice(col, dbTypeChoices[dbType])
			fieldTypes = append(fieldTypes, userChoice)
		}

		// CREATE THE TABLE
		createTableString := createQueryString(tableName, fieldTypes, newFirstLine)
		if err := createTable(db, createTableString); err != nil {
			fmt.Println("error", err)
		}

		f.Close()
		// Split the file into multiples
		sliceOfFiles := splitFile(filename)

		start := time.Now()
		for i, file := range sliceOfFiles {
			// READ THE LINES OF THE CSV
			f, err := os.Open(file)
			r = csv.NewReader(f)
			if err != nil {
				fmt.Println("error processing a split file:", err)
			}
			if i == 0 {
				_, err = r.Read()
			}
			wg.Add(1)
			go func(db *sql.DB, tableName string, lenRecord int, r *csv.Reader) {
				insertLines(db, tableName, lenRecord, r)
				wg.Done()
			}(db, tableName, lenRecord, r)
		}
		wg.Wait()
		stop := time.Since(start)
		fmt.Println("time taken: ", stop)
	}
}
