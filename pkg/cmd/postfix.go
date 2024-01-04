package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"io"
	"os"
	"time"
)

//			Colors for terminal prompt			//

var Red = color.New(color.FgRed)
var Black = color.New(color.BgBlack)
var boldRed = Red.Add(color.Bold)
var red = color.New(color.BgRed).SprintFunc()
var greenFunc = color.New(color.FgGreen).SprintFunc()

// Loader message to indicate work in progress
var calculationSpinner = wow.New(os.Stdout, spin.Get(spin.Dots), greenFunc(" Calculating ..."))
var clearingSpinner = wow.New(os.Stdout, spin.Get(spin.Dots), red(" Clearing session ..."))

// RunInterpreter starts calculating the expression
func RunInterpreter() {
	if len(os.Args) == 0 {
		help()
		return
	}

	// Get command line arguments
	args := os.Args[1:]

	session := &Session{}

	// Load previous session from the json file if one exists
	jsonFile, err := os.Open("session.json")
	if err != nil {
		session = NewSession()
	} else {
		time.Sleep(1 * time.Second)

		byteValue, _ := io.ReadAll(jsonFile)

		err = json.Unmarshal(byteValue, &session)
		if err != nil {
			panic("cannot unmarshal session")
		}
	}

	defer jsonFile.Close()

	/*
		Commands
	*/
	command := args[0]

	switch command {
	case "input":
		// a little spinner to indicate work in progress
		// sleeping for a second-two to imitate
		calculationSpinner.Start()
		time.Sleep(1 * time.Second)
		fmt.Println()
		session.Input(args[1:])
	case "stack":
		session.PrintStack()
	case "symbol_table":
		session.PrintSymbolTable()
	case "clear":
		// a little spinner to indicate work in progress
		clearingSpinner.Start()
		time.Sleep(1 * time.Second)
		session.Clear()
		return
	case "help":
		help()
		time.Sleep(1 * time.Second)
	default:
		boldRed.Println(args[0], "command is not found, try using help")
		return
	}

	// Write current session in json file.
	sessionBytes, _ := json.Marshal(session)
	err = os.WriteFile("session.json", sessionBytes, 0644)
	if err != nil {
		Red.Println("error writing to the file")
	}
}

// help prints the help menu
func help() {
	h := `
	Commands:
		$ postfix input 		- input a Postfix++ expression
		$ postfix stack 		- prints the current state of stack
		$ postfix symbol_table 	- prints the current state of symbol table
		$ postfix clear  		- clears the current session purging the variables
	Usage:
		$ postfix input <expression_item> <expression_item> ... <expression_item>
	Examples:
		$ postfix input 3 4 5 + \* && postfix print
	
		$ postfix input A 3 =
		$ postfix input B 3 =
		$ postfix input A B \*
`
	fmt.Println(h)
}
