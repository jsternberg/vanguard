package command

import (
	"fmt"
	"os"

	"github.com/jsternberg/vanguard/vanguard"
)

type ExecCommand struct{}

func (c *ExecCommand) Help() string {
	return `usage: vanguard exec <file>`
}

func (c *ExecCommand) Run(args []string) int {
	if len(args) != 1 {
		fmt.Println("must specify exactly 1 argument")
		return 1
	}

	f, err := os.Open(args[0])
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer f.Close()

	loader := &vanguard.YamlLoader{}
	plan, err := loader.Load(f)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	runner := vanguard.NewRunner(os.Stdout)
	for _, task := range plan.Tasks {
		runner.C <- task
	}
	runner.Close()

	runner.Wait()
	return 0
}

func (c *ExecCommand) Synopsis() string {
	return "Executes a build plan from YAML"
}
