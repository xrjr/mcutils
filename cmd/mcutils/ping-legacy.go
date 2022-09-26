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

func (cmd PingLegacyCommand) Execute(params []string, jsonFormat bool) bool {
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

	if jsonFormat {
		cmd.jsonOutput(infos, latency)
	}

	return cmd.basicOutput(infos, latency)
}

func (PingLegacyCommand) basicOutput(infos ping.LegacyPingInfos, latency int) bool {
	fmt.Printf("Protocol Version : %d\n", infos.ProtocolVersion)
	fmt.Printf("Minecraft Version : %s\n", infos.MinecraftVersion)
	fmt.Printf("MOTD : %s\n", infos.MOTD)
	fmt.Printf("Online Players : %d\n", infos.OnlinePlayers)
	fmt.Printf("Max Players : %d\n", infos.MaxPlayers)
	fmt.Printf("Latency : %d ms\n", latency)

	return true
}

func (PingLegacyCommand) jsonOutput(infos ping.LegacyPingInfos, latency int) bool {
	res := struct {
		ping.LegacyPingInfos
		Latency int `json:"latency"`
	}{
		LegacyPingInfos: infos,
		Latency:         latency,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(res)
	if err != nil {
		return false
	}

	return true
}
