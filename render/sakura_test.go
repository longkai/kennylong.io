package render

import (
	"errors"
	"io"
	"testing"
)

func TestGet(t *testing.T) {
	s := NewSakura(``).(*Sakura)
	var (
		withError   bool
		calledTimes int
		wantTimes   int
	)
	s.render = func(id string, in io.Reader) (interface{}, error) {
		calledTimes++
		if withError {
			return nil, errors.New("balabala")
		}
		return []byte("balabala"), nil
	}

	// arrange testdata, should we need mock?
	key, ts := "balabala", timestamp(1234567)
	s.index[key] = ts
	m := new(Meta)
	m.Body = []byte("balabala")
	s.list.Set(ts, m)
	// note non-cache

	t.Run("Retry", func(t *testing.T) {
		withError, wantTimes = true, 3 // max retry
		// don't care what the result
		for i := 0; i < wantTimes; i++ {
			s.Get(key)
		}
		if calledTimes != wantTimes {
			t.Errorf("call Sakura.Get(_) with error %d times, got %d, want %d", wantTimes, calledTimes, wantTimes)
		}
	})

	// reset cache
	calledTimes = 0
	delete(s.cache, key)
	t.Run("Cache", func(t *testing.T) {
		// without error and cache enabled
		withError, wantTimes = false, 1
		n := 100
		for i := 0; i < n; i++ {
			s.Get(key)
		}
		if calledTimes != wantTimes {
			t.Errorf("call Sakura.Get(_) without error %d times, got %d, want %d", n, calledTimes, wantTimes)
		}
	})
}
