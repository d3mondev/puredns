package filetest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStubWriter(t *testing.T) {
	t.Run("internal buffer correctly updated", func(t *testing.T) {
		w := NewStubWriter(nil)

		buf := []byte("test")
		n, err := w.Write(buf)

		assert.Nil(t, err)
		assert.Equal(t, len(buf), n)
		assert.Equal(t, buf, w.Buffer)
	})

	t.Run("generate write error", func(t *testing.T) {
		wantErr := errors.New("error")
		w := NewStubWriter(wantErr)

		buf := []byte("test")
		n, err := w.Write(buf)

		assert.Equal(t, wantErr, err)
		assert.Equal(t, 0, n)
		assert.Equal(t, []byte(nil), w.Buffer)
	})
}
