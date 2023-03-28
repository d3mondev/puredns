package cmd

import (
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
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
}
