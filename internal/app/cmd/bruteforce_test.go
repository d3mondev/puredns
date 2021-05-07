package cmd

import (
	"io/fs"
	"os"
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBruteforceArgs_TwoArgs(t *testing.T) {
	context = ctx.NewCtx()
	resolveOptions = ctx.DefaultResolveOptions()

	parseBruteforceArgs([]string{"wordlist.txt", "example.com"})

	err := resolveOptions.Validate()
	assert.Nil(t, err)
}

func TestParseBruteforceArgs_NoArgs(t *testing.T) {
	context = ctx.NewCtx()
	resolveOptions = ctx.DefaultResolveOptions()

	parseBruteforceArgs([]string{})

	err := resolveOptions.Validate()
	assert.ErrorIs(t, err, ctx.ErrNoDomain)
}

func TestParseBruteforceArgs_NoDomain(t *testing.T) {
	context = ctx.NewCtx()
	resolveOptions = ctx.DefaultResolveOptions()

	parseBruteforceArgs([]string{"wordlist.txt"})

	err := resolveOptions.Validate()
	assert.ErrorIs(t, err, ctx.ErrNoDomain)
}

func TestParseBruteforceArgs_DomainFile(t *testing.T) {
	context = ctx.NewCtx()
	resolveOptions = ctx.DefaultResolveOptions()
	resolveOptions.DomainFile = "domains.txt"

	parseBruteforceArgs([]string{"wordlist.txt"})

	err := resolveOptions.Validate()
	assert.Nil(t, err)
}

func TestParseBruteforceArgs_Stdin(t *testing.T) {
	context = ctx.NewCtx()
	resolveOptions = ctx.DefaultResolveOptions()

	r, w, err := os.Pipe()
	require.Nil(t, err)
	require.Nil(t, w.Close())
	filetest.OverrideStdin(t, r)

	parseBruteforceArgs([]string{"domain.com"})

	err = resolveOptions.Validate()
	assert.Nil(t, err)
}

func TestParseBruteforceArgs_StdinNoDomain(t *testing.T) {
	context = ctx.NewCtx()
	resolveOptions = ctx.DefaultResolveOptions()

	r, w, err := os.Pipe()
	require.Nil(t, err)
	require.Nil(t, w.Close())
	filetest.OverrideStdin(t, r)

	parseBruteforceArgs([]string{})

	err = resolveOptions.Validate()
	assert.ErrorIs(t, err, ctx.ErrNoDomain)
}

func TestParseBruteforceArgs_StdinDomainFile(t *testing.T) {
	context = ctx.NewCtx()
	resolveOptions = ctx.DefaultResolveOptions()
	resolveOptions.DomainFile = "domains.txt"

	r, w, err := os.Pipe()
	require.Nil(t, err)
	require.Nil(t, w.Close())
	filetest.OverrideStdin(t, r)

	parseBruteforceArgs([]string{})

	err = resolveOptions.Validate()
	assert.Nil(t, err)
}

func TestRunBruteforce_OK(t *testing.T) {
	resolvers := filetest.CreateFile(t, "8.8.8.8\n")
	wordlist := filetest.CreateFile(t, "")

	context = ctx.NewCtx()
	resolveOptions = ctx.DefaultResolveOptions()
	resolveOptions.ResolverFile = resolvers.Name()

	cmd := newCmdBruteforce()
	err := runBruteforce(cmd, []string{wordlist.Name(), "example.com"})
	assert.Nil(t, err)
}

func TestRunBruteforce_ValidateError(t *testing.T) {
	context = ctx.NewCtx()
	resolveOptions = ctx.DefaultResolveOptions()

	cmd := newCmdBruteforce()
	err := runBruteforce(cmd, []string{})
	assert.ErrorIs(t, err, ctx.ErrNoDomain)
}

func TestRunBruteforce_InitializeError(t *testing.T) {
	context = ctx.NewCtx()
	resolveOptions = ctx.DefaultResolveOptions()

	cmd := newCmdBruteforce()
	err := runBruteforce(cmd, []string{"wordlist.txt", "example.com"})
	assert.ErrorIs(t, err, fs.ErrNotExist)
}
