package git

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
)

const (
	add = `A`
	mod = `M`
	del = `D`
)

var (
	re = regexp.MustCompile(`^([AMD])\s+(.*)`)
)

// Diff `git diff --name-status HEAD~ HEAD`
func Diff(repo string) (adds, mods, dels []string, err error) {
	script := fmt.Sprintf("cd %s && git diff --name-status HEAD~ HEAD", repo)
	log.Printf("run shell script: %q", script)
	cmd := exec.Command("/bin/sh", "-c", script)
	b, err := cmd.Output()
	if err != nil {
		return
	}
	adds, mods, dels = diff(bytes.NewReader(b))
	return
}

func diff(in io.Reader) (adds, mods, dels []string) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		pair := re.FindStringSubmatch(line)
		if len(pair) != 3 {
			continue
		}
		switch pair[1] {
		case add:
			adds = append(adds, pair[2])
		case mod:
			mods = append(mods, pair[2])
		case del:
			dels = append(dels, pair[2])
		default:
			log.Printf("Unknown git diff line: %q", line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("parsing git diff: %v", err)
	}
	return
}
