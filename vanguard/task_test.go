package vanguard

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mock struct {
	count int
	delay time.Duration
}

func (m *mock) Prepare(config interface{}) (bool, error) {
	m.count++
	return true, nil
}

func (m *mock) Run() error {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	return nil
}

func TestTask(t *testing.T) {
	assert := assert.New(t)

	var wg sync.WaitGroup
	module := &mock{}
	task := NewTask(module, nil)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task.Prepare()
		}()
	}
	wg.Wait()

	assert.Equal(1, module.count)
}
