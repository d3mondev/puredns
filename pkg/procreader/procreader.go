package procreader

import (
	"io"
)

// ProcReader is a procedural reader that generates its data from a callback function.
type ProcReader struct {
	callback  Callback
	remainder []byte
	err       error
}

var _ io.Reader = (*ProcReader)(nil)

// Callback is a callback function that generates data in the form of a slice of bytes.
// Size is a hint as to how much data is requested. If the callback returns more data, the
// excess will be buffered by the reader for subsequent Read calls. If no data is left,
// the Callback function must returns an io.EOF error.
type Callback func(size int) ([]byte, error)

// New creates a new ProcReader.
func New(callback Callback) *ProcReader {
	return &ProcReader{
		callback: callback,
	}
}

// Read requests data from the callback until either the buffer is full, or an error like EOF occurs.
func (r *ProcReader) Read(p []byte) (int, error) {
	// Buffer cannot be nil, otherwise an error like EOF will never be returned
	if p == nil {
		return 0, io.ErrShortBuffer
	}

	var written int
	total := len(p)

	for {
		var data []byte

		// Get the data to write to the buffer
		if r.remainder != nil {
			data = r.remainder
			r.remainder = nil
		} else {
			data, r.err = r.callback(total - written)
		}

		// Write to buffer
		n := copy(p[written:], data)
		written += n

		// Could not write entire data, save remainder and exit
		if n < len(data) {
			r.remainder = data[n:]
			return written, nil
		}

		// Could save entire data, but buffer is full
		if written == len(p) {
			return written, r.err
		}

		// Error or EOF while creating data, return
		if r.err != nil {
			return written, r.err
		}
	}
}
