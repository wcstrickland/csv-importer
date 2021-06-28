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
	"time"
)

func main() {
	// flag stuff
	cmdLineDB := flag.String("t", "", "selected database type")
	dbConnString := flag.String("c", "", "URI/DSN")
	quietFlag := flag.Bool("quiet", false, "suppress confirmation messages")
	maxConns := flag.Int("m", 0, "max number of connections: defaults to unlimited")
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
		db.SetMaxOpenConns(*maxConns)
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
		start := time.Now()
		createTableString := createQueryString(tableName, fieldTypes, newFirstLine)
		_, err = db.Exec(createTableString)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println("TABLE CREATED SUCCESSFULLY")
		//		if err := createTable(db, createTableString); err != nil {
		//			fmt.Println("error", err)
		//		}

		jobs := make(chan job)
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go insertWorker(i, db, jobs)
		}

		// READ THE LINES OF THE CSV
		insertLines(db, tableName, lenRecord, r, jobs)
		close(jobs)
		wg.Wait()
		stop := time.Since(start)
		fmt.Println("time taken: ", stop)
	}
}
