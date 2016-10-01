package controller

import "testing"

func TestRevAsset(t *testing.T) {
	_cdn, _v := cdn, v
	defer func() { cdn, v = _cdn, _v }()
	cdn = "//awesome.com"
	v = "v1"
	tests := []struct {
		input, want string
	}{
		{"a.js", "a-v1.js"},
		{"a.min.js", "a.min-v1.js"},
		{"a.jpg", "a.jpg"},
	}
	for _, test := range tests {
		if got := revAsset(test.input); got != test.want {
			t.Errorf("revAsset(%q) = %s, want %s", test.input, got, test.want)
		}
	}
}
