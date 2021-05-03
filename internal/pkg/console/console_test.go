package console

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func redirectOutput() *bytes.Buffer {
	buffer := bytes.NewBuffer([]byte{})
	Output = buffer

	return buffer
}

type spyExitHandler struct {
	seen int
}

func (s *spyExitHandler) Exit(int) {
	s.seen++
}

func TestMessage(t *testing.T) {
	buffer := redirectOutput()
	Message("foo %s", "bar")
	assert.True(t, strings.Contains(buffer.String(), "foo bar"))
}

func TestSuccess(t *testing.T) {
	buffer := redirectOutput()
	Success("foo %s", "bar")
	assert.True(t, strings.Contains(buffer.String(), "foo bar"))
}

func TestWarning(t *testing.T) {
	buffer := redirectOutput()
	Warning("foo %s", "bar")
	assert.True(t, strings.Contains(buffer.String(), "foo bar"))
}

func TestError(t *testing.T) {
	buffer := redirectOutput()
	Error("foo %s", "bar")
	assert.True(t, strings.Contains(buffer.String(), "foo bar"))
}

func TestPrintf(t *testing.T) {
	buffer := redirectOutput()
	Printf("foo %s", "bar")
	assert.True(t, strings.Contains(buffer.String(), "foo bar"))
}

func TestFatal(t *testing.T) {
	buffer := redirectOutput()
	spyExitHandler := spyExitHandler{}
	ExitHandler = spyExitHandler.Exit

	Fatal("foo %s", "bar")
	assert.True(t, strings.Contains(buffer.String(), "foo bar"))
	assert.Equal(t, 1, spyExitHandler.seen)
}
