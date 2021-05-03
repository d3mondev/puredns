package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCmdSponsors(t *testing.T) {
	cmd := newCmdSponsors()
	assert.NotNil(t, cmd)
}
