package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/xrjr/mcutils/pkg/ping"
)

type PingLegacy1_6_4Command struct{}

func (PingLegacy1_6_4Command) MinNumberOfArguments() int {
	return 2
}

func (PingLegacy1_6_4Command) MaxNumberOfArguments() int {
	return 2
}

func (PingLegacy1_6_4Command) Usage() string {
	return "<hostname> <port>"
}

func (PingLegacy1_6_4Command) Execute(params []string) bool {
	port, err := strconv.Atoi(params[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid port.")
		return false
	}

	infos, latency, err := ping.PingLegacy1_6_4(params[0], port)
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
