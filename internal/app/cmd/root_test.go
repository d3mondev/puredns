package cmd

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/internal/pkg/console"
	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
)

func overrideStdin(t *testing.T, f *os.File) {
	old := os.Stdin
	os.Stdin = f

	t.Cleanup(func() {
		os.Stdin = old
	})
}

func TestNewCmdRoot(t *testing.T) {
	cmd := newCmdRoot()
	assert.NotNil(t, cmd)
}

func TestPreRun(t *testing.T) {
	t.Run("quiet mode", func(t *testing.T) {
		context = ctx.NewCtx()
		context.Options.Quiet = true

		cmd := newCmdResolve()
		preRun(cmd, []string{})

		assert.Equal(t, console.Output, io.Discard)
	})
}

func TestMust(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		must(nil)
	})

	t.Run("panics on error", func(t *testing.T) {
		assert.Panics(t, func() { must(errors.New("error")) })
	})
}

func TestHasStdin(t *testing.T) {
	t.Run("default stdin", func(t *testing.T) {
		want := false
		got := hasStdin()
		assert.Equal(t, want, got)
	})

	t.Run("with stdin", func(t *testing.T) {
		r, _, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		os.Stdin = r

		want := true
		got := hasStdin()
		assert.Equal(t, want, got)
	})

	t.Run("with file as stdin", func(t *testing.T) {
		file := filetest.CreateFile(t, "")
		overrideStdin(t, file)

		want := false
		got := hasStdin()
		assert.Equal(t, want, got)
	})

	t.Run("with nil as stdin", func(t *testing.T) {
		overrideStdin(t, nil)

		want := false
		got := hasStdin()
		assert.Equal(t, want, got)
	})
}
