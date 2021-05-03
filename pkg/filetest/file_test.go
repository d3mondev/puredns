package filetest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFile_Empty(t *testing.T) {
	file := CreateFile(t, "")
	lines := ReadFile(t, file.Name())

	assert.NotNil(t, file)
	assert.Equal(t, []string{}, lines)
}

func TestCreateFile_Content(t *testing.T) {
	file := CreateFile(t, "foo\nbar")
	lines := ReadFile(t, file.Name())

	assert.NotNil(t, file)
	assert.Equal(t, []string{"foo", "bar"}, lines)
}

func TestCreateDir(t *testing.T) {
	dir := CreateDir(t)
	assert.NotEqual(t, "", dir)
}

func TestReadFile_OK(t *testing.T) {
	file := CreateFile(t, "line1\nline2\nline3")
	lines := ReadFile(t, file.Name())
	assert.Equal(t, []string{"line1", "line2", "line3"}, lines)
}

func TestReadFile_Empty(t *testing.T) {
	lines := ReadFile(t, "")
	assert.Equal(t, []string{}, lines)
}

func TestClearFile(t *testing.T) {
	file := CreateFile(t, "foo\nbar")

	ClearFile(t, file)

	lines := ReadFile(t, file.Name())
	assert.Equal(t, []string{}, lines)
}
