package helper_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/longkai/xiaolongtongxue.com/helper"
)

func TestTry(t *testing.T) {
	t.Skip("integration test: system timer")
	var i int
	f := func() (interface{}, error) {
		if i == 1 {
			return nil, nil
		}
		i++
		return nil, fmt.Errorf("times %d fail", i)
	}

	if _, err := helper.Try(1, f); err == nil {
		t.Errorf("should fail")
	}

	i = 0

	if _, err := helper.Try(2, f); err != nil {
		t.Errorf("2 time fail")
	}
}

func TestTimeout(t *testing.T) {
	t.Skip("integetaion test: system timeout")
	f := func() (interface{}, error) {
		return nil, fmt.Errorf("alwasy fail")
	}
	_, err := helper.Try(100, f) // a very large trying times.
	if err == nil {
		t.Errorf("must fail, got nil")
		return
	}
	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("expect timeout, got %v", err)
	}
}
