package massdns

import (
	"io"
	"strings"
)

// StdoutHandler read complete lines from the massdns output and sends them one by
// one to the callback function.
type StdoutHandler struct {
	callback  Callback
	remainder string
}

var _ io.Writer = (*StdoutHandler)(nil)

// Callback is a callback function that receives lines from the massdns output.
type Callback interface {
	Callback(line string) error
	Close()
}

// NewStdoutHandler returns a new OutputHandler that can be used to receive massdns' stdout.
func NewStdoutHandler(callback Callback) *StdoutHandler {
	return &StdoutHandler{
		callback: callback,
	}
}

// Write detects strings terminated by a \n character in the specified buffer and
// sends them to the callback function.
func (w *StdoutHandler) Write(p []byte) (n int, err error) {
	var builder strings.Builder
	builder.WriteString(w.remainder)

	for n = 0; n < len(p); n++ {
		// If we reach the end of a line, send the line to the callback function
		// even if the line is empty
		if p[n] == byte('\n') {
			line := builder.String()
			builder.Reset()

			if err := w.callback.Callback(line); err != nil {
				return n + 1, err
			}

			continue
		}

		// Build the string from the current data
		builder.WriteByte(p[n])
	}

	// Keep the remainder of the line for the next call to Write
	w.remainder = builder.String()

	return n, nil
}

// Close closes the callback interface.
func (w *StdoutHandler) Close() {
	w.callback.Close()
}
