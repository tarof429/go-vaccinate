package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/tarof429/go-vaccinate/vaccinate"
)

var (
	s vaccinate.Simulator
)

func mainCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "load", Description: "Load configuration from ~/.vaccinate"},
		{Text: "run", Description: "Run simulation"},
		{Text: "quit", Description: "Quit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func load() error {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return s.Load(user.HomeDir)
}

func run() {
	s.Run()
}

// Print usage and exit
func usage() {
	fmt.Println("Usage: simulator [--terminal|--console]")
	fmt.Println()
	fmt.Println("--terminal will run the simulation and just print the results.")
	fmt.Println("--console will run the simulation and display the results using a plot and table")
	fmt.Println()
	os.Exit(1)
}

func main() {

	if len(os.Args) != 2 {
		usage()
	}

	arg := strings.Trim(os.Args[1], "")

	if arg == "--terminal" {
		s = &vaccinate.TerminalSimulator{}
	} else if arg == "--console" {
		s = &vaccinate.ConsoleSimulator{}
	} else {
		usage()
	}

	// Check the simulator type
	_, ok := s.(*vaccinate.TerminalSimulator)

	if ok {
		load()
		run()
	} else {
		for {
			fmt.Println("Please select command")
			t := prompt.Input("> ", mainCompleter)

			args := strings.Fields(t)
			if len(args) < 1 {
				continue
			}
			cmd := args[0]
			var err error

			if strings.ToLower(cmd) == "load" {
				err = load()
				if err != nil {
					fmt.Println("Error while attempting to read config file: " + err.Error())
				}
			} else if strings.ToLower(cmd) == "run" {
				run()
			} else if strings.ToLower(cmd) == "quit" {
				break
			}
		}
	}
}
