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

func (cmd PingCommand) Execute(params []string, jsonFormat bool) bool {
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

	if jsonFormat {
		return cmd.jsonOutput(properties, latency)
	}

	return cmd.basicOutput(properties, latency)
}

func (PingCommand) basicOutput(properties ping.JSON, latency int) bool {
	jsonProperties, err := json.MarshalIndent(properties, "", "\t")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s.\n", err.Error())
		return false
	}

	fmt.Println("Properties :", string(jsonProperties))
	fmt.Printf("Latency : %d ms\n", latency)

	return true
}

func (PingCommand) jsonOutput(properties ping.JSON, latency int) bool {
	res := struct {
		Properties ping.JSON `json:"properties"`
		Latency    int       `json:"latency"`
	}{
		Properties: properties,
		Latency:    latency,
	}

	encoder := json.NewEncoder(os.Stdout)
	err := encoder.Encode(res)

	if err != nil {
		return false
	}

	return true
}
