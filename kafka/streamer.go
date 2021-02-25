package kafka

import (
	"context"
	"fmt"
	"sync"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	"github.com/spy16/pkg/log"
)

// Streamer streams kafka messages, decodes and applies some function to it.
type Streamer struct {
	log.Logger

	Apply         ApplyFn `json:"-"`
	Servers       string  `json:"servers"`
	Topic         string  `json:"topic"`
	Workers       int     `json:"workers"`
	StartOffset   string  `json:"start_offset"`
	ConsumerGroup string  `json:"consumer_group"`
}

// ApplyFn can be set on a streamer to process each message.
type ApplyFn func(ctx context.Context, key, val []byte) error

// Run spans the workers that consume from kafka and apply the configured
// function to each message.
func (st *Streamer) Run(ctx context.Context) error {
	kConf, err := st.init()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for i := 0; i < st.Workers; i++ {
		con, err := kafka.NewConsumer(kConf)
		if err != nil {
			return fmt.Errorf("failed to create consumer: %v", err)
		}

		if err := con.Subscribe(st.Topic, nil); err != nil {
			return fmt.Errorf("failed to subscribe to '%s': %v", st.Topic, err)
		}

		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			if err := st.runWorker(ctx, id, con); err != nil {
				st.Warnf("worker %d exited due to error: %v", id, err)
			} else {
				st.Infof("worker %d finished successfully", id)
			}
		}(i)
	}
	wg.Wait()

	st.Infof("all workers exited, streamer shutting down")
	return nil
}

func (st *Streamer) runWorker(ctx context.Context, workerID int, con *kafka.Consumer) error {
	defer func() { _ = con.Close() }()

	for {
		select {
		case <-ctx.Done():
			return nil

		case e, ok := <-con.Events():
			if !ok {
				return fmt.Errorf("consumer channel closed")
			}
			st.handleEvent(ctx, workerID, e, con)
		}
	}
}

func (st *Streamer) init() (*kafka.ConfigMap, error) {
	if st.Workers == 0 {
		st.Workers = 1
	}
	return &kafka.ConfigMap{
		"bootstrap.servers":               st.Servers,
		"group.id":                        st.ConsumerGroup,
		"enable.auto.commit":              false,
		"auto.offset.reset":               st.StartOffset,
		"socket.keepalive.enable":         true,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"enable.partition.eof":            true,
	}, nil
}

func (st *Streamer) handleEvent(ctx context.Context, workerID int, ev kafka.Event, con *kafka.Consumer) {
	switch e := ev.(type) {
	case kafka.AssignedPartitions:
		_ = con.Assign(e.Partitions)
		st.Debugf("partition %s assigned to %d", e.Partitions, workerID)

	case kafka.RevokedPartitions:
		_ = con.Unassign()
		st.Debugf("partition %s revoked from %d", e.Partitions, workerID)

	case *kafka.Message:
		if err := st.Apply(ctx, e.Key, e.Value); err != nil {
			st.Errorf("apply failed (partition=%s): %v", e.TopicPartition, err)
			return
		}
		_, _ = con.CommitMessage(e)

	case kafka.PartitionEOF:
		st.Infof("reached EOF of partition=%s", e)

	case kafka.Error:
		st.Warnf("got error from kafka: %v", e)
	}
}
