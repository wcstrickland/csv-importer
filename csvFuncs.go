package main

import (
	"fmt"
	"strconv"
)

//func getDBType() string {
//	var DBType string
//	fmt.Println("what kind of Database?")
//	fmt.Scanln(&DBType)
//}

// getTableName takes no arguments and returns a string
// getTableName takes user input to assign a variable for table name
func getTableName() string {
	var tableName string
	fmt.Println("what should the table be named?")
	fmt.Scanln(&tableName)
	fmt.Println("your table is named:", tableName)
	return tableName
}

// getUserChoice takes a map[int]string and returns string
// getUserChoice takes user input and parses the string to int
// to use in the map to return the valid choice or reprompt the user if the
// choice is invalid
func getUserChoice(choice string, validChoices map[int]string) string {
	var userChoice string
	var parsedChoice int
	fmt.Printf("What is the type of %s\n", choice)
	fmt.Println("The valid choices are:")
	for i, v := range validChoices { //TODO map prints out of order. https://stackoverflow.com/questions/12108215/golang-map-prints-out-of-order
		fmt.Println(i, ":", v)
	}
	fmt.Println("Please input your choice")
	fmt.Scanln(&userChoice)
	parsedChoice, _ = strconv.Atoi(userChoice)    // converts string input to int
	validChoice, ok := validChoices[parsedChoice] // stores choice
	// loop based on validity of input
	for {
		if ok == false {
			fmt.Printf("What is the type of %s\n", choice)
			fmt.Println("Invalid choice:")
			fmt.Println("The valid choices are:")
			for i, v := range validChoices {
				fmt.Println(i, ":", v)
			}
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

// parseValueByChoice takes a string and returns interface{}
// the string input is evaluated
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
