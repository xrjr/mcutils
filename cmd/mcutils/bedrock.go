package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xrjr/mcutils/pkg/bedrock"
)

type BedrockCommand struct{}

func (BedrockCommand) MinNumberOfArguments() int {
	return 2
}

func (BedrockCommand) MaxNumberOfArguments() int {
	return 2
}

func (BedrockCommand) Usage() string {
	return "<hostname> <port>"
}

func (BedrockCommand) Execute(params []string) bool {
	port, err := strconv.Atoi(params[1])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Invalid port.")
		return false
	}

	stat, err := bedrock.Stat(params[0], port)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s.\n", err.Error())
		return false
	}

	fmt.Printf("MOTD: %s\n", stat.MOTD)
	fmt.Printf("Version: %s\n", stat.Version)
	fmt.Printf("Protocol: %s\n", stat.Protocol)
	fmt.Printf("Num Players: %d\n", stat.Players)
	fmt.Printf("Max Players: %d\n", stat.MaxPlayers)
	return true
}
