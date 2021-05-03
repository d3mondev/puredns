package procreader

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	r := New(func(int) ([]byte, error) { return []byte{}, nil })
	assert.NotNil(t, r)
}

func TestRead_Buffer(t *testing.T) {
	var cb = func(int) ([]byte, error) {
		return []byte("this is a test"), io.EOF
	}

	tests := []struct {
		name       string
		haveBuffer []byte
		wantRead   []byte
		wantErr    error
	}{
		{name: "no buffer", haveBuffer: nil, wantRead: nil, wantErr: io.ErrShortBuffer},
		{name: "small buffer", haveBuffer: make([]byte, 1), wantRead: []byte("this is a test"), wantErr: io.EOF},
		{name: "big buffer", haveBuffer: make([]byte, 256), wantRead: []byte("this is a test"), wantErr: io.EOF},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := New(cb)

			var err error
			var readBuf []byte

			for err == nil {
				var n int
				n, err = r.Read(test.haveBuffer)
				readBuf = append(readBuf, test.haveBuffer[:n]...)
			}

			assert.ErrorIs(t, err, test.wantErr)
			assert.Equal(t, test.wantRead, readBuf)
		})
	}
}

func TestRead_MultipleCallbacks(t *testing.T) {
	data := [][]byte{
		[]byte("first callback"),
		[]byte("second callback"),
	}

	var cb = func(int) ([]byte, error) {
		val := data[0]
		data = data[1:]

		var err error
		if len(data) == 0 {
			err = io.EOF
		}

		return val, err
	}

	r := New(cb)

	var err error
	var readBuf []byte

	buffer := make([]byte, 1)
	for err == nil {
		var n int
		n, err = r.Read(buffer)
		readBuf = append(readBuf, buffer[:n]...)
	}

	assert.ErrorIs(t, err, io.EOF)
	assert.Equal(t, "first callbacksecond callback", string(readBuf))
}
