# errors

[![GoDoc](https://godoc.org/github.com/spy16/canister/errors?status.svg)](https://godoc.org/github.com/spy16/canister/errors)
[![Go Report Card](https://goreportcard.com/badge/github.com/spy16/canister/errors)](https://goreportcard.com/report/github.com/spy16/canister/errors)

Commonly used error definitions for golang.


## `Error` Type

`Error` type implements the `error` interface and provides certain
other methods and fields to describe the error more effectively.

```go
type Error struct {
    Code    int                     `json:"code"`
    Type    string                  `json:"type"`
    Details map[string]interface{}  `json:"details,omitempty"`
    Message string                  `json:"message,omitempty"`
}
```

See `error.go` file for more information.