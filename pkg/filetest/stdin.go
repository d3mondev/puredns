package filetest

import (
	"os"
	"testing"
)

// OverrideStdin overrides os.Stdin with the specified file and restores it after the test.
func OverrideStdin(t *testing.T, f *os.File) {
	old := os.Stdin
	os.Stdin = f

	t.Cleanup(func() {
		os.Stdin = old
	})
}
