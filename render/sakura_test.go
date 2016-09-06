package render

import (
	"errors"
	"io"
	"testing"
)

func TestGet(t *testing.T) {
	s := NewSakura().(*Sakura)
	var (
		withError   bool
		calledTimes int
		wantTimes   int
	)
	s.Render = func(in io.Reader) (interface{}, error) {
		calledTimes++
		if withError {
			return nil, errors.New("balabala")
		}
		return []byte("balabala"), nil
	}

	// arrange testdata, should we need mock?
	key, ts := "balabala", timestamp(1234567)
	s.index[key] = ts
	s.list.Set(ts, new(Meta))
	// note non-cache

	t.Run("Retry", func(t *testing.T) {
		withError, wantTimes = true, 123
		// don't care what the result
		for i := 0; i < wantTimes; i++ {
			s.Get(key)
		}
		// at this time calledTimes should be 3
		if calledTimes != wantTimes {
			t.Errorf("call Sakura.Get(_) with error %d times, got %d, want %d", calledTimes, calledTimes, wantTimes)
		}
	})

	calledTimes = 0
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
