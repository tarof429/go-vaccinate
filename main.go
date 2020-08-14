package main

import (
	"fmt"
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

func main() {

	for {
		fmt.Println("Please select command")
		t := prompt.Input("> ", mainCompleter)
		//fmt.Println("You selected " + t)
		args := strings.Fields(t)
		cmd := args[0]

		if strings.ToLower(cmd) == "load" {
			break
		} else if strings.ToLower(cmd) == "run" {
			vaccinate.Run()
		} else if strings.ToLower(cmd) == "quit" {
			break
		}
	}
}
