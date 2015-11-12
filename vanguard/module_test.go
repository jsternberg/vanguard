package vanguard

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileModule(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir(os.TempDir(), "vanguard-test")
	if err != nil {
		assert.Fail(err.Error())
	}
	defer os.RemoveAll(dir)

	path := fmt.Sprintf("%s/file.txt", dir)

	module := File()
	changed, err := module.Prepare(map[string]interface{}{
		"path": path,
		"mode": 0644,
	})

	assert.True(changed)
	assert.NoError(err)

	if assert.NoError(module.Run()) {
		_, err = os.Stat(path)
		assert.NoError(err)
	}
}
