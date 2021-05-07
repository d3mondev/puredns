package fileoperation

import (
	"io/fs"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
)

func TestReadLines(t *testing.T) {
	file := filetest.CreateFile(t, "foo\nbar")

	lines, err := ReadLines(file.Name())
	assert.Nil(t, err)
	assert.Equal(t, []string{"foo", "bar"}, lines)
}

func TestReadLines_FileNotFound(t *testing.T) {
	_, err := ReadLines("thisfiledoesnotexist.txt")
	assert.ErrorIs(t, err, fs.ErrNotExist)
}
