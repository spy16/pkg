# Lua

An easy-to-use wrapper for `gopher-lua` that provides easy interoperability between Lua and Go.

## Usage

```go
package main

import "github.com/spy16/pkg/lua"

func main() {
	luaState, _ := lua.New(
		lua.Path("/Users/bob/lua-lib"),
		lua.Module("http", httpClient{}),
	)

	_ = luaState.Execute(`print("hello")`)
}
```
