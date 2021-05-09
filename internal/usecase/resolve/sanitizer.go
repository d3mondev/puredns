package resolve

import (
	"strings"
)

// DefaultSanitizer is the default sanitizer function. It transforms the domain to lower characters,
// and ensures only valid characters are present. Returns an empty string if the domain fails sanitization.
func DefaultSanitizer(domain string) string {
	// Set to lowercase
	domain = strings.ToLower(domain)

	// Remove *.
	domain = strings.TrimPrefix(domain, "*.")

	// Keep only domains containing [a-z0-9.-]
	// Faster than using a regular expression
	for i := 0; i < len(domain); i++ {
		char := domain[i]

		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || (char == '-') || (char == '.') {
			continue
		}

		domain = ""
		break
	}

	return domain
}
