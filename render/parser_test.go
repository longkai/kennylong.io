package render

import (
	"bytes"
	j "encoding/json"
	"log"
	"os"
	"strings"
	"testing"
)

func TestNormalParse(t *testing.T) {
	src := `./testdata/normal.md`
	f, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	wants := struct {
		title string
		body  []byte
		link  []byte
		err   error
	}{
		`Title`, []byte(`Body`), []byte(`https://`), nil,
	}

	title, body, json, err := parse(f)

	if err != wants.err {
		log.Fatal(err)
	}
	if !strings.Contains(title, wants.title) {
		t.Errorf("title: strings.Contains(%q, %q) = false\n", title, wants.title)
	}
	if !bytes.Contains(body, wants.body) {
		t.Errorf("body: bytes.Contains(%s, %s) = false\n", body, wants.body)
	}
	if !bytes.Contains(body, wants.link) {
		t.Errorf("body: bytes.Contains(%s, %s) = false\n", body, wants.link)
	}

	m := map[string]interface{}{}
	if err := j.Unmarshal(json, &m); err != nil {
		t.Errorf("json.Unmarshal(%s) fail: %v\n", err)
	}
}

func TestParseNoTitle(t *testing.T) {
	md := `
body1

body2

### EOF
` + "```json" +
		`{"key": "val"}` + "```" +
		`
[1]: https://xiaolongtongxue.com
[anchor]: https://lliant.com`

	if !re.MatchString(md) {
		t.Errorf("MatchString(%q) = false\n", md)
	}

	title, _, json, err := parse(strings.NewReader(md))
	if err != nil {
		log.Fatal(err)
	}

	want := "body1"
	if title != want {
		t.Errorf("parse(%q), no title pattern matches, got %q, want %q\n", md, title, want)
	}

	m := map[string]interface{}{}
	if err := j.Unmarshal(json, &m); err != nil {
		t.Errorf("json.Unmarshal(%s) fail: %v\n", err)
	}
}

func TestParseNoJSON(t *testing.T) {
	src := `./testdata/no_json.md`
	f, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	wantBody, wantLink, wantJSON := `body`, `https://`, ``
	_, body, json, err := parse(f)
	if err != nil {
		log.Fatal(err)
	}

	if !bytes.Contains(body, []byte(wantBody)) {
		t.Errorf("bytes.Contains(%s, %s) = false\n", body, wantBody)
	}

	if !bytes.Contains(body, []byte(wantLink)) {
		t.Errorf("bytes.Contains(%s, %s) = false\n", body, wantLink)
	}

	if string(json) != wantJSON {
		t.Errorf("parse(%q), got json: %s, want %s\n", src, json, wantJSON)
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte(""), ""},
		{[]byte("a"), "a"},
		{[]byte("abc"), "cba"},
		{[]byte("abbc"), "cbba"},
	}

	for _, test := range tests {
		tmp := test.input
		if reverse(test.input); string(test.input) != test.want {
			t.Errorf("reverse(%q) = %q, want %q\n", tmp, string(test.input), test.want)
		}
	}
}
