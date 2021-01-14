package lua

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// Option can be provided to New() to customise initialization of Lua state.
type Option func(l *Lua) error

// Context sets a context to be used by the lua state.
func Context(ctx context.Context) Option {
	return func(l *Lua) error {
		if l.cancel != nil {
			l.cancel()
			l.ctx = nil
		}
		l.ctx, l.cancel = context.WithCancel(ctx)
		l.state.SetContext(l.ctx)
		return nil
	}
}

// Path appends the given directories to LUA_PATH in correct format.
func Path(dirs ...string) Option {
	const luaPathVar = "LUA_PATH"

	return func(_ *Lua) error {
		curValue := os.Getenv(luaPathVar)
		for i, dir := range dirs {
			dirs[i] = fmt.Sprintf("%s/?.lua", strings.TrimRight(dir, "/"))
		}
		newValue := curValue + ";" + strings.Join(dirs, ";")
		return os.Setenv(luaPathVar, newValue)
	}
}

// Globals sets all values in the map as global variables in the lua state.
// See Module() to create a module from the map.
func Globals(vals map[string]interface{}) Option {
	return func(l *Lua) error {
		for key, val := range vals {
			key = strings.TrimSpace(key)

			l.state.SetGlobal(key, luar.New(l.state, val))
		}
		return nil
	}
}

// Module defines a Lua module with public struct fields and methods exported.
// Use Globals() if you need to set the struct or the map as a single value.
func Module(name string, structOrMap interface{}) Option {
	name = strings.TrimSpace(name)
	exports := createExportsMap(structOrMap)
	funcs := getLFuncs(exports)

	return func(l *Lua) error {
		l.state.PreloadModule(name, func(state *lua.LState) int {
			mod := state.SetFuncs(state.NewTable(), funcs)

			for key, val := range exports {
				if _, found := funcs[key]; !found {
					state.SetField(mod, key, luar.New(l.state, val))
				}
			}

			state.Push(mod)
			return 1
		})
		return nil
	}
}

func createExportsMap(structVal interface{}) map[string]interface{} {
	rv := reflect.ValueOf(structVal)
	rt := rv.Type()

	res := map[string]interface{}{}
	switch rt.Kind() {
	case reflect.Struct:
		for i := 0; i < rt.NumField(); i++ {
			res[rt.Field(i).Name] = rv.Field(i).Interface()
		}

		for i := 0; i < rt.NumMethod(); i++ {
			res[rt.Method(i).Name] = rv.Method(i).Interface()
		}

	case reflect.Map:
		if rt.Key().Kind() != reflect.String {
			panic(fmt.Errorf("map key must be string, not '%s'", rt.Key().Kind()))
		}
		iter := rv.MapRange()
		for iter.Next() {
			res[iter.Key().String()] = iter.Value().Interface()
		}

	default:
		panic(fmt.Errorf("arg must be a struct or map, not '%s'", rt.String()))
	}

	return res
}

func getLFuncs(vals map[string]interface{}) map[string]lua.LGFunction {
	res := map[string]lua.LGFunction{}
	for k, v := range vals {
		if lgf, ok := v.(lua.LGFunction); ok {
			res[k] = lgf
		}
	}
	return res
}
