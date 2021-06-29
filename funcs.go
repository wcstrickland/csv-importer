package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"strings"
)

// printSortedMap takes a map[int]string and has no return value
// the function allows a map whose keys are integers to be printed in order by those integers
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
// for each prohibited char the string replaces all instances of said char
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
// to use in the map to return the valid choice or re prompt the user if the
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
		if !ok {
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

// getSqlInfo takes no arguments and returns a set of strings and integers  supplied by the user used to construct a db driver string
func getSQLInfo() (string, string, string, string, string, string) {
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

// connectSqlite attempts to make a connection to a database with user supplied information
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
	db, err = sql.Open("sqlite3", liteDsn)
	return db, err
}

// connectPostgres attempts to make a connection to a database with user supplied information
func connectPostgres() (*sql.DB, error) {
	host, user, password, dbname, port, sslMode := getSQLInfo()
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

// connectMysql attempts to make a connection to a database with user supplied information
func connectMysql() (*sql.DB, error) {
	host, user, password, dbname, _, _ := getSQLInfo()
	mysqlInfo := fmt.Sprintf("%s:%s@(%s)/%s", user, password, host, dbname)
	db, err = sql.Open("mysql", mysqlInfo)
	return db, err
}

// tableString(tableName string, fieldTypes, newFirstLine []string) string
// given a table name, the types of each field, and the label for each field the function returns an SQL statemetnt
// for creating a table in the form of a string
func tableString(tableName string, fieldTypes, newFirstLine []string) string {
	tableString := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(", tableName)
	for i := 0; i < len(fieldTypes); i++ {
		tableString += fmt.Sprintf("%s %s, ", newFirstLine[i], fieldTypes[i])
	}
	tableString = strings.TrimSuffix(tableString, ", ")
	tableString += ")"
	return tableString
}

// batchString(batchSize int, tableName string, lenRecord int) string
// creates an SQL multiple insert statement with given number of fields(lenRecord) and desired number of insertions(batchSize)
func batchString(batchSize int, tableName string, lenRecord int) string {
	phSlice := make([]string, batchSize)
	xs := make([]string, 3)
	xs[0] = fmt.Sprintf("INSERT INTO %s ", tableName)
	xs[1] = "VALUES "
	for i := 0; i < batchSize; i++ {
		ph := "("
		ph += strings.Repeat("?, ", lenRecord)
		ph = strings.TrimSuffix(ph, ", ")
		ph += "),"
		phSlice[i] = ph
	}
	phs := strings.Join(phSlice, " ")
	xs[2] = strings.TrimSuffix(phs, ",")
	return strings.Join(xs, " ")
}

// line counter takes r io.Reader and returns (int,error)
// efficient use of buffers and bytes.Count returns the number of lines with minimal resources
func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// insertLines takes db *sql.DB, tableName string, lenRecord int, r *csv.Reader, jobs chan<- job and has no return
// a csv.Reader is looped over. Values are collected into batches of 1000 and an insert query is generated
// the query and values comproise a job which is sent onto a channel to be consumed by workers
// if EOF is reached a query of appropriate size is constructed with values and sent. the function then exits
func insertLines(db *sql.DB, tableName string, lenRecord int, r *csv.Reader, jobs chan<- job) {
	for {
		vals := make([]interface{}, 1000*lenRecord)
		for i := 0; i < 1000; i++ {
			record, err := r.Read()
			if err == io.EOF {
				if i == 0 {
					return
				}
				vals = vals[:i*lenRecord]
				query := batchString(i, tableName, lenRecord)
				j := job{
					query: query,
					vals:  vals,
				}
				jobs <- j
				return
			}
			if i == 0 {
				for j, v := range record {
					vals[j] = v
				}
			} else {
				for j, v := range record {
					vals[(lenRecord*i)+j] = v
				}
			}
		}
		query := batchString(1000, tableName, lenRecord)
		j := job{
			query: query,
			vals:  vals,
		}
		jobs <- j
	}
}

// insert worker takes an id(int), db, <-chan job, chan<- int
// the worker pulls jobs from a channel and performs db insertions
// then sends an integer result out. +1 indicates a sucessful job -1 indicates the worker is closing
func insertWorker(lenRecord int, db *sql.DB, jobs <-chan job, results chan<- resultSignal) {
	for job := range jobs {
		_, err = db.Exec(job.query, job.vals...)
		if err != nil {
			fmt.Println("error at worker level", err)
			panic(err)
		}
		r := resultSignal{
			lines:  (len(job.vals) / lenRecord),
			signal: 1,
		}
		results <- r
	}
	r := resultSignal{
		lines:  0,
		signal: -1,
	}
	results <- r
}

// loading bar takes the string components desired to represent the bar, the desired width of the bar,
// the total work the bar is measuring progress of, a chanel to recieve instances of work completed, and total number of workers contributing to this work
// it renders a visual representation of progress on a work load
func loadingBar(bar, tip string, width int, totalWork int, workIn chan resultSignal, workers int) int {
	linesDone := 0
	workDone := 1
	doneSigs := workers // num of workers on workload
	var percentage float64
	for {
		if totalWork > 1000 {
			percentage = (float64(workDone) / float64(totalWork)) * 1000
		} else {
			percentage = (float64(workDone) / float64(totalWork)) //work done is divided by total work to achieve a percentage.
		}
		progress := percentage * float64(width) //this is multipled by the desired width
		rounded := int(progress)                //and rounded to represent the number of progress bars to be printed
		if rounded == 0 {
			rounded = 1
		}
		select { // if a value is recieved it is checked
		case v := <-workIn:
			if v.signal == 1 { // 1 indicates a successful job thus incrementing work done,
				workDone += v.signal
				linesDone += v.lines
			} else if v.signal == -1 { // -1 is a signal from the worker that it has closed thus decrementing the number of outstanding workers
				doneSigs--
			}
		default: // if no signal is recieved a display bar is rendered (this allows the chanel not to block)
			if doneSigs == 0 { // if work is finished (all workers have reported closing) the bar is completed and the function exits
				fmt.Printf("\r[%s%s%s]%d%%", strings.Repeat(bar, width+1), tip, strings.Repeat(" ", 1), 100)
				return linesDone
			} // otherwise the bar represents the ratio of work done to total work
			fmt.Printf("\r[%s%s%s]%2.f%%", strings.Repeat(bar, rounded), tip, strings.Repeat(" ", width-rounded), percentage*100)
		}
	}
}
