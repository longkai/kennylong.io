package git

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	Renamed   = 'R'
	Deleted   = 'D'
	Modified  = 'M'
	Untracked = '?'
	// not interested other status

	spiltor = " -> "
)

func ParseRename(s string) (string, string) {
	files := strings.Split(s, spiltor)
	if len(files) == 1 {
		return files[0], ""
	}
	return files[0], files[1]
}

func Status(path string) (map[byte][]string, error) {
	script := fmt.Sprintf("cd %s; git status --porcelain", path)
	log.Printf("exec shell script, %q\n", script)
	cmd := exec.Command("/bin/bash", "-c", script)
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", b)
	return status(b), nil
}

func status(output []byte) map[byte][]string {
	scanner := bufio.NewScanner(bytes.NewReader(output))
	m := make(map[byte][]string)
	for scanner.Scan() {
		code, path := separate(scanner.Bytes())
		switch code {
		case Renamed:
			fallthrough
		case Deleted:
			fallthrough
		case Modified:
			fallthrough
		case Untracked:
			m[code] = append(m[code], path)
		default:
			log.Printf("skip git status code %q, path is %s\n", code, path)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	return m
}

// separate status code and path
func separate(b []byte) (byte, string) {
	if len(b) < 4 {
		return 0, ""
	}
	// input may like below

	// A  a.txt
	//  M a/b/c.txt
	// ?? path/to/file

	if b[0] == ' ' {
		return b[1], string(b[3:])
	} else {
		return b[0], string(b[3:])
	}
}
