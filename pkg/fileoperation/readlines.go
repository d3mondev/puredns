package fileoperation

import (
	"bufio"
	"io"
	"os"
)

// ReadLines read lines from a text file and returns them as a slice.
func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readlines(file)
}

func readlines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
