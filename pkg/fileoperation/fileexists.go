package fileoperation

import (
	"os"
)

// FileExists returns true if a file exists on disk.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
