package command

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

type RunCommand struct {
	wg sync.WaitGroup
}

func (c *RunCommand) Help() string {
	return `usage: vanguard run <file>`
}

func (c *RunCommand) Run(args []string) int {
	if len(args) != 1 {
		fmt.Println("must specify exactly 1 argument")
		return 1
	}

	in, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Println(err)
		return 1
	}

	hosts := []string{"127.0.0.1:8700"}
	for _, host := range hosts {
		c.wg.Add(1)
		go c.provision(host, bytes.NewBuffer(in))
	}
	c.wg.Wait()
	return 0
}

func (c *RunCommand) provision(host string, in io.Reader) {
	defer c.wg.Done()

	url := fmt.Sprintf("http://%s/v1/provision", host)
	_, err := http.Post(url, "application/yaml", in)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (c *RunCommand) Synopsis() string {
	return "Runs a build plan on a remote vanguard agent"
}
