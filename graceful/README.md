# graceful

[![GoDoc](https://godoc.org/github.com/spy16/canister/graceful?status.svg)](https://godoc.org/github.com/spy16/canister/graceful)
[![Go Report Card](https://goreportcard.com/badge/github.com/spy16/canister/graceful)](https://goreportcard.com/report/github.com/spy16/canister/graceful)

A wrapper around `http.Server` with Graceful-Shutdown enabled

## Usage

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/spy16/canister/graceful"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(wr http.ResponseWriter, req *http.Request) {
		time.Sleep(5 * time.Second)
		json.NewEncoder(wr).Encode(map[string]string{
			"status": "ok",
		})
	})
	srv := graceful.NewServer(mux, syscall.SIGINT, syscall.SIGTERM)
	srv.Addr = ":8080"
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
```