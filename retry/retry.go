package retry

import (
	"context"
	"time"
)

const Forever = -1

func Retry(ctx context.Context, maxAttempts int, backoff Backoff, fn func() error) error {
	var waitTime time.Duration
	var attemptsDone int

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-time.After(waitTime):
			if err := fn(); err == nil || (attemptsDone >= maxAttempts && maxAttempts != Forever) {
				return err
			}
			attemptsDone++
			waitTime = backoff.WaitFor(attemptsDone)
		}
	}
}

// Backoff implementations define the backoff strategy for retry.
type Backoff interface {
	WaitFor(retriesDone int) time.Duration
}

// backoffFunc is an adaptor to allow using ordinary Go functions as Backoff
// strategy.
type backoffFunc func(retriesDone int) time.Duration

func (bf backoffFunc) WaitFor(retriesDone int) time.Duration {
	return bf(retriesDone)
}
