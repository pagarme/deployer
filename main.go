package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"

	_ "github.com/pagarme/deployer/deploy/nomad"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	commands := Commands()

	cli := &cli.CLI{
		Args:     args,
		Commands: commands,
		HelpFunc: cli.BasicHelpFunc("deployer"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
