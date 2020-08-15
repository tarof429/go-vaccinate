package main

import (
	"fmt"
	"log"
	"os/user"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/tarof429/go-triage/vaccinate"
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

	return vaccinate.Load(user.HomeDir, attr)
}

func run(attr *vaccinate.PersonListAttributes) error {
	return vaccinate.Run(attr)
}

func main() {

	var attr vaccinate.PersonListAttributes

	for {
		fmt.Println("Please select command")
		t := prompt.Input("> ", mainCompleter)
		//fmt.Println("You selected " + t)
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
