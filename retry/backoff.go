package retry

import (
	"math"
	"time"
)

func ExpBackoff(base float64, initialTimeout, maxTimeout time.Duration) Backoff {
	var waitTime time.Duration

	return backoffFunc(func(attempt int) time.Duration {
		if attempt == 0 {
			return time.Duration(0)
		} else if waitTime >= maxTimeout {
			return maxTimeout
		}

		waitTime = time.Duration(
			float64(initialTimeout.Nanoseconds()) * math.Pow(base, float64(attempt)),
		)

		if waitTime > maxTimeout {
			return maxTimeout
		}

		return waitTime
	})
}

func ConstBackoff(interval time.Duration) Backoff {
	return backoffFunc(func(_ int) time.Duration {
		return interval
	})
}
