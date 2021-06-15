package main

import (
	"fmt"
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

// parseValueByChoice takes a string and returns interface{}
// the string input is evaluated and each case performs appropriate strconv.Method()
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
