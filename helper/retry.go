package helper

import (
	"time"
)

// Retry a simple exponential retry strategy, starting from 1s.
func Retry(n int, f func() (interface{}, error)) (interface{}, error) {
	var i = 0
	var val interface{}
	var err error
	var delay = time.Second
	for {
		val, err = f()
		i++
		if err == nil {
			break
		}
		if i >= n {
			break
		}
		_ = <-time.After(delay)
		delay *= 2
	}
	return val, err
}
