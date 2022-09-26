package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/xrjr/mcutils/pkg/query"
)

type QueryBasicCommand struct{}

func (QueryBasicCommand) MinNumberOfArguments() int {
	return 2
}

func (QueryBasicCommand) MaxNumberOfArguments() int {
	return 2
}

func (QueryBasicCommand) Usage() string {
	return "<hostname> <port>"
}

func (cmd QueryBasicCommand) Execute(params []string, jsonFormat bool) bool {
	port, err := strconv.Atoi(params[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid port.")
		return false
	}

	bs, err := query.QueryBasic(params[1], port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s.\n", err.Error())
		return false
	}

	if jsonFormat {
		return cmd.jsonOutput(bs)
	}

	return cmd.basicOutput(bs)
}

func (QueryBasicCommand) basicOutput(bs query.BasicStat) bool {
	fmt.Printf("MOTD : %s\n", bs.MOTD)
	fmt.Printf("Game Type : %s\n", bs.GameType)
	fmt.Printf("Map : %s\n", bs.Map)
	fmt.Printf("Num Players : %d\n", bs.NumPlayers)
	fmt.Printf("Max Players : %d\n", bs.MaxPlayers)
	fmt.Printf("Host Port : %d\n", int(bs.HostPort))
	fmt.Printf("Host IP : %s\n", bs.HostIP)

	return true
}

func (QueryBasicCommand) jsonOutput(bs query.BasicStat) bool {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(bs)
	if err != nil {
		return false
	}

	return true
}
