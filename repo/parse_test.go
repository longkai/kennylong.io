package repo_test

import (
	"testing"

	"github.com/longkai/xiaolongtongxue.com/repo"
)

var parser = new(repo.DocParser)

func TestNormalParse(t *testing.T) {
	d, err := parser.Parse("./testdata/normal.org")
	if err != nil {
		t.Errorf("parse fail: %v", err)
	}
	if d.Title != "Title" {
		t.Errorf("title = %q, want %q", d.Title, "Title")
	}
	if d.Weather != "ok" {
		t.Errorf("weather = %q, want %q", d.Weather, "ok")
	}
}

func TestComplexParse(t *testing.T) {
	d, err := parser.Parse("./testdata/complex.org")
	if err != nil {
		t.Errorf("parse fail: %v", err)
	}
	if d.Title != "Title" {
		t.Errorf("title = %q, want %q", d.Title, "Title")
	}
	if d.Weather != "cold" {
		t.Errorf("weather = %q, want %q", d.Weather, "cold")
	}
}

func TestBadFormat(t *testing.T) {
	_, err := parser.Parse("./testdata/bad.md")
	if err == nil {
		t.Errorf("should fail")
	}
}
