package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"

	_ "github.com/pagarme/deployer/build/rocker"
	_ "github.com/pagarme/deployer/deploy/lambda"
	_ "github.com/pagarme/deployer/deploy/nomad"
	_ "github.com/pagarme/deployer/scm/git"
)

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {
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
