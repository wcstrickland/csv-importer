package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"strings"
)

// printSortedMap takes a map[int]string and has no return value
// the func allows a map whose keys are ints to be printed in order by those ints
func printSortedMap(m map[int]string) {
	var si []int
	for i := range m {
		si = append(si, i)
	}
	sort.Ints(si)
	for _, v := range si {
		fmt.Printf("%d: %s\n", v, m[v])
	}
}

// sanitize takes a string and returns a string
// for each prohbited char the string replaces all instances of said char
func sanitize(s string) string {
	prohibited := []string{";", ":", " ", "|", "-", "*", "/", "<", ">", ",", "=", "`", "~", "!", "?", "^", "(", ")"}
	for _, v := range prohibited {
		s = strings.ReplaceAll(s, v, "_")
	}
	return s
}

// getTableName takes no arguments and returns a string
// getTableName takes user input to assign a variable for table name
func getTableName() string {
	var tableName string
	fmt.Println("what should the table be named?")
	fmt.Scanln(&tableName)
	return tableName
}

// getUserChoice takes a string and a map[int]string and returns string
// getUserChoice takes user input and parses the string to int
// to use in the map to return the valid choice or reprompt the user if the
// choice is invalid
func getUserChoice(choice string, validChoices map[int]string) string {
	var userChoice string
	var parsedChoice int
	fmt.Printf("\nWhat is the type of %s\n", choice)
	fmt.Println("The valid choices are:")
	printSortedMap(validChoices)
	fmt.Println("Please input your choice")
	fmt.Scanln(&userChoice)
	parsedChoice, _ = strconv.Atoi(userChoice)    // converts string input to int
	validChoice, ok := validChoices[parsedChoice] // stores choice
	// loop based on validity of input
	for {
		if ok == false {
			fmt.Printf("\nWhat is the type of %s\n", choice)
			fmt.Println("Invalid choice:")
			printSortedMap(validChoices)
			fmt.Println("Please input your choice")
			fmt.Scanln(&userChoice)
			parsedChoice, _ = strconv.Atoi(userChoice)
			validChoice, ok = validChoices[parsedChoice]
		} else {
			break
		}
	}
	return validChoice
}

// getSqlInfo takes no arguments and returns a set of strings and ints used to construct a db driver string
func getSqlInfo() (string, string, string, string, string) {
	var host, user, password, dbname, port string
	fmt.Println("\nPlease enter host")
	fmt.Scanln(&host)
	fmt.Println("\nPlease enter port")
	fmt.Scanln(&port)
	fmt.Println("\nPlease enter user")
	fmt.Scanln(&user)
	fmt.Println("\nPlease enter password")
	fmt.Scanln(&password)
	fmt.Println("\nPlease enter dbname")
	fmt.Scanln(&dbname)
	return host, user, password, dbname, port
}

// connectToDBtype takes a string representing a type of db and returns a *sql.DB and an error
// the function handles returning a db connection for multiple db types
func connectToDBtype(dbtype string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	switch dbtype {
	case "Postgres":
		host, user, password, dbname, port := getSqlInfo()
		psqlInfoMap := map[string]string{
			"host":     fmt.Sprintf("host=%s", host),
			"user":     fmt.Sprintf("user=%s", user),
			"password": fmt.Sprintf("password=%s", password),
			"dbname":   fmt.Sprintf("dbname=%s", dbname),
			"port":     fmt.Sprintf("port=%s", port),
		}
		psqlInfo := ""
		for _, v := range psqlInfoMap {
			lastChar := v[len(v)-1:]
			if lastChar != "=" {
				psqlInfo += fmt.Sprint(v, " ")
			}
		}
		db, err = sql.Open("postgres", psqlInfo)
		return db, err
	case "MySQL":
		host, user, password, dbname, _ := getSqlInfo()
		mysqlInfo := fmt.Sprintf("%s:%s@(%s)/%s", user, password, host, dbname)
		db, err = sql.Open("mysql", mysqlInfo)
		return db, err
	case "SQLite":
		//var sqliteFileName string
		sqliteFileName := ""
		fmt.Println("\nwhat SQLite file do you want to use?")
		fmt.Println("If your file is outside of this directory Please provide an absolute path to the file:\n")
		fmt.Scanln(&sqliteFileName)
		sqliteFileName = fmt.Sprintf("%s.db", sqliteFileName)
		if _, err = os.Open(sqliteFileName); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				_, err = os.Create(sqliteFileName)
				fmt.Println("\nThe file you requested did not exist, but has now been created\n")
				if err != nil {
					fmt.Println("error:", err)
				}
			}
		}
		liteDsn := fmt.Sprintf("file:%s", sqliteFileName)
		db, err = sql.Open("sqlite", liteDsn)
		return db, err
	}
	return db, err
}

// parseValueByChoice takes a string and returns interface{}
// the string input is evaluated and each case performs appropriate strconv.Method()
// possibly not needed if the database query is a string and the type conversion is passed to the
// database to handle. Possibly useful outside of this project
func parseValueByChoice(choice string, value string) interface{} {
	switch choice {
	case "string":
		return value
	case "int":
		v1, _ := strconv.Atoi(value)
		return v1
	case "float":
		v2, _ := strconv.ParseFloat(value, 64)
		return v2
	}
	return "error"
}
