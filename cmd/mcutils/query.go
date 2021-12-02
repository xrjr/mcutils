package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xrjr/mcutils/pkg/query"
)

type QueryCommand struct{}

func (QueryCommand) MinNumberOfArguments() int {
	return 3
}

func (QueryCommand) MaxNumberOfArguments() int {
	return 3
}

func (QueryCommand) Usage() string {
	return "<basic|full> <hostname> <port>"
}

func (QueryCommand) Execute(params []string) bool {
	if params[0] == "basic" {
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

		fmt.Printf("MOTD : %s\n", bs.MOTD)
		fmt.Printf("Game Type : %s\n", bs.GameType)
		fmt.Printf("Map : %s\n", bs.Map)
		fmt.Printf("Num Players : %d\n", bs.NumPlayers)
		fmt.Printf("Max Players : %d\n", bs.MaxPlayers)
		fmt.Printf("Host Port : %d\n", int(bs.HostPort))
		fmt.Printf("Host IP : %s\n", bs.HostIP)
	} else if params[0] == "full" {
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

		fmt.Println("Properties :")
		for k, v := range fs.Properties {
			fmt.Fprintf(os.Stderr, " - %s : %s\n", k, v)
		}

		fmt.Println("Online Players :")
		for _, v := range fs.OnlinePlayers {
			fmt.Fprintf(os.Stderr, " - %s\n", v)
		}
	} else {
		return false
	}
	return true
}
