package vanguard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInventory(t *testing.T) {
	assert := assert.New(t)

	inventory := NewInventory()
	assert.Equal([]string{}, inventory.Hosts())

	inventory.AddHost(Host{Name: "default", Addr: "127.0.0.1:8700"})
	assert.Equal([]string{"default"}, inventory.Hosts())

	host, ok := inventory.GetHost("default")
	if assert.True(ok) {
		assert.Equal("default", host.Name)
		assert.Equal("127.0.0.1:8700", host.Addr)
	}
}
