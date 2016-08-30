package render

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"time"
)

// Meta metadata for the markdown.
type Meta struct {
	Title        string    `json:"title"`
	Tags         []string  `json:"tags"`
	Date         time.Time `json:"date"`
	Weather      string    `json:"weather"`
	Summary      string    `json:"summary"`
	Location     string    `json:"location"`
	Background   string    `json:"background"`
	RenderOption int       `json:"render_option"`
}

var parseJSON = func(in io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func parseMd(in io.Reader) (*Meta, []byte, error) {
	title, body, _json, err := parse(in)
	if err != nil {
		return nil, nil, err
	}

	m := new(Meta)
	m.Title = title
	if err = parseJSON(bytes.NewReader(_json), m); err != nil {
		log.Print(err)
	}
	return m, body, nil
}
