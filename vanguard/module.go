package vanguard

import (
	"os"

	"github.com/mitchellh/mapstructure"
)

type Module interface {
	Prepare(config interface{}) (bool, error)

	Run() error
}

type file struct {
	Path string
	Mode os.FileMode
}

func File() Module {
	return &file{}
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
