package app

import (
	"os"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
)

func TestHasStdin_Default(t *testing.T) {
	got := HasStdin()
	assert.Equal(t, false, got)
}

func TestHasStdin_With(t *testing.T) {
	r, _, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdin = r

	got := HasStdin()
	assert.Equal(t, true, got)
}

func TestHasStdin_File(t *testing.T) {
	file := filetest.CreateFile(t, "")
	filetest.OverrideStdin(t, file)

	got := HasStdin()
	assert.Equal(t, false, got)
}

func TestHasStdin_Nil(t *testing.T) {
	filetest.OverrideStdin(t, nil)

	got := HasStdin()
	assert.Equal(t, false, got)
}
