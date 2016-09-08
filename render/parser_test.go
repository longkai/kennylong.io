package render

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"
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
			[]byte(`section`),
		}, []byte(`https://`),
	}

	title, body, _yaml, err := parse(f)

	if err != nil {
		t.Errorf("parse(%q) fail: %v\n", src, err)
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
	if err := yaml.Unmarshal(_yaml, &m); err != nil {
		t.Errorf("yaml.Unmarshal(%s) fail: %v\n", _yaml, err)
	}
}

func TestParseNoTitle(t *testing.T) {
	md := `
body1

body2

### EOF
` + "```yaml" +
		`key: val` + "```" +
		`
[1]: https://xiaolongtongxue.com
[anchor]: https://lliant.com`

	if !re.MatchString(md) {
		t.Errorf("MatchString(%q) = false\n", md)
	}

	title, _, _yaml, err := parse(strings.NewReader(md))
	if err != nil {
		log.Fatal(err)
	}

	want := "body1"
	if title != want {
		t.Errorf("parse(%q), no title pattern matches, got %q, want %q\n", md, title, want)
	}

	m := map[string]interface{}{}
	if err := yaml.Unmarshal(_yaml, &m); err != nil {
		t.Errorf("yaml.Unmarshal(%s) fail: %v\n", _yaml, err)
	}
}
