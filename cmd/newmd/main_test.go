package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatName(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"hello world", "hello-world"},
		{"HELLO WORLD", "hello-world"},
		{"hello  world", "hello--world"},
		{"hello-world", "hello-world"},
	}
	for _, test := range tests {
		if got := formatName(test.input); got != test.want {
			t.Errorf("formatName(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestNewMD(t *testing.T) {
	out := new(bytes.Buffer)
	title := "balabala"
	err := newMD(title, out)
	if err != nil {
		t.Errorf("newMD(%q, _) fail: %v", title, err)
	}
	got := out.String()
	if !strings.Contains(got, title) {
		t.Errorf("newMD(%q,_) = %s, not contains %q", title, got, title)
	}
}
