package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func main() {

	// get table name
	var tableName string
	fmt.Println("what should the table be named?")
	fmt.Scanln(&tableName)
	fmt.Println("your table is named:", tableName)

	// create valid choices map
	validChoices := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	// ask for user input
	fmt.Println(`The valid choices are:
        1: one
        2: two
        3: three
        Please input your choice:`)
	// parse user input for validation with map
	var userChoice string
	fmt.Scanln(&userChoice)
	parsedChoice, _ := strconv.Atoi(userChoice)
	validChoice, ok := validChoices[parsedChoice]
	// loop based on validity of input
	for {
		if ok == false {
			fmt.Println(`Invalid Choice:
The valid choices are:
        1: one
        2: two
        3: three
Please input your choice:`)
			fmt.Scanln(&userChoice)
			parsedChoice, _ = strconv.Atoi(userChoice)
			validChoice, ok = validChoices[parsedChoice]
		} else {
			break
		}
	}
	fmt.Println(validChoice)

	// open csv file
	f, err := os.Open("kidney.csv")
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
	fmt.Println(firstLine)

	// read lines temporarily using a loop to work with smaller numbers of lines
	for i := 0; i < 1; i++ {
		record, err := r.Read()
		if err != nil {
			fmt.Println("error:", err)
		}
		// range over an record to access colums
		for _, col := range record { // throw away the index
			fd, err := strconv.ParseFloat(col, 32)
			if err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Printf("%[1]T %[1]v\n", fd)
			}
		}
	}
}
