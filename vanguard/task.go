package vanguard

import (
	"fmt"
	"io"
	"sync"
)

type Task struct {
	Name string

	module Module
	config interface{}

	prepared bool
	lock     sync.Mutex
	changed  bool
	err      error
}

func NewTask(module Module, config interface{}) *Task {
	return &Task{
		module: module,
		config: config,
	}
}

func (t *Task) Prepare() (bool, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.prepared {
		t.changed, t.err = t.module.Prepare(t.config)
		t.prepared = true
	}
	return t.changed, t.err
}

func (t *Task) Run(w io.Writer) error {
	fmt.Fprintf(w, "Running %s[%s]\n", t.module.Name(), t.Name)
	return t.module.Run()
}
