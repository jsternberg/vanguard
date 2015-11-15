package vanguard

import (
	"io"
	"io/ioutil"
	"sync"
)

type Runner struct {
	C chan *Task
	w io.Writer

	prepareCh chan *Task
	executeCh chan *Task

	wg sync.WaitGroup
}

// Instantiates a new runner and starts the underlying channels for reading
// from the task channel and sending tasks to the prepare and execution
// pipelines.
func NewRunner(w io.Writer) *Runner {
	if w == nil {
		w = ioutil.Discard
	}

	r := &Runner{
		C:         make(chan *Task, 100),
		w:         w,
		prepareCh: make(chan *Task, 100),
		executeCh: make(chan *Task, 100),
	}
	r.wg.Add(3)
	go r.read()
	go r.prepare()
	go r.execute()
	return r
}

func (r *Runner) read() {
	defer r.wg.Done()
	for task := range r.C {
		r.prepareCh <- task
		r.executeCh <- task
	}
	close(r.prepareCh)
	close(r.executeCh)
}

func (r *Runner) prepare() {
	defer r.wg.Done()
	for task := range r.prepareCh {
		task.Prepare()
	}
}

func (r *Runner) execute() {
	defer r.wg.Done()
	for task := range r.executeCh {
		changed, err := task.Prepare()
		if changed && err == nil {
			task.Run(r.w)
		}
	}
}

// Waits for all tasks to finish before returning. You must call
// Close before calling Wait.
func (r *Runner) Wait() error {
	r.wg.Wait()
	return nil
}

// Closes the input channel for new tasks.
func (r *Runner) Close() error {
	close(r.C)
	return nil
}
