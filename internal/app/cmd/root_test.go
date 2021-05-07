package cmd

import (
	"errors"
	"io"
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/internal/pkg/console"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdRoot(t *testing.T) {
	cmd := newCmdRoot()
	assert.NotNil(t, cmd)
}

func TestPreRun_Quiet(t *testing.T) {
	context = ctx.NewCtx()
	context.Options.Quiet = true

	cmd := newCmdResolve()
	preRun(cmd, []string{})

	assert.Equal(t, console.Output, io.Discard)
}

func TestMust_OK(t *testing.T) {
	must(nil)
}

func TestMust_Panics(t *testing.T) {
	assert.Panics(t, func() { must(errors.New("error")) })
}
