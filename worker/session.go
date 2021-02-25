package worker

import (
	"context"
	"errors"
	"sync"

	"github.com/spy16/pkg/log"
)

type workerSession struct {
	log.Logger

	proc     Proc
	workers  int
	OnFinish func(job Job)
}

// Run spawns the workers to consume from the stream and executed the
// registered proc.
func (ws *workerSession) Run(ctx context.Context, stream <-chan Job) error {
	if err := ws.init(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &sync.WaitGroup{}
	for i := 0; i < ws.workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := ws.worker(ctx, stream); err != nil {
				ws.Debugf("worker-%d exited (cause: %s)", id, err)
			}
		}(i)
	}
	wg.Wait()

	return nil
}

func (ws *workerSession) worker(ctx context.Context, ch <-chan Job) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case j, ok := <-ch:
			if !ok {
				return errors.New("stream closed")
			}

			ws.processOne(ctx, j)
		}
	}
}

func (ws *workerSession) processOne(ctx context.Context, job Job) {
	// TODO: collect metrics?
	err := ws.proc.Exec(ctx, job)
	if err != nil {
		job.Error = err
	}
	ws.OnFinish(job)
	if job.Ack != nil {
		job.Ack(err)
	}
}

func (ws *workerSession) init() error {
	if ws.proc == nil {
		ws.proc = noOpProc
	}

	if ws.workers <= 0 {
		ws.workers = 1
	}

	return nil
}
