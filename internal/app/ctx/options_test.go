package ctx

import (
	"os"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultGlobalOptions(t *testing.T) {
	opts := DefaultGlobalOptions()
	assert.NotNil(t, opts)
}

func TestResolveOptionsValidate_OK(t *testing.T) {
	have := DefaultResolveOptions()
	want := DefaultResolveOptions()

	err := have.Validate()

	assert.Nil(t, err)
	assert.Equal(t, want, have)
}

func TestResolveOptionsValidate_NoPublic(t *testing.T) {
	have := DefaultResolveOptions()
	have.TrustedOnly = true

	want := DefaultResolveOptions()
	want.TrustedOnly = true
	want.SkipValidation = true

	err := have.Validate()

	assert.Nil(t, err)
	assert.Equal(t, want, have)
}

func TestResolveOptionsValidate_BruteforceNoDomain(t *testing.T) {
	have := DefaultResolveOptions()
	have.Mode = Bruteforce
	have.Wordlist = "wordlist.txt"

	err := have.Validate()
	assert.ErrorIs(t, ErrNoDomain, err)
}

func TestResolveOptionsValidate_BruteforceDomain(t *testing.T) {
	have := DefaultResolveOptions()
	have.Mode = Bruteforce
	have.Wordlist = "wordlist.txt"
	have.Domain = "example.com"

	err := have.Validate()

	assert.Nil(t, err)
}

func TestResolveOptionsValidate_BruteforceDomainFile(t *testing.T) {
	have := DefaultResolveOptions()
	have.Mode = Bruteforce
	have.Wordlist = "wordlist.txt"
	have.DomainFile = "domains.txt"

	err := have.Validate()

	assert.Nil(t, err)
}

func TestResolveOptionsValidate_BruteforceNoWordlist(t *testing.T) {
	have := DefaultResolveOptions()
	have.Mode = Bruteforce
	have.Domain = "example.com"

	err := have.Validate()

	assert.ErrorIs(t, ErrNoWordlist, err)
}

func TestResolveOptionsValidate_BruteforceWordlistStdin(t *testing.T) {
	have := DefaultResolveOptions()
	have.Mode = Bruteforce
	have.Domain = "example.com"

	r, w, err := os.Pipe()
	require.Nil(t, err)
	require.Nil(t, w.Close())
	filetest.OverrideStdin(t, r)

	err = have.Validate()

	assert.Nil(t, err)
}
