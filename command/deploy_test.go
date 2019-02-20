package command

import (
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	deployCommand := &DeployCommand{}
	result := deployCommand.Help()

	expectedResult := strings.TrimSpace(`
Usage:
	deployer command [options] <path>

Available Commands:
  deploy    Deploy an application using a configuration file

Options:
  --ref     Source Code Management hash to be fetched (default: master)
  --env     Environment to be used (default: main)
`)

	if result != expectedResult {
		t.Error("Expected:", expectedResult, "\nGot:", result)
	}
}

func TestSynopsis(t *testing.T) {
	deployCommand := &DeployCommand{}
	result := deployCommand.Synopsis()
	expectedResult := "Executes a deploy"

	if result != expectedResult {
		t.Error("Expected:", expectedResult, "\nGot:", result)
	}
}
