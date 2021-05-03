package fileoperation

import (
	"bufio"
	"io"
	"os"
)

// WriteLines writes all lines to a text file, truncating the file if it already exists.
func WriteLines(lines []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return writelines(lines, file, 64*1024)
}

func writelines(lines []string, w io.Writer, bufferSize int) error {
	writer := bufio.NewWriterSize(w, bufferSize)

	for _, line := range lines {
		if _, err := writer.Write([]byte(line + "\n")); err != nil {
			return err
		}
	}

	return writer.Flush()
}
