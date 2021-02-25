# RabbitMQ

## Streamer

```go
package main

import (
	"context"
	"log"

	"github.com/spy16/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

func main() {
	rs := &rabbitmq.Streamer{
		Workers:  10,
		Addr:     "localhost:5672",
		Consumer: "pkg-consumer",
		AutoAck:  false,
		Queue:    "events",
		Process: func(d amqp.Delivery) {
			log.Printf("message: %s", string(d.Body))
			_ = d.Ack(false)
		},
	}

	if err := rs.Run(context.Background()); err != nil {
		log.Fatalf("streamer exited: %v", err)
	}
}
```
