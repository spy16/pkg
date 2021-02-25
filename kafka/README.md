# Kafka

## Streamer

```go
package main

import (
	"context"
	"log"

	"github.com/spy16/pkg/kafka"
)

func main() {
	ks := kafka.Streamer{
		Workers:       10,
		Topic:         "events",
		Servers:       "localhost:9092",
		ConsumerGroup: "pkg-simple-consumer",
		Apply: func(ctx context.Context, key, val []byte) error {
			log.Printf("%s = %s", string(key), string(val))
			return nil
		},
	}

	if err := ks.Run(context.Background()); err != nil {
		log.Fatalf("streamer exited: %v", err)
	}
}
```
