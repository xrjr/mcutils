package main

import (
	"encoding/json"
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

func (cmd PingBedrockCommand) Execute(params []string, jsonFormat bool) bool {
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

	if jsonFormat {
		return cmd.jsonOutput(pong, latency)
	}

	return cmd.jsonOutput(pong, latency)
}

func (PingBedrockCommand) basicOutput(pong bedrock.UnconnectedPong, latency int) bool {
	fmt.Printf("Game Name : %s\n", pong.GameName)
	fmt.Printf("MOTD : %s\n", pong.MOTD)
	fmt.Printf("Protocol Version : %d\n", pong.ProtocolVersion)
	fmt.Printf("Minecraft Version : %s\n", pong.MinecraftVersion)
	fmt.Printf("Online Players : %d\n", pong.OnlinePlayers)
	fmt.Printf("Max Players : %d\n", pong.MaxPlayers)
	fmt.Printf("Server ID : %s\n", pong.ServerID)
	fmt.Printf("Level Name : %s\n", pong.LevelName)
	fmt.Printf("Game Mode : %s\n", pong.GameMode)
	fmt.Printf("Game Mode (Numeric) : %d\n", pong.GameModeNumeric)
	fmt.Printf("IPv4 Port : %d\n", pong.IPv4Port)
	fmt.Printf("IPv6 Port : %d\n", pong.IPv6Port)
	fmt.Printf("Latency : %d ms\n", latency)

	return true
}

func (PingBedrockCommand) jsonOutput(pong bedrock.UnconnectedPong, latency int) bool {
	res := struct {
		bedrock.UnconnectedPong
		Latency int `json:"latency"`
	}{
		UnconnectedPong: pong,
		Latency:         latency,
	}

	encoder := json.NewEncoder(os.Stdout)
	err := encoder.Encode(res)

	if err != nil {
		return false
	}

	return true
}
