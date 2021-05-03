package filetest

import "io"

// StubReader implements the io.ReadCloser interface.
// It is used during tests to return the content of its internal buffer,
// or to return an error as needed.
type StubReader struct {
	Buffer []byte
	Index  int
	Err    error
}

// NewStubReader returns a new StubReader object that reads from the buffer specified.
// If err is not nil, it will return an error on the first read. If the end of the buffer is
// reached, it will return an io.EOF error.
func NewStubReader(buffer []byte, err error) *StubReader {
	return &StubReader{
		Buffer: buffer,
		Err:    err,
	}
}

// Read implements the io.Reader interface.
func (r *StubReader) Read(p []byte) (int, error) {
	if r.Err != nil {
		return 0, r.Err
	}

	count := copy(p, r.Buffer[r.Index:])
	r.Index += count

	if r.Err == nil && r.Index >= len(r.Buffer) {
		r.Err = io.EOF
	}

	return count, r.Err
}

// Close implements the io.Closer interface.
func (r *StubReader) Close() error {
	return nil
}
