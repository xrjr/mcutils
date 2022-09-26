package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/xrjr/mcutils/pkg/ping"
)

type PingCommand struct{}

func (PingCommand) MinNumberOfArguments() int {
	return 2
}

func (PingCommand) MaxNumberOfArguments() int {
	return 2
}

func (PingCommand) Usage() string {
	return "<hostname> <port>"
}

func (PingCommand) Execute(params []string, jsonFormat bool) bool {
	port, err := strconv.Atoi(params[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid port.")
		return false
	}

	properties, latency, err := ping.Ping(params[0], port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s.\n", err.Error())
		return false
	}

	jsonProperties, err := json.MarshalIndent(properties, "", "\t")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s.\n", err.Error())
		return false
	}

	fmt.Println("Properties :", string(jsonProperties))
	fmt.Printf("Latency : %d ms\n", latency)
	return true
}
