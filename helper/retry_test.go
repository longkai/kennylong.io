package helper_test

import (
	"fmt"
	"testing"

	"github.com/longkai/xiaolongtongxue.com/helper"
)

func TestRetry(t *testing.T) {
	t.Skip("integration test involves system timer")
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
