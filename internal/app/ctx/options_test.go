package ctx

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	have.NoPublicResolvers = true

	want := DefaultResolveOptions()
	want.NoPublicResolvers = true
	want.SkipValidation = true

	err := have.Validate()

	assert.Nil(t, err)
	assert.Equal(t, want, have)
}
