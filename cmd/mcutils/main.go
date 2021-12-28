package main

import (
	"fmt"
	"os"
)

type Command interface {
	MinNumberOfArguments() int
	MaxNumberOfArguments() int
	Execute(params []string) bool
	Usage() string
}

var (
	commands map[string]Command = map[string]Command{
		"ping":              PingCommand{},
		"query":             QueryCommand{},
		"rcon":              RconCommand{},
		"ping-legacy":       PingLegacyCommand{},
		"ping-legacy-1.6.4": PingLegacy1_6_4Command{},
	}
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage : %s <command> <params...>\n", os.Args[0])
		return
	}

	command, ok := commands[os.Args[1]]

	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command %s. Existing commands are :\n", os.Args[1])

		for k := range commands {
			fmt.Fprintf(os.Stderr, " - %s\n", k)
		}

		return
	}

	argsNumber := len(os.Args) - 2
	if argsNumber < command.MinNumberOfArguments() || argsNumber > command.MaxNumberOfArguments() {
		fmt.Fprintf(os.Stderr, "Invalid number of arguments (min=%d, max=%d).\n", command.MinNumberOfArguments(), command.MaxNumberOfArguments())
		fmt.Fprintf(os.Stderr, "Usage : %s %s %s\n", os.Args[0], os.Args[1], command.Usage())
		return
	}

	if !command.Execute(os.Args[2:]) {
		fmt.Fprintf(os.Stderr, "Usage : %s %s %s\n", os.Args[0], os.Args[1], command.Usage())
	}
}
