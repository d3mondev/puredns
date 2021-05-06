package ctx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCtx(t *testing.T) {
	ctx := NewCtx()
	assert.NotNil(t, ctx)
}
