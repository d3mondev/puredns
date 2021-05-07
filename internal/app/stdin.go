package app

import "os"

// HasStdin returns true if there is a valid stdin present.
func HasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	if stat.Mode()&os.ModeNamedPipe == 0 {
		return false
	}

	return true
}
