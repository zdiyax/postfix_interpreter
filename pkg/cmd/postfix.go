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
var blackFunc = color.New(color.BgBlack).SprintFunc()
var greenFunc = color.New(color.FgGreen).SprintFunc()

// Loader message to indicate work in progress
var calculationSpinner = wow.New(os.Stdout, spin.Get(spin.Dots), greenFunc(" Calculating ..."))
var clearingSpinner = wow.New(os.Stdout, spin.Get(spin.Dots), blackFunc(" Clearing session ..."))

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
		Black.Println("no session, opening a new one")
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
		session.Input(args[1:])
	case "print":
		// a little spinner to indicate work in progress
		calculationSpinner.Start()
		// sleeping for a second-two to imitate
		time.Sleep(2 * time.Second)
		session.Print()
	case "clear":
		// a little spinner to indicate work in progress
		clearingSpinner.Start()
		time.Sleep(1 * time.Second)
		session.Clear()
		return
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
		$ postfix print 		- prints the current state of stack
		$ postfix variables  	- prints current variables
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