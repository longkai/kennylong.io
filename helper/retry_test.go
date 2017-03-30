package helper_test

import (
	"fmt"
	"github.com/longkai/xiaolongtongxue.com/helper"
	"testing"
)

func TestRetry(t *testing.T) {
	var i int
	f := func() (interface{}, error) {
		if i == 1 {
			return nil, nil
		}
		i++
		return nil, fmt.Errorf("times %d fail", i)
	}

	if _, err := helper.Retry(1, f); err == nil {
		t.Errorf("should fail")
	}

	i = 0

	if _, err := helper.Retry(2, f); err != nil {
		t.Errorf("2 time fail")
	}
}
