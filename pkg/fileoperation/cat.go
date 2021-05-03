package fileoperation

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Cat reads files sequentially and sends the output to the specified writer.
func Cat(filenames []string, w io.Writer) error {
	readers := []io.Reader{}

	for _, name := range filenames {
		file, err := os.Open(name)
		if err != nil {
			return err
		}
		defer file.Close()

		readers = append(readers, file)
	}

	return CatIO(readers, w)
}

// CatIO reads sequentially from readers and sends the output to the specified writer.
func CatIO(readers []io.Reader, w io.Writer) error {
	for _, r := range readers {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
				return err
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}

	return nil
}
