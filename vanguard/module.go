package vanguard

import (
	"os"

	"github.com/mitchellh/mapstructure"
)

// An idempotent unit of work for orchestrating a provision.
type Module interface {
	// Unique name of the module.
	Name() string

	// Prepare the module for running. This will load the configuration,
	// do error checking, and perform any actions that can be pipelined.
	// The Prepare method should *never* perform an action. It should return
	// true if the Run method should be invoked and an error if there is some
	// problem.
	Prepare(config interface{}) (bool, error)

	// Invoke the module. This will always perform the action. The Prepare method
	// must be called before Prepare.
	Run() error
}

type file struct {
	Path string
	Mode os.FileMode
}

// The "file" module.
// The file module is used for creating and setting permissions
// for a file on disk.
func File() Module {
	return &file{}
}

func (file *file) Name() string {
	return "file"
}

func (file *file) Prepare(config interface{}) (bool, error) {
	if err := mapstructure.Decode(config, file); err != nil {
		return false, err
	}

	st, err := os.Stat(file.Path)
	if err != nil || st.Mode() != file.Mode {
		return true, nil
	}
	return false, nil
}

func (file *file) Run() error {
	f, err := os.Open(file.Path)
	if err != nil {
		f, err = os.Create(file.Path)
		if err != nil {
			return err
		}
	}
	return f.Chmod(file.Mode)
}
