package helper_test

import (
	"github.com/longkai/xiaolongtongxue.com/helper"
	"testing"
)

func TestExist(t *testing.T) {
	t.Skip("skip integration test")
	tests := []struct {
		input string
		want  bool
	}{
		{"/etc/hosts", true},
		{"/etc/balabala", false},
	}

	for _, test := range tests {
		if got := helper.Exists(test.input); got != test.want {
			t.Errorf("helper.Exists(%q) = %t", test.input, got)
		}
	}
}
