package main

import (
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

	fmt.Printf("Protocol Version : %d\n", infos.ProtocolVersion)
	fmt.Printf("Minecraft Version : %s\n", infos.MinecraftVersion)
	fmt.Printf("MOTD : %s\n", infos.MOTD)
	fmt.Printf("Online Players : %d\n", infos.OnlinePlayers)
	fmt.Printf("Max Players : %d\n", infos.MaxPlayers)
	fmt.Printf("Latency : %d ms\n", latency)
	return true
}
