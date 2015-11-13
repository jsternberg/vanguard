package vanguard

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	assert := assert.New(t)
	modules := []*mock{
		&mock{delay: 200 * time.Millisecond},
		&mock{delay: 300 * time.Millisecond},
		&mock{delay: 100 * time.Millisecond},
	}

	runner := NewRunner()
	for _, m := range modules {
		runner.C <- NewTask(m, nil)
	}
	runner.Close()

	runner.Wait()
	for _, m := range modules {
		assert.Equal(1, m.count)
	}
}
