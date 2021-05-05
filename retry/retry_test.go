package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(suite *testing.T) {
	suite.Parallel()

	suite.Run("Success", func(t *testing.T) {
		backOff := ExpBackoff(2, 100*time.Millisecond, 500*time.Millisecond)
		ctx := context.Background()
		handlerCalled := false
		handler := func() error {
			handlerCalled = true
			return nil
		}
		_ = Retry(ctx, 3, backOff, handler)
		assert.Equal(suite, true, handlerCalled)
	})

	suite.Run("RetriesFailure", func(t *testing.T) {
		backOff := ExpBackoff(2, 100*time.Millisecond, 500*time.Millisecond)
		ctx := context.Background()
		handlerCalledCount := 0
		handler := func() error {
			handlerCalledCount++
			return errors.New("some random error")
		}
		err := Retry(ctx, 3, backOff, handler)
		assert.Equal(t, 4, handlerCalledCount)
		assert.Contains(t, "some random error", err.Error())
	})

	suite.Run("ContextCancelled", func(t *testing.T) {
		backOff := ExpBackoff(2, 100*time.Millisecond, 5000*time.Millisecond)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		handler := func() error {
			return errors.New("some random error")
		}
		err := Retry(ctx, 13, backOff, handler)
		assert.Equal(t, context.Canceled, err)
	})

	suite.Run("SuccessAfterRetry", func(t *testing.T) {
		backOff := ExpBackoff(2, 100*time.Millisecond, 500*time.Millisecond)
		ctx := context.Background()
		handlerCalledCount := 0
		handler := func() error {
			handlerCalledCount++
			if handlerCalledCount == 2 {
				return nil
			}
			return errors.New("some random error")
		}
		_ = Retry(ctx, 4, backOff, handler)
		assert.Equal(t, 2, handlerCalledCount)
	})
}
