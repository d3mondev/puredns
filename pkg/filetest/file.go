package filetest

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

// CreateFile creates a temporary file with the content specified used during the test.
// The file is closed and deleted after the test is done running.
func CreateFile(t *testing.T, content string) *os.File {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := file.WriteString(content); err != nil {
		t.Fatal(err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})

	return file
}

// CreateDir creates a temporary directory.
// The directory is deleted after the test is done running.
func CreateDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.Remove(dir)
	})

	return dir
}

// ReadFile reads a text file and returns each line in a slice.
// If the file name is empty, returns an empty slice.
func ReadFile(t *testing.T, name string) []string {
	lines := []string{}

	if name == "" {
		return lines
	}

	file, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	return lines
}

// ClearFile truncates the content of a file.
func ClearFile(t *testing.T, file *os.File) {
	if err := file.Truncate(0); err != nil {
		t.Fatal(err)
	}
}
