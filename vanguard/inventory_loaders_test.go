package vanguard

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDynamicInventory(t *testing.T) {
	assert := assert.New(t)

	tmpfile, err := ioutil.TempFile(os.TempDir(), "vanguard-inventory")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	io.WriteString(tmpfile, `#!/usr/bin/python
print """{
	"default": [
		"default"
	]
}
"""`)
	tmpfile.Chmod(0700)
	tmpfile.Close()

	inventory, err := DynamicInventory(tmpfile.Name())
	if assert.NoError(err) {
		assert.Equal([]string{"default"}, inventory.Hosts())
	}
}
