package vanguard

import "fmt"

// The default module factory. All of the builtin
// modules are automatically registered with this factory.
var DefaultModuleFactory = NewModuleFactory()

type ModuleFunc func() Module

// A factory for instantiating new Module instances from
// a registry.
type ModuleFactory interface {
	// Instantiate the Module registered with the name.
	New(name string) (Module, error)

	// Register a Module creation function. The name is inferred
	// from the Module created from the ModuleFunc.
	Register(f ModuleFunc)
}

// Create a new module from the DefaultModuleFactory.
func NewModule(name string) (Module, error) {
	return DefaultModuleFactory.New(name)
}

// Register a module with the DefaultModuleFactory.
func RegisterModule(f ModuleFunc) {
	DefaultModuleFactory.Register(f)
}

type moduleFactory struct {
	builtins map[string]ModuleFunc
}

// Instantiate a new module factory with no default modules.
func NewModuleFactory() ModuleFactory {
	return &moduleFactory{
		builtins: make(map[string]ModuleFunc),
	}
}

func (mf *moduleFactory) New(name string) (Module, error) {
	f, ok := mf.builtins[name]
	if !ok {
		return nil, fmt.Errorf("unknown module: %s", name)
	}
	return f(), nil
}

func (mf *moduleFactory) Register(f ModuleFunc) {
	name := f().Name()
	_, ok := mf.builtins[name]
	if ok {
		fmt.Printf("warning: registered %s multiple times\n", name)
	}
	mf.builtins[name] = f
}

func init() {
	RegisterModule(File)
}
