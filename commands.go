package main

import (
	"github.com/jsternberg/vanguard/command"
	"github.com/mitchellh/cli"
)

var Commands = map[string]cli.CommandFactory{
	"exec": func() (cli.Command, error) {
		return &command.ExecCommand{}, nil
	},
}
