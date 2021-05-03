package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type spyExitHandler struct {
	count    int
	lastCode int
}

func (s *spyExitHandler) Exit(code int) {
	s.count++
	s.lastCode = code
}

func TestMain(t *testing.T) {
	spyExit := spyExitHandler{}

	os.Args = []string{os.Args[0], "--version"}
	exitHandler = spyExit.Exit

	main()

	assert.Equal(t, 0, spyExit.count)
	assert.Equal(t, 0, spyExit.lastCode)
}

func TestMainError(t *testing.T) {
	spyExit := spyExitHandler{}

	os.Args = []string{os.Args[0], "invalid-command"}
	exitHandler = spyExit.Exit

	main()

	assert.Equal(t, 1, spyExit.count)
	assert.Equal(t, 1, spyExit.lastCode)
}
