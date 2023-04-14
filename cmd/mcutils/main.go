package main

import (
	"flag"
	"fmt"
	"os"
)

type Command interface {
	MinNumberOfArguments() int
	MaxNumberOfArguments() int
	Execute(params []string, jsonFormat bool) bool
	Usage() string
}

var (
	commands map[string]Command = map[string]Command{
		"ping":              PingCommand{},
		"query-basic":       QueryBasicCommand{},
		"query-full":        QueryFullCommand{},
		"rcon":              RconCommand{},
		"ping-legacy":       PingLegacyCommand{},
		"ping-legacy-1.6.4": PingLegacy1_6_4Command{},
		"ping-bedrock":      PingBedrockCommand{},
		"version":           VersionCommand{},
	}
)

func main() {
	var jsonFormat *bool = flag.Bool("json", false, "")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "Usage : %s [--json] <command> <params...>\n", os.Args[0])

		fmt.Fprint(os.Stderr, "Existing commands are :\n")
		for k := range commands {
			fmt.Fprintf(os.Stderr, " - %s\n", k)
		}

		os.Exit(1)
		return
	}

	command, ok := commands[flag.Arg(0)]

	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command %s. Run '%s help' to see existing commands.\n", os.Args[1])

		os.Exit(1)
		return
	}

	commandArgsNumber := len(flag.Args()) - 1
	if commandArgsNumber < command.MinNumberOfArguments() || commandArgsNumber > command.MaxNumberOfArguments() {
		fmt.Fprintf(os.Stderr, "Invalid number of arguments (current=%d, min=%d, max=%d).\n", commandArgsNumber, command.MinNumberOfArguments(), command.MaxNumberOfArguments())
		fmt.Fprintf(os.Stderr, "Usage : %s %s %s\n", os.Args[0], os.Args[1], command.Usage())
		os.Exit(1)
		return
	}

	if !command.Execute(flag.Args()[1:], *jsonFormat) {
		fmt.Fprintf(os.Stderr, "Usage : %s [--json] %s %s\n", os.Args[0], flag.Arg(0), command.Usage())
		os.Exit(1)
	}
}
