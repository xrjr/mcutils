package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/xrjr/mcutils/pkg/ping"
)

type PingLegacyCommand struct{}

func (PingLegacyCommand) MinNumberOfArguments() int {
	return 2
}

func (PingLegacyCommand) MaxNumberOfArguments() int {
	return 2
}

func (PingLegacyCommand) Usage() string {
	return "<hostname> <port>"
}

func (PingLegacyCommand) Execute(params []string) bool {
	port, err := strconv.Atoi(params[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid port.")
		return false
	}

	infos, latency, err := ping.PingLegacy(params[0], port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s.\n", err.Error())
		return false
	}

	jsonProperties, err := json.MarshalIndent(infos, "", "\t")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s.\n", err.Error())
		return false
	}

	fmt.Println(string(jsonProperties))
	fmt.Printf("Latency : %d ms\n", latency)
	return true
}
