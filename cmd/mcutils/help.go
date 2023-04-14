package main

import (
	"fmt"
	"os"
)

type HelpCommand struct{}

func (HelpCommand) MinNumberOfArguments() int {
	return 0
}

func (HelpCommand) MaxNumberOfArguments() int {
	return 0
}

func (HelpCommand) Usage() string {
	return ""
}

func (cmd HelpCommand) Execute(_ []string, jsonFormat bool) bool {
	return cmd.basicOutput()
}

func (HelpCommand) basicOutput() bool {
	fmt.Println("Existing commands are :")
	for k := range commands {
		fmt.Fprintf(os.Stderr, " - %s\n", k)
	}

	return true
}
