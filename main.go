package main

import (
	"fmt"
	"log"
	"os/user"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/tarof429/go-vaccinate/vaccinate"
)

var (
	s    vaccinate.Simulator
	attr vaccinate.PersonListAttributes
)

func mainCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "load", Description: "Load configuration from ~/.vaccinate"},
		{Text: "run", Description: "Run simulation"},
		{Text: "quit", Description: "Quit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func load(attr *vaccinate.PersonListAttributes) error {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return s.Load(user.HomeDir, attr)
}

func run(attr *vaccinate.PersonListAttributes) error {

	return s.Run(attr)
}

func main() {

	s = vaccinate.ConsoleSimulator{}

	for {
		fmt.Println("Please select command")
		t := prompt.Input("> ", mainCompleter)

		args := strings.Fields(t)
		cmd := args[0]
		var err error

		if strings.ToLower(cmd) == "load" {
			err = load(&attr)
			if err != nil {
				fmt.Println("Error while attempting to read config file: " + err.Error())
			}
		} else if strings.ToLower(cmd) == "run" {
			err = run(&attr)

			if err != nil {
				fmt.Println(err.Error())
			}
		} else if strings.ToLower(cmd) == "quit" {
			break
		}
	}
}
