package worker

import "github.com/spy16/pkg/log"

// Option can be provided to Run() to customise run behaviour of the worker.
type Option func(ws *workerSession) error

func withDefaults(opts []Option) []Option {
	return append([]Option{
		WithLogger(log.StdLogger{}),
		WithProc(nil, 1),
	}, opts...)
}

// WithProc sets the Proc (processor) to be invoked for each Job.
func WithProc(proc Proc, workerCount int) Option {
	return func(ws *workerSession) error {
		if proc == nil {
			proc = noOpProc
		}
		ws.proc = proc
		ws.workers = workerCount
		return nil
	}
}

// WithLogger sets the logger to be used by worker. If nil, logging is disabled.
func WithLogger(lg log.Logger) Option {
	return func(ws *workerSession) error {
		if lg == nil {
			lg = log.NoOpLogger{}
		}
		ws.Logger = lg
		return nil
	}
}
