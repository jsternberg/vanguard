package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

func realMain() int {
	app := cli.NewCLI("vanguard", Version)
	app.Args = os.Args[1:]
	app.Commands = Commands

	status, err := app.Run()
	if err != nil {
		fmt.Println(err)
		return 2
	}
	return status
}

func main() {
	os.Exit(realMain())
}
