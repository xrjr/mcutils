package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/xrjr/mcutils/pkg/query"
)

type QueryFullCommand struct{}

func (QueryFullCommand) MinNumberOfArguments() int {
	return 2
}

func (QueryFullCommand) MaxNumberOfArguments() int {
	return 2
}

func (QueryFullCommand) Usage() string {
	return "<hostname> <port>"
}

func (cmd QueryFullCommand) Execute(params []string, jsonFormat bool) bool {
	port, err := strconv.Atoi(params[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid port.")
		return false
	}

	fs, err := query.QueryFull(params[1], port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s.\n", err.Error())
		return false
	}

	if jsonFormat {
		return cmd.jsonOutput(fs)
	}

	return cmd.basicOutput(fs)
}

func (QueryFullCommand) basicOutput(fs query.FullStat) bool {
	fmt.Println("Properties :")
	for k, v := range fs.Properties {
		fmt.Fprintf(os.Stderr, " - %s : %s\n", k, v)
	}

	fmt.Println("Online Players :")
	for _, v := range fs.OnlinePlayers {
		fmt.Fprintf(os.Stderr, " - %s\n", v)
	}

	return true
}

func (QueryFullCommand) jsonOutput(fs query.FullStat) bool {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(fs)
	if err != nil {
		return false
	}

	return true
}