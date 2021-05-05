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

With Retries using the `retry` package:

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/spy16/pkg/kafka"
	"github.com/spy16/pkg/retry"
)

func main() {
	backOff := retry.ExpBackoff(1.3, 0, 10*time.Second)

	ks := kafka.Streamer{
		Workers:       10,
		Topic:         "events",
		Servers:       "localhost:9092",
		ConsumerGroup: "pkg-simple-consumer",
		Apply: func(ctx context.Context, key, val []byte) error {
			// instead of 100, use retry.Forever to do infinite retries.
			return retry.Retry(ctx, 100, backOff, func() error {
				log.Printf("%s = %s", string(key), string(val))
				return nil
			})
		},
	}

	if err := ks.Run(context.Background()); err != nil {
		log.Fatalf("streamer exited: %v", err)
	}
}
```
