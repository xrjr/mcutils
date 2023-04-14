package main

import (
	"fmt"
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
	return cmd.basicOutput()
}

func (VersionCommand) basicOutput() bool {
	bi, _ := debug.ReadBuildInfo()
	fmt.Printf("mcutils %s %s/%s\n", bi.Main.Version, runtime.GOOS, runtime.GOARCH)

	return true
}
