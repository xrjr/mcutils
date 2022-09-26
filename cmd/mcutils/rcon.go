package main

import (
	"encoding/json"
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

func (cmd RconCommand) Execute(params []string, jsonFormat bool) bool {
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

	if jsonFormat {
		return cmd.jsonOutput(response)
	}

	return cmd.basicOutput(response)
}

func (RconCommand) basicOutput(response string) bool {
	fmt.Println(response)

	return true
}

func (RconCommand) jsonOutput(response string) bool {
	res := struct {
		Response string `json:"response"`
	}{
		Response: response,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(res)
	if err != nil {
		return false
	}

	return true
}
