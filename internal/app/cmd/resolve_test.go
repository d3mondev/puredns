package cmd

import (
	"os"
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
)

func TestRunResolve(t *testing.T) {
	t.Run("no argument", func(t *testing.T) {
		context = ctx.NewCtx()
		cmd := newCmdResolve()

		err := runResolve(cmd, []string{})

		assert.NotNil(t, err)
	})

	t.Run("file that does not exist", func(t *testing.T) {
		context = ctx.NewCtx()
		cmd := newCmdResolve()

		err := runResolve(cmd, []string{"thisfiledoesnotexist.txt"})

		assert.NotNil(t, err)
	})

	t.Run("with stdin", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		err = w.Close()
		if err != nil {
			t.Fatal(err)
		}

		filetest.OverrideStdin(t, r)

		context = ctx.NewCtx()
		cmd := newCmdResolve()

		err = runResolve(cmd, []string{})

		assert.NotNil(t, err)
	})
}
