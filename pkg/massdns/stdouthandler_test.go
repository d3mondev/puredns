package massdns

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubCallback struct {
	data    []string
	returns error
	closed  bool
}

func (c *stubCallback) Callback(line string) error {
	if c.returns != nil {
		return c.returns
	}

	c.data = append(c.data, line)

	return nil
}

func (c *stubCallback) Close() {
	c.closed = true
}

func TestOutputHandlerNew(t *testing.T) {
	var cb stubCallback
	handler := NewStdoutHandler(&cb)
	assert.NotNil(t, handler)
}

func TestOutputHandlerWrite(t *testing.T) {
	callbackError := errors.New("error")

	tests := []struct {
		name        string
		haveBuffers [][]byte
		haveError   error
		want        []string
		wantErr     error
	}{
		{
			name: "empty write",
		},
		{
			name: "no newline",
			haveBuffers: [][]byte{
				[]byte("line"),
			},
		},
		{
			name: "with newline",
			haveBuffers: [][]byte{
				[]byte("line\n"),
			},
			want: []string{"line"},
		},
		{
			name: "multiple lines",
			haveBuffers: [][]byte{
				[]byte("line1\nline2\nline3\n"),
			},
			want: []string{"line1", "line2", "line3"},
		},
		{
			name: "partial line",
			haveBuffers: [][]byte{
				[]byte("line1\nli"),
				[]byte("ne2\nline3\n"),
			},
			want: []string{"line1", "line2", "line3"},
		},
		{
			name: "callback error",
			haveBuffers: [][]byte{
				[]byte("line\n"),
			},
			haveError: callbackError,
			wantErr:   callbackError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cb := &stubCallback{returns: test.haveError}
			handler := NewStdoutHandler(cb)

			for _, buf := range test.haveBuffers {
				n, err := handler.Write(buf)
				assert.Equal(t, len(buf), n)
				assert.ErrorIs(t, err, test.wantErr)
			}

			assert.Equal(t, test.want, cb.data)
		})
	}
}

func TestOutputHandlerClose(t *testing.T) {
	var cb stubCallback
	handler := NewStdoutHandler(&cb)
	handler.Close()
	assert.True(t, cb.closed)
}
