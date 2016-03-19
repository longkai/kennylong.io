package render

import (
	"io/ioutil"
	"testing"
)

func TestMarkdownRederRead(t *testing.T) {
	m := New("hello, world")
	if _, err := ioutil.ReadAll(m); err != nil {
		t.Errorf("reading fail, %v\n", err)
	}
}

func TestMarkdownRender(t *testing.T) {
	m := New("Hello world github/longkai#1 **cool**, and #1!")
	if _, err := m.Render(); err != nil {
		t.Errorf("render fail, %v\n", err)
	}
}
