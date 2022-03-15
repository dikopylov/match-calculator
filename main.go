package main

import (
	"flag"
	"occurrence-calculator/internal/cli/commands"
	"occurrence-calculator/internal/model/infrastructure/cli"
	"os"
	"runtime"
)

var sourceFlag cli.Source

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	registerFlags()

	command := commands.NewMatchCounterCommand(os.Stdin, os.Stdout)
	command.AddFlag(cli.FlagType, sourceFlag)
	command.Execute()
}

func registerFlags() {
	flag.Var(&sourceFlag, cli.FlagType, "Type of source: url or file")

	flag.Parse()
}
