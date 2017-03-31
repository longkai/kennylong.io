package helper_test

import (
	"testing"

	"github.com/longkai/xiaolongtongxue.com/helper"
)

func TestExist(t *testing.T) {
	if !*isIntegration {
		t.Skip("integration test: file system")
	}
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
