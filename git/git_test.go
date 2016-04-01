package git

import (
	"sort"
	"testing"
)

func TestExecStatus(t *testing.T) {
	// just testing cmd exec successfully, don' t care the result
	path := "."
	_, err := Status(path)
	if err != nil {
		t.Errorf("Status(%q) = (_, %v), want (_, nil)\n", path, err)
	}
}

func TestStatus(t *testing.T) {
	test := struct {
		input []byte
		want  map[byte][]string
	}{
		[]byte(`M  a.txt
R  a.txt -> b.txt
 M a.txt
D  b.txt
?? a/b/c.txt
A  a.txt
`), map[byte][]string{
			Renamed:   []string{"a.txt -> b.txt"},
			Deleted:   []string{"b.txt"},
			Modified:  []string{"a.txt", "a.txt"},
			Untracked: []string{"a/b/c.txt"},
		},
	}

	if got := status(test.input); !mapEquals(got, test.want) {
		t.Errorf("status(\"%s\") = %v, want %v\n", test.input, got, test.want)
	}
}

func mapEquals(m1, m2 map[byte][]string) bool {
	keys1, keys2 := make([]int, 0), make([]int, 0)
	// no generic, hehe
	for k := range m1 {
		keys1 = append(keys1, int(k))
	}
	for k := range m2 {
		keys2 = append(keys2, int(k))
	}
	// test keys are equal
	sort.Ints(keys1)
	sort.Ints(keys2)
	if !intsEquals(keys1, keys2) {
		return false
	}
	// test values are equal
	for _, k := range keys1 {
		v1, v2 := m1[byte(k)], m2[byte(k)]
		if !stringsEquals(v1, v2) {
			return false
		}
	}
	return true
}

func intsEquals(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func stringsEquals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSeparate(t *testing.T) {
	tests := []struct {
		input    []byte
		wantCode byte
		wantPath string
	}{
		{[]byte(`A  a.txt`), 'A', "a.txt"},
		{[]byte(` M a.txt`), 'M', "a.txt"},
		{[]byte(`?? a/b/c.txt`), '?', "a/b/c.txt"},
		{[]byte(""), 0, ""},
		{[]byte("A  "), 0, ""},
		{[]byte("A   "), 'A', " "},
	}

	for _, test := range tests {
		if code, path := separate(test.input); code != test.wantCode || path != test.wantPath {
			t.Errorf("separate(%q) = (%q, %q), want (%q, %q)\n", test.input, code, path, test.wantCode, test.wantPath)
		}
	}
}

func TestFromRename(t *testing.T) {
	tests := []struct {
		input    string
		wantFrom string
		wantTo   string
	}{
		{"a.go -> b.go", "a.go", "b.go"},
		{"a.go ->", "a.go ->", ""},
		{"  -> b.go", " ", "b.go"},
		{"a.go b.go", "a.go b.go", ""},
	}

	for _, test := range tests {
		if gotFrom, gotTo := ParseRename(test.input); gotTo != test.wantTo || gotFrom != test.wantFrom {
			t.Errorf("ParseRename(%q) = (%q, %q), want (%q, %q)\n", test.input, gotFrom, gotTo, test.wantFrom, test.wantTo)
		}
	}
}
