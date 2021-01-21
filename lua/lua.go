package lua

import (
	"context"
	"fmt"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// New returns a new Lua instance initialised. Options can be provided to
// bind globals, register modules etc.
func New(opts ...Option) (*Lua, error) {
	opts = append([]Option{Context(context.Background())}, opts...)
	l := &Lua{state: lua.NewState()}

	for _, opt := range opts {
		if err := opt(l); err != nil {
			return nil, err
		}
	}
	return l, nil
}

// Lua is a wrapper around lua-state and provides functions for managing state
// and executing lua code.
type Lua struct {
	ctx    context.Context
	cancel context.CancelFunc
	state  *lua.LState
}

// Execute the given lua script string. Use Call() for calling a function for
// result.
func (l *Lua) Execute(src string) error { return l.state.DoString(src) }

// ExecuteFile reads and executes the lua file. Use Call() for calling a function
// for result.
func (l *Lua) ExecuteFile(fileName string) error { return l.state.DoFile(fileName) }

// Call a lua function by its name. Args are automatically converted to
// appropriate types using the Luar library
func (l *Lua) Call(name string, args ...interface{}) (lua.LValue, error) {
	fn := l.state.GetGlobal(name)

	lfn, ok := fn.(*lua.LFunction)
	if !ok {
		return nil, fmt.Errorf("%s is not a function", name)
	}

	l.state.Push(lfn)
	for _, arg := range args {
		l.state.Push(luar.New(l.state, arg))
	}
	err := l.state.PCall(len(args), 1, nil)
	if err != nil {
		return nil, err
	}

	top := l.state.GetTop()
	retVal := l.state.Get(top)
	return retVal, nil
}

// Destroy releases all resources held by the lua state and marks the instance
// closed for usage.
func (l *Lua) Destroy() {
	if l.state == nil {
		return
	}
	l.cancel()
	l.state.Close()
	l.state = nil
}

// State returns the internal LState.
func (l *Lua) State() *lua.LState { return l.state }
