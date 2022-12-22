package resolve

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMassResolverNew(t *testing.T) {
	r := NewDefaultMassResolver("")
	assert.NotNil(t, r)
}

func TestMassResolverResolve_OK(t *testing.T) {
	r := NewDefaultMassResolver("")

	err := r.Resolve(strings.NewReader("example.com"), "", 0, "", 10)
	assert.EqualError(t, err, "exec: no command", "should not call massdns because of invalid path")
}

func TestMassResolverResolve_WithTotal(t *testing.T) {
	r := NewDefaultMassResolver("")

	err := r.Resolve(strings.NewReader("example.com"), "", 100, "", 10)
	assert.EqualError(t, err, "exec: no command")
}
