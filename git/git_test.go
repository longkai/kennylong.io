package git

import (
	"strings"
	"testing"
)

func TestDiffCmd(t *testing.T) {
	// exe the command in the current repo is okay.
	if _, _, _, err := Diff(`.`); err != nil {
		t.Errorf("Diff(.): %v", err)
	}
}

func TestDiff(t *testing.T) {
	input := `
A       cp2.txt
D       cp1.txt
M       dir/inner.txt
D       renamef1.txt
A       re\ namef1.txt
D       re namef1.txt
?       what.txt
`

	wants := struct {
		a []string
		m []string
		d []string
	}{
		[]string{"cp2.txt", `re\ namef1.txt`},
		[]string{"dir/inner.txt"},
		[]string{"cp1.txt", "renamef1.txt", "re namef1.txt"},
	}

	if a, m, d := diff(strings.NewReader(input)); !sliceEq(a, wants.a) || !sliceEq(m, wants.m) || !sliceEq(d, wants.d) {
		t.Errorf("diff(%s) = (%s, %s, %s), want (%s, %s, %s)", input, a, m, d, wants.a, wants.m, wants.d)
	}
}

func sliceEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
