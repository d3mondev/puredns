package fileoperation

import (
	"bufio"
	"io"
	"os"
)

// Copy copies a file from the source filename to the destination filename.
func Copy(src string, dest string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()

	return copyByBuffer(source, destination, 64*1024)
}

func copyByBuffer(r io.Reader, w io.Writer, bufferSize int) error {
	buf := make([]byte, bufferSize)
	writer := bufio.NewWriterSize(w, bufferSize)

	for {
		n, err := r.Read(buf)

		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := writer.Write(buf[:n]); err != nil {
			return err
		}
	}

	return writer.Flush()
}
