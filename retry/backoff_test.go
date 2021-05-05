package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpBackoff(t *testing.T) {
	backOff := ExpBackoff(2, 100*time.Millisecond, 500*time.Millisecond)

	assert.Equal(t, time.Duration(0), backOff.WaitFor(0))
	assert.Equal(t, 200*time.Millisecond, backOff.WaitFor(1))
	assert.Equal(t, 400*time.Millisecond, backOff.WaitFor(2))
	assert.Equal(t, 500*time.Millisecond, backOff.WaitFor(3))
	assert.Equal(t, 500*time.Millisecond, backOff.WaitFor(4))
	assert.Equal(t, 500*time.Millisecond, backOff.WaitFor(10))
	assert.Equal(t, 500*time.Millisecond, backOff.WaitFor(5000000))
}

func TestConstBackoff(t *testing.T) {
	backOff := ConstBackoff(2 * time.Second)
	assert.Equal(t, 2*time.Second, backOff.WaitFor(0))
	assert.Equal(t, 2*time.Second, backOff.WaitFor(1))
	assert.Equal(t, 2*time.Second, backOff.WaitFor(2))
	assert.Equal(t, 2*time.Second, backOff.WaitFor(3))
	assert.Equal(t, 2*time.Second, backOff.WaitFor(4))
	assert.Equal(t, 2*time.Second, backOff.WaitFor(5))
	assert.Equal(t, 2*time.Second, backOff.WaitFor(10))
}
