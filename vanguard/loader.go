package vanguard

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

// An interface for loading a plan from a reader.
// The underlying parser of the reader is up to the specific
// Loader implementation.
type Loader interface {
	// Loads the plan from the passed in reader.
	// Returns an error if there's an underlying problem
	// with the reader or a valid plan could not be loaded.
	Load(r io.Reader) (*Plan, error)
}

// Loads a plan from a YAML document.
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
		task := &Task{}
		metadata := mapstructure.Metadata{}
		config := mapstructure.DecoderConfig{
			Result:   task,
			Metadata: &metadata,
		}
		decoder, err := mapstructure.NewDecoder(&config)
		if err != nil {
			return nil, err
		}

		if err := decoder.Decode(rawTask); err != nil {
			return nil, err
		}

		if len(metadata.Unused) != 1 {
			return nil, fmt.Errorf("requires exactly 1 extra attribute as the action")
		}

		name := config.Metadata.Unused[0]
		task.module, err = DefaultModuleFactory.New(name)
		if err != nil {
			return nil, err
		}
		task.config = rawTask[name]

		tasks[i] = task
	}

	return &Plan{
		Tasks: tasks,
	}, nil
}
