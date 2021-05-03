package fileoperation

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// AppendWord appends a separator and a word to each lines present in a text file.
func AppendWord(src string, dest string, sep string, word string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	return appendWordIO(srcFile, destFile, sep, word)
}

func appendWordIO(r io.Reader, w io.Writer, sep string, word string) error {
	scanner := bufio.NewScanner(r)
	writer := bufio.NewWriter(w)

	for scanner.Scan() {
		line := scanner.Text() + sep + word
		fmt.Fprintln(writer, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return writer.Flush()
}
