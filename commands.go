package main

import (
	"github.com/mitchellh/cli"

	"github.com/pagarme/deployer/command"
)

func Commands() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"deploy": func() (cli.Command, error) {
			return &command.DeployCommand{}, nil
		},
	}
}
