package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
)

type VersionCommand struct{}

func (VersionCommand) MinNumberOfArguments() int {
	return 0
}

func (VersionCommand) MaxNumberOfArguments() int {
	return 0
}

func (VersionCommand) Usage() string {
	return ""
}

func (cmd VersionCommand) Execute(_ []string, jsonFormat bool) bool {

	if jsonFormat {
		return cmd.jsonOutput(version())
	}

	return cmd.basicOutput(version())
}

func (VersionCommand) basicOutput(version string) bool {
	fmt.Println("Version :", version)

	return true
}

func (VersionCommand) jsonOutput(version string) bool {
	res := struct {
		Version string `json:"version"`
	}{
		Version: version,
	}

	encoder := json.NewEncoder(os.Stdout)
	err := encoder.Encode(res)

	if err != nil {
		return false
	}

	return true
}

func version() string {
	bi, _ := debug.ReadBuildInfo()
	return "mcutils " + bi.Main.Version + " " + runtime.GOOS + "/" + runtime.GOARCH
}
