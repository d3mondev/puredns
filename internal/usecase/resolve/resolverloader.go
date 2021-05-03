package resolve

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
)

// DefaultResolverLoader loads resolvers from a text file.
type DefaultResolverLoader struct{}

// NewDefaultResolverFileLoader creates a new ResolverFileLoader instance.
func NewDefaultResolverFileLoader() *DefaultResolverLoader {
	return &DefaultResolverLoader{}
}

// Load parses the specified filename to load resolvers and saves them to the program context.
func (l *DefaultResolverLoader) Load(ctx *ctx.Ctx, filename string) error {
	if filename == "" {
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	resolvers, err := load(file)

	if len(resolvers) > 0 {
		ctx.Options.TrustedResolvers = resolvers
	}

	return err
}

func load(r io.Reader) ([]string, error) {
	resolvers := []string{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		resolver := strings.TrimSpace(scanner.Text())
		if resolver == "" {
			continue
		}

		resolvers = append(resolvers, resolver)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return resolvers, nil
}
