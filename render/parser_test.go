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
		bodys [][]byte
		link  []byte
	}{
		`标题`, [][]byte{
			[]byte(`Body`),
			[]byte(`header2`),
		}, []byte(`https://`),
	}

	title, body, json, err := parse(f)

	if err != nil {
		t.Errorf("parse(%q) fail: %v\n", err)
	}
	if !strings.Contains(title, wants.title) {
		t.Errorf("title: strings.Contains(%q, %q) = false\n", title, wants.title)
	}
	for _, b := range wants.bodys {
		if !bytes.Contains(body, b) {
			t.Errorf("body: bytes.Contains(%s, %s) = false\n", body, b)
		}
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
