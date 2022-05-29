package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xrjr/mcutils/pkg/bedrock"
)

type PingBedrockCommand struct{}

func (PingBedrockCommand) MinNumberOfArguments() int {
	return 2
}

func (PingBedrockCommand) MaxNumberOfArguments() int {
	return 2
}

func (PingBedrockCommand) Usage() string {
	return "<hostname> <port>"
}

func (PingBedrockCommand) Execute(params []string) bool {
	port, err := strconv.Atoi(params[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid port.")
		return false
	}

	pong, latency, err := bedrock.Ping(params[0], port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s.\n", err.Error())
		return false
	}

	fmt.Printf("Game Name : %s\n", pong.GameName)
	fmt.Printf("MOTD : %s\n", pong.MOTD)
	fmt.Printf("Protocol Version : %d\n", pong.ProtocolVersion)
	fmt.Printf("Minecraft Version : %s\n", pong.MinecraftVersion)
	fmt.Printf("Online Players : %d\n", pong.OnlinePlayers)
	fmt.Printf("Max Players : %d\n", pong.MaxPlayers)
	fmt.Printf("Server ID : %s\n", pong.ServerID)
	fmt.Printf("Map : %s\n", pong.Map)
	fmt.Printf("Game Mode : %s\n", pong.GameMode)
	fmt.Printf("IPv4 Port : %d\n", pong.IPv4Port)
	fmt.Printf("IPv6 Port : %d\n", pong.IPv6Port)
	fmt.Printf("Latency : %d ms\n", latency)
	return true
}
