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

// Runs the underlying module. This will invoke Prepare automatically
// and will only call the underlying module if it needs to be changed.
func (t *Task) Run(w io.Writer) error {
	changed, err := t.Prepare()
	if err != nil {
		fmt.Fprintf(w, "failed: %s[%s]\n  %s\n", t.module.Name(), t.Name, err)
		return err
	}

	if !changed {
		fmt.Fprintf(w, "ok: %s[%s]\n", t.module.Name(), t.Name)
		return nil
	}

	if err := t.module.Run(); err != nil {
		fmt.Fprintf(w, "failed: %s[%s]\n  %s\n", t.module.Name(), t.Name, err)
		return err
	}
	fmt.Fprintf(w, "changed: %s[%s]\n", t.module.Name(), t.Name)
	return nil
}
