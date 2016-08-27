package helper

import (
	"bytes"
	"testing"
)

func TestJSON2Yaml(t *testing.T) {
	test := struct {
		input string
		want  string
	}{
		`{"foo": "bar"}`,
		`{foo:bar}`,
	}

	if got, err := JSON2Yaml([]byte(test.input)); err != nil || bytes.HasSuffix(got, []byte(test.want)) {
		t.Errorf("JSON2Yaml(%q) = (%v, %q), want (nil, %q)\n", test.input, err, string(got), test.want)
	}
}
