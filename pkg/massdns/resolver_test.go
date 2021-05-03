package massdns

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubRunner struct {
	returns error
}

func (r stubRunner) Run(lineReader io.Reader, output string, resolvers string, qps int) error {
	buf := make([]byte, 1024)

	for {
		_, err := lineReader.Read(buf)
		if err != nil {
			break
		}
	}

	return r.returns
}

func TestResolve(t *testing.T) {
	tests := []struct {
		name            string
		haveRunnerError error
		wantErr         bool
	}{
		{name: "success"},
		{name: "runner error handling", haveRunnerError: errors.New("runner error"), wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resolver := NewResolver("massdns")
			resolver.runner = stubRunner{returns: test.haveRunnerError}

			gotErr := resolver.Resolve(strings.NewReader(""), "", "", 10)

			assert.Equal(t, test.wantErr, gotErr != nil, gotErr)
		})
	}
}

func TestCurrent(t *testing.T) {
	resolver := NewResolver("massdns")
	resolver.runner = stubRunner{}

	gotCurrent := resolver.Current()
	assert.Equal(t, 0, gotCurrent)

	resolver.Resolve(strings.NewReader("example.com\n"), "", "", 0)
	gotCurrent = resolver.Current()
	assert.Equal(t, 1, gotCurrent)
}
