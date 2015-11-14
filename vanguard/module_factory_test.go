package vanguard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuleFactory(t *testing.T) {
	assert := assert.New(t)

	mf := &moduleFactory{
		builtins: make(map[string]ModuleFunc),
	}
	mf.Register(func() Module { return &mock{} })

	module, err := mf.New("mock")
	if assert.NoError(err) {
		assert.IsType(&mock{}, module)
	}

	module, err = mf.New("error")
	assert.Nil(module)
	assert.Error(err)
}

func TestBuiltins(t *testing.T) {
	assert := assert.New(t)

	if module, err := DefaultModuleFactory.New("file"); assert.NoError(err) {
		assert.IsType(File(), module)
	}
}
