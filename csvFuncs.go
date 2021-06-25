package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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
func getSqlInfo() (string, string, string, string, string, string) {
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
	fmt.Println("\nSSL mode?(enable or disable)")
	fmt.Scanln(&sslMode)
	return host, user, password, dbname, port, sslMode
}

// connectToDBtype takes a string representing a type of db and returns a *sql.DB and an error
// the function handles returning a db connection for multiple db types
func connectToDBtype(dbType string) (*sql.DB, error) {
	switch dbType {
	case "postgres":
		db, err := connectPostgres()
		return db, err
	case "mysql":
		db, err := connectMysql()
		return db, err
	case "sqlite":
		db, err := connectSqlite()
		return db, err
	}
	return db, err
}

func connectSqlite() (*sql.DB, error) {
	sqliteFileName := ""
	fmt.Println("\nwhat SQLite file do you want to use?")
	fmt.Println("If your file is outside of this directory Please provide an absolute path to the file:\n")
	fmt.Scanln(&sqliteFileName)
	sqliteFileName = fmt.Sprintf("%s.db", sqliteFileName)
	if _, err = os.Open(sqliteFileName); err != nil {
		if errors.Is(err, fs.ErrNotExist) { // os.O_Open|os.O_Create?
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

func connectPostgres() (*sql.DB, error) {
	host, user, password, dbname, port, sslMode := getSqlInfo()
	psqlInfoMap := map[string]string{
		"host":     fmt.Sprintf("host=%s", host),
		"user":     fmt.Sprintf("user=%s", user),
		"password": fmt.Sprintf("password=%s", password),
		"dbname":   fmt.Sprintf("dbname=%s", dbname),
		"port":     fmt.Sprintf("port=%s", port),
		"sslmode":  fmt.Sprintf("sslmode=%s", sslMode),
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
}

func connectMysql() (*sql.DB, error) {
	host, user, password, dbname, _, _ := getSqlInfo()
	mysqlInfo := fmt.Sprintf("%s:%s@(%s)/%s", user, password, host, dbname)
	db, err = sql.Open("mysql", mysqlInfo)
	return db, err
}

func createQueryString(tableName string, fieldTypes, newFirstLine []string) string {
	createQueryString := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(", tableName)
	for i := 0; i < len(fieldTypes); i++ {
		createQueryString += fmt.Sprintf("%s %s, ", newFirstLine[i], fieldTypes[i])
	}
	createQueryString = strings.TrimSuffix(createQueryString, ", ")
	createQueryString += ")"
	return createQueryString
}

func createTable(db *sql.DB, query string) error {
	// time out context
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	fmt.Println("TABLE CREATED SUCCESSFULLY")
	return nil
}

func qString(tableName string, newFirstLine []string) string {
	xs := make([]string, 4)
	xs[0] = fmt.Sprintf("INSERT INTO %s ", tableName)
	xs[1] = "VALUES ("
	ph := strings.Repeat("?, ", len(newFirstLine))
	xs[2] = strings.TrimSuffix(ph, ", ")
	xs[3] = ")"
	return strings.Join(xs, " ")
}

func insertRow(db *sql.DB, query string, record []string) (sql.Result, error) {
	convertedRow := make([]interface{}, len(record))
	for i, v := range record {
		convertedRow[i] = v
	}
	result, err := db.Exec(query, convertedRow...)
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	return result, err
}

// unsafe query subject to sql injection
//func injectQueryString(queryPrefix string, curLine []string) string {
//	xs := make([]string, 3)
//	var vals strings.Builder
//	for _, v := range curLine {
//		fmt.Fprintf(&vals, "'%s', ", v)
//	}
//	str1 := vals.String()
//	str1 = str1[:vals.Len()-2]
//	xs[0] = queryPrefix
//	xs[1] = str1
//	xs[2] = ")"
//	return strings.Join(xs, " ")
//}
