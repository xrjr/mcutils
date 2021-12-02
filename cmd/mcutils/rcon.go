package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xrjr/mcutils/pkg/rcon"
)

type RconCommand struct{}

func (RconCommand) MinNumberOfArguments() int {
	return 4
}

func (RconCommand) MaxNumberOfArguments() int {
	return 4
}

func (RconCommand) Usage() string {
	return "<hostname> <port> <password> <command>"
}

func (RconCommand) Execute(params []string) bool {
	port, err := strconv.Atoi(params[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid port.")
		return false
	}

	response, err := rcon.Rcon(params[0], port, params[2], params[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s.\n", err.Error())
		return false
	}

	fmt.Println(response)
	return true
}
