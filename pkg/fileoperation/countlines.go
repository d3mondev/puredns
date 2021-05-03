package fileoperation

import (
	"bytes"
	"io"
	"os"
)

// CountLines counts the number of lines in a file.
func CountLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return count(file)
}

func count(reader io.Reader) (int, error) {
	buf := make([]byte, 64*1024)
	count := 0
	sep := []byte{'\n'}

	for {
		c, err := reader.Read(buf)
		count += bytes.Count(buf[:c], sep)

		if err == io.EOF {
			return count, nil
		}

		if err != nil {
			return 0, err
		}
	}
}
