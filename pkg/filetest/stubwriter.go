package filetest

import "io"

// StubWriter implements the io.Writer interface.
// It is used during tests to examine the content of the data written,
// and to return fake errors as needed.
type StubWriter struct {
	Buffer []byte
	Count  int
	Err    error
}

var _ io.WriteCloser = (*StubWriter)(nil)

// NewStubWriter returns a new StubWriter.
func NewStubWriter(err error) *StubWriter {
	return &StubWriter{
		Err: err,
	}
}

// Write implements the io.Writer interface.
func (w *StubWriter) Write(p []byte) (n int, err error) {
	if w.Err != nil {
		return 0, w.Err
	}

	w.Buffer = append(w.Buffer, p...)
	w.Count += len(p)

	return len(p), nil
}

// Close implements the io.Closer interface.
func (w *StubWriter) Close() error {
	return nil
}
