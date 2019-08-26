package repo

import (
	"bufio"
	"errors"
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

	err = unmarshal(f, &doc)

	return doc, err
}

var parseDoc = func(in io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

func unmarshal(in io.Reader, d *Doc) error {
	title, yml, body, err := parse(in)
	if err != nil {
		return err
	}
	if yml == nil {
		return errors.New("yaml block is empty")
	}

	// Default title if not specify in yaml block.
	d.Title = title
	d.rawBody = body

	return parseDoc(bytes.NewReader(yml), d)
}

var (
	emptyLineRegex    = regexp.MustCompile(`^\s*$`)
	eofMarkRegex      = regexp.MustCompile(`(?i)\s+EOF\s*$`)
	ymlStartMarkRegex = regexp.MustCompile(`ya?ml\s*$`)
	ymlEndMarkRegex   = regexp.MustCompile(`^\s*date:`)
)

// Parse a title and yaml meta block from a reader.
// The logic is:
// A headline followed by a yaml code block, and ends with a `data` field.
// Empty line in between will be stripped.
func parse(in io.Reader) (title string, yml []byte, body string, err error) {
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		line := scanner.Text()
		if !emptyLineRegex.MatchString(line) {
			title = line
			// Treat the first utf8 rune great and equal than '0'
			// as the starting title, if any.
			// Luckily, The most common mark up symbols are less than '0'
			// in ascii: '#', '*', etc.
			for i, c := range line {
				if c >= '0' {
					title = line[i:]
					break
				}
			}
			break
		}
	}

	buf := make([]byte, 0, 256)
	var tmp bytes.Buffer
	for scanner.Scan() {
		if eofMarkRegex.Match(scanner.Bytes()) {
		try:
			// Skip empty lines.
			for scanner.Scan() && emptyLineRegex.Match(scanner.Bytes()) {
			}
			// Current line must be a source start mark up.
			if !ymlStartMarkRegex.Match(scanner.Bytes()) {
				continue
			}
			// Gather yaml source.
			for scanner.Scan() {
				b := scanner.Bytes()
				// If we find another `EOF`, drop current and try parsing again.
				if eofMarkRegex.Match(b) {
					buf = buf[:0] // Reset buffer.
					goto try
				}
				buf = append(buf, b...)
				buf = append(buf, '\n')
				// Meta yaml block must end with a `date` field.
				if ymlEndMarkRegex.Match(b) {
					yml = buf
					body = tmp.String()
					return
				}
			}
		} else {
			tmp.WriteString(scanner.Text())
			tmp.WriteByte('\n')
		}
	}
	err = scanner.Err()
	return
}
