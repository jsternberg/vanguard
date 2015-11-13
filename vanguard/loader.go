package vanguard

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Loader interface {
	Load(r io.Reader) (*Plan, error)
}

type YamlLoader struct{}

func (l *YamlLoader) Load(r io.Reader) (*Plan, error) {
	in, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	raw := rawPlan{}
	if err := yaml.Unmarshal(in, &raw); err != nil {
		return nil, err
	}

	tasks := make([]*Task, len(raw.Tasks))
	for i, rawTask := range raw.Tasks {
		if len(rawTask) != 1 {
			return nil, fmt.Errorf("missing action")
		}

		var module Module
		for k, v := range rawTask {
			switch k {
			case "file":
				module = File()
			default:
				return nil, fmt.Errorf("unknown module: %s", k)
			}

			tasks[i] = NewTask(module, v)
			break
		}
	}

	return &Plan{
		Tasks: tasks,
	}, nil
}
