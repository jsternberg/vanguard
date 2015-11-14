package main

import (
	"github.com/jsternberg/vanguard/command"
	"github.com/mitchellh/cli"
)

var Commands = map[string]cli.CommandFactory{
	"agent": func() (cli.Command, error) {
		return &command.AgentCommand{}, nil
	},
	"exec": func() (cli.Command, error) {
		return &command.ExecCommand{}, nil
	},
	"run": func() (cli.Command, error) {
		return &command.RunCommand{}, nil
	},
}
