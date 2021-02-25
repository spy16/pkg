package rabbitmq

import (
	"context"
	"errors"
	"sync"

	"github.com/streadway/amqp"
)

// Streamer implements a RabbitMQ streamer. Messages are streamed to the Process
// function. Process is responsible for Ack/Nack.
type Streamer struct {
	Addr     string                `json:"addr"`
	Queue    string                `json:"queue"`
	Consumer string                `json:"consumer"`
	Workers  int                   `json:"workers"`
	AutoAck  bool                  `json:"auto_ack"`
	Process  func(d amqp.Delivery) `json:"-"`
}

func (st *Streamer) Run(ctx context.Context) error {
	conn, err := amqp.Dial(st.Addr)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func() { _ = ch.Close() }()

	msgCh, err := ch.Consume(st.Queue, st.Consumer, st.AutoAck,
		false, false, false, nil)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < st.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			st.runWorker(ctx, msgCh)
		}()
	}
	wg.Wait()

	return ctx.Err()
}

func (st *Streamer) runWorker(ctx context.Context, msgCh <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return

		case d, more := <-msgCh:
			if !more {
				return
			}
			st.Process(d)
		}
	}
}

func (st *Streamer) init() error {
	if st.Addr == "" {
		return errors.New("addr cannot be empty")
	}

	if st.Queue == "" {
		return errors.New("queue cannot be empty")
	}

	if st.Workers == 0 {
		st.Workers = 1
	}

	return nil
}
