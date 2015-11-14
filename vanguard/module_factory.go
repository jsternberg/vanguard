package vanguard

import "fmt"

var DefaultModuleFactory = &moduleFactory{
	builtins: make(map[string]ModuleFunc),
}

type ModuleFunc func() Module

type ModuleFactory interface {
	New(name string) (Module, error)

	Register(name string, f ModuleFunc)
}

type moduleFactory struct {
	builtins map[string]ModuleFunc
}

func (mf *moduleFactory) New(name string) (Module, error) {
	f, ok := mf.builtins[name]
	if !ok {
		return nil, fmt.Errorf("unknown module: %s", name)
	}
	return f(), nil
}

func (mf *moduleFactory) Register(name string, f ModuleFunc) {
	_, ok := mf.builtins[name]
	if ok {
		fmt.Printf("warning: registered %s multiple times\n", name)
	}
	mf.builtins[name] = f
}

func init() {
	DefaultModuleFactory.Register("file", File)
}
