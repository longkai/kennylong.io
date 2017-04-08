package helper

import (
	"os"
)

// Exists test a path existence on the OS.
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
