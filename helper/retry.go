package helper

import (
	"errors"
	"time"
)

// Try a simple exponential retry strategy whose first retry starting after 1s.
func Try(n int, f func() (interface{}, error)) (interface{}, error) {
	if n < 1 {
		return nil, errors.New("trying times should be >= 1")
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
		time.Sleep(delay)
		delay <<= 1
	}
	return val, err
}
