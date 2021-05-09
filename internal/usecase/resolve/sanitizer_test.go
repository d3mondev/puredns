package resolve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSanitizer(t *testing.T) {
	tests := []struct {
		name       string
		haveDomain string
		wantDomain string
	}{
		{name: "valid domain", haveDomain: "example.com", wantDomain: "example.com"},
		{name: "tolower transform", haveDomain: "EXAMPLE.COM", wantDomain: "example.com"},
		{name: "invalid characters", haveDomain: "example+.com", wantDomain: ""},
		{name: "wildcard", haveDomain: "*.example.com", wantDomain: "example.com"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := DefaultSanitizer(test.haveDomain)
			assert.Equal(t, test.wantDomain, got)
		})
	}
}
