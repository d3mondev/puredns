package filetest

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStubReaderRead(t *testing.T) {
	t.Run("read full buffer", func(t *testing.T) {
		buffer := []byte("foo")
		r := NewStubReader(buffer, nil)

		output := make([]byte, len(buffer))
		count, err := r.Read(output)

		assert.Equal(t, err, io.EOF)
		assert.Equal(t, len(output), count)
		assert.Equal(t, buffer, output)
	})

	t.Run("read part of buffer", func(t *testing.T) {
		buffer := []byte("foo")
		r := NewStubReader(buffer, nil)

		output := make([]byte, 1)
		count, err := r.Read(output)

		assert.Nil(t, err)
		assert.Equal(t, len(output), count)
		assert.Equal(t, []byte("f"), output)
	})

	t.Run("generate read error", func(t *testing.T) {
		buffer := []byte("foo")
		readErr := errors.New("read error")
		r := NewStubReader(buffer, readErr)

		output := make([]byte, 1)
		count, err := r.Read(output)

		assert.Equal(t, readErr, err)
		assert.Equal(t, 0, count)
		assert.Equal(t, []byte{0x0}, output)
	})
}

func TestStubReaderClose(t *testing.T) {
	r := NewStubReader(nil, nil)
	err := r.Close()
	assert.Nil(t, err)
}
