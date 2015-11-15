package vanguard

import (
	"fmt"
	"io"
	"sync"
)

// A thread-safe wrapper around a module with extra metadata
// needed for all tasks.
type Task struct {
	Name string

	module Module
	config interface{}

	prepared bool
	lock     sync.Mutex
	changed  bool
	err      error
}

// Runs the Prepare method from the underlying module and passes
// in the task configuration. This will only invoke Prepare exactly once
// and is thread-safe.
func (t *Task) Prepare() (bool, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.prepared {
		t.changed, t.err = t.module.Prepare(t.config)
		t.prepared = true
	}
	return t.changed, t.err
}

// Runs the underlying module. You must call Prepare before invoking this method.
func (t *Task) Run(w io.Writer) error {
	fmt.Fprintf(w, "Running %s[%s]\n", t.module.Name(), t.Name)
	return t.module.Run()
}
