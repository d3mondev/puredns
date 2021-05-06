package cmd

import (
	"os"
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/stretchr/testify/assert"
)

func TestRunBruteforce(t *testing.T) {
	context = ctx.NewCtx()
	resolveOptions = ctx.DefaultResolveOptions()

	t.Run("missing wordlist", func(t *testing.T) {
		cmd := newCmdBruteforce()
		err := runBruteforce(cmd, []string{"domain.com"})
		assert.NotNil(t, err)
	})

	t.Run("file that does not exist", func(t *testing.T) {
		cmd := newCmdBruteforce()
		err := runBruteforce(cmd, []string{"thisfiledoesnotexist.txt", "domain.com"})
		assert.NotNil(t, err)
	})

	t.Run("with stdin", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		err = w.Close()
		if err != nil {
			t.Fatal(err)
		}

		overrideStdin(t, r)

		cmd := newCmdBruteforce()
		err = runBruteforce(cmd, []string{"domain.com"})
		assert.NotNil(t, err)
	})
}
