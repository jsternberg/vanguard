package vanguard

import "sync"

type Runner struct {
	C chan *Task

	prepareCh chan *Task
	executeCh chan *Task

	wg sync.WaitGroup
}

func NewRunner() *Runner {
	r := &Runner{
		C:         make(chan *Task, 100),
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
			task.Run()
		}
	}
}

func (r *Runner) Wait() error {
	r.wg.Wait()
	return nil
}

func (r *Runner) Close() error {
	close(r.C)
	return nil
}
