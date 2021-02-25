package worker

import (
	"context"
)

// Run starts a worker and runs until context is cancelled or stream returns io.EOF.
func Run(ctx context.Context, stream <-chan Job, opts ...Option) error {
	ws := workerSession{}
	for _, opt := range withDefaults(opts) {
		if err := opt(&ws); err != nil {
			return err
		}
	}
	return ws.Run(ctx, stream)
}

// Proc represents the processor to be applied on each Job by the worker.
type Proc interface {
	Exec(ctx context.Context, job Job) error
}

// ProcFn is a adaptor type to implement Proc using simple Go func values.
type ProcFn func(ctx context.Context, job Job) error

func (pFn ProcFn) Exec(ctx context.Context, job Job) error { return pFn(ctx, job) }

var noOpProc = ProcFn(func(_ context.Context, _ Job) error {
	/* do nothing. */
	return nil
})

type Stream func(ctx context.Context) (<-chan Job, error)

// StreamFn turns a simple Go func into a Stream by invoking it in an infinite
// loop running on an independent GoRoutine. Stream exits on first error from
// the 'fn'.
func StreamFn(buffer int, fn func(ctx context.Context) (*Job, error)) Stream {
	return func(ctx context.Context) (<-chan Job, error) {
		ch := make(chan Job, buffer)

		go func() {
			defer close(ch)

			for {
				j, err := fn(ctx)
				if err != nil {
					return
				}

				select {
				case <-ctx.Done():
					return
				case ch <- *j:
				}
			}
		}()

		return ch, nil
	}
}
