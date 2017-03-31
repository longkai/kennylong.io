package helper

import (
	"time"
)

// TODO: let caller specify timeout properties when necessary.
const threshold = time.Minute

// RetryTimeoutError return when the retry delay exceeds the max threshold.
type RetryTimeoutError time.Duration

func (t RetryTimeoutError) Error() string { return "timeout: exceeds " + string(t) }

// Try a simple exponential retry strategy whose first retry starting after 1s.
// The max delay timeout is 1min, if exceeds, a timeout error will be returned.
func Try(n int, f func() (interface{}, error)) (interface{}, error) {
	if n < 1 {
		panic("trying times should be >= 1")
	}
	var (
		i     int
		val   interface{}
		err   error
		delay = time.Second
	)
	for {
		val, err = f()
		i++
		if err == nil {
			break
		}
		if i >= n {
			break
		}
		// time.After() seems expensive than Sleep(), which introduces a channel.
		if delay > threshold {
			return nil, RetryTimeoutError(threshold)
		}
		time.Sleep(delay)
		delay <<= 1
	}
	return val, err
}
