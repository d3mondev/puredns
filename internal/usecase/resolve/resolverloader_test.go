package resolve

import (
	"errors"
	"testing"
	"testing/iotest"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolverLoader(t *testing.T) {
	file := filetest.CreateFile(t, "")
	_, err := file.WriteString("8.8.8.8\n  \n1.1.1.1\n4.4.4.4")
	require.Nil(t, err)

	ctx := ctx.NewCtx()

	loader := NewDefaultResolverFileLoader()
	err = loader.Load(ctx, file.Name())

	assert.Nil(t, err)
	assert.ElementsMatch(t, ctx.Options.TrustedResolvers, []string{"8.8.8.8", "1.1.1.1", "4.4.4.4"})
}

func TestResolverLoaderFileOpenError(t *testing.T) {
	ctx := ctx.NewCtx()
	loader := NewDefaultResolverFileLoader()

	err := loader.Load(ctx, "thisfiledoesnotexit.txt")

	assert.NotNil(t, err)
}

func TestResolverScannerError(t *testing.T) {
	reader := iotest.ErrReader(errors.New("read error"))

	_, err := load(reader)

	assert.NotNil(t, err)
}
