package repo

import (
	"bufio"
	"io"
	"io/ioutil"
	"regexp"

	"bytes"
	"os"

	"gopkg.in/yaml.v2"
)

// Parser parse the raw document into application representation.
type Parser interface {
	Parse(file string) (Doc, error)
}

// DocParser parse the document.
type DocParser struct {
}

// Parse parse implementation.
func (p *DocParser) Parse(path string) (Doc, error) {
	doc := Doc{}
	f, err := os.Open(path)
	if err != nil {
		return doc, err
	}
	defer f.Close()

	if err = unmarshal(f, &doc); err != nil {
		return doc, err
	}

	return doc, nil
}

var parseYAML = func(in io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

func unmarshal(in io.Reader, d *Doc) error {
	title, _yaml, err := parse(in)
	if err != nil {
		return err
	}

	// Default title if not specify in yaml block.
	d.Title = title
	return parseYAML(bytes.NewReader(_yaml), d)
}

var (
	emptyLineRegex    = regexp.MustCompile(`^\s*$`)
	eofMarkRegex      = regexp.MustCompile(`(?i)\s+EOF\s*$`)
	ymlStartMarkRegex = regexp.MustCompile(`\s*ya?ml\s*$`)
	ymlEndMarkRegex   = regexp.MustCompile(`\s*date:`)
)

// Parse a reader
func parse(in io.Reader) (title string, yml []byte, err error) {
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		line := scanner.Text()
		if !emptyLineRegex.MatchString(line) {
			title = line
			// Treat the first utf8 code point >= '0' as the starting title, if any.
			// The most common mark up symbols are less than '0' in ascii: '#', '*', etc.
			for i, c := range line {
				if c >= '0' {
					title = line[i:]
					break
				}
			}
			break
		}
	}

	for scanner.Scan() {
		if eofMarkRegex.Match(scanner.Bytes()) {
			// Skip empty lines.
			for scanner.Scan() && emptyLineRegex.Match(scanner.Bytes()) {
			}
			// Current line must be a source start mark up.
			if !ymlStartMarkRegex.Match(scanner.Bytes()) {
				continue
			}
			// Gather yaml source.
			buf := make([]byte, 0, 256)
			for scanner.Scan() {
				b := scanner.Bytes()
				buf = append(buf, b...)
				buf = append(buf, '\n')
				// Meta yaml block must end with a `date` field.
				if ymlEndMarkRegex.Match(b) {
					yml = buf
					return
				}
			}
		}
	}
	err = scanner.Err()
	return
}
