package git

import (
	"errors"
	"strings"
	"testing"
)

func TestExecer(t *testing.T) {
	saved := execer
	defer func() { execer = saved }()

	var wantError bool
	execer = func(script string) ([]byte, error) {
		var err error
		if wantError {
			err = errors.New(`balh`)
		}
		return nil, err
	}

	f := func(t *testing.T) {
		_, _, _, err := Diff(`don't care`, `balh`)
		if got := err != nil; got != wantError {
			t.Errorf("Diff(...): %v, wantError %t", err, wantError)
		}
	}

	wantError = true
	t.Run(`WantError`, f)

	wantError = false
	t.Run(`NoError`, f)
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
