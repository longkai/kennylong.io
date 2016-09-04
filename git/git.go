package git

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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

var execer = func(script string) ([]byte, error) {
	os.Stdout.WriteString(script)
	b, err := exec.Command(`/bin/sh`, `-c`, script).CombinedOutput()
	os.Stdout.Write(b) // if b is nil, nothing outputs
	return b, err
}

// Pull `git pull`
func Pull(repo string) error {
	script := fmt.Sprintf("git -C %s pull", repo)
	_, err := execer(script)
	return err
}

// Rev `git rev-parse --short HEAD`
func Rev(repo string) (string, error) {
	script := fmt.Sprintf("git -C %s rev-parse --short HEAD", repo)
	b, err := execer(script)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(b)), nil
}

// Diff `git diff --name-status [hash] to HEAD`
func Diff(repo, hash string) (adds, mods, dels []string, err error) {
	script := fmt.Sprintf("git -C %s diff --name-status %s HEAD", repo, hash)
	b, err := execer(script)
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
		triple := re.FindStringSubmatch(line)
		if len(triple) != 3 {
			continue
		}
		switch triple[1] {
		case add:
			adds = append(adds, triple[2])
		case mod:
			mods = append(mods, triple[2])
		case del:
			dels = append(dels, triple[2])
		default:
			log.Printf("unknown git diff line: %q", line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("parsing git diff: %v", err)
	}
	return
}
