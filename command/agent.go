package command

import (
	"fmt"
	"net"

	"github.com/jsternberg/vanguard/vanguard"
)

type AgentCommand struct{}

func (c *AgentCommand) Help() string {
	return `usage: vanguard agent`
}

func (c *AgentCommand) Run(args []string) int {
	l, err := net.Listen("tcp", ":8700")
	if err != nil {
		fmt.Println(err)
		return 1
	}

	server := vanguard.NewServer()
	if err := server.Serve(l); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

func (c *AgentCommand) Synopsis() string {
	return "Runs the vanguard agent"
}
