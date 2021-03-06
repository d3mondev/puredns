package programbanner

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/internal/pkg/console"
	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	buffer := new(bytes.Buffer)
	console.Output = buffer

	ctx := ctx.NewCtx()
	service := NewService(ctx)
	service.Print()

	assert.True(t, strings.Contains(buffer.String(), ctx.ProgramName))
	assert.True(t, strings.Contains(buffer.String(), ctx.ProgramVersion))
}

func TestPrintWithResolveOptions(t *testing.T) {
	tests := []struct {
		name     string
		haveCtx  ctx.Ctx
		haveOpts ctx.ResolveOptions
		want     string
	}{
		{name: "stdin", haveCtx: ctx.Ctx{Stdin: os.Stdin}, haveOpts: ctx.ResolveOptions{}, want: "stdin"},
		{name: "resolve mode", haveCtx: ctx.Ctx{}, haveOpts: ctx.ResolveOptions{DomainFile: "domains.txt", Mode: 0}, want: "domains.txt"},
		{name: "bruteforce mode", haveCtx: ctx.Ctx{}, haveOpts: ctx.ResolveOptions{Wordlist: "wordlist.txt", Mode: 1}, want: "wordlist.txt"},
		{name: "trusted resolvers", haveCtx: ctx.Ctx{}, haveOpts: ctx.ResolveOptions{ResolverTrustedFile: "trusted.txt"}, want: "trusted.txt"},
		{name: "rate", haveCtx: ctx.Ctx{}, haveOpts: ctx.ResolveOptions{RateLimit: 777}, want: "777"},
		{name: "batch size", haveCtx: ctx.Ctx{}, haveOpts: ctx.ResolveOptions{WildcardBatchSize: 5555}, want: "5555"},
		{name: "write domains", haveCtx: ctx.Ctx{}, haveOpts: ctx.ResolveOptions{WriteDomainsFile: "domains_out.txt"}, want: "domains_out.txt"},
		{name: "write massdns", haveCtx: ctx.Ctx{}, haveOpts: ctx.ResolveOptions{WriteMassdnsFile: "massdns_out.txt"}, want: "massdns_out.txt"},
		{name: "write wildcards", haveCtx: ctx.Ctx{}, haveOpts: ctx.ResolveOptions{WriteWildcardsFile: "wildcards_out.txt"}, want: "wildcards_out.txt"},
		{name: "skip sanitize", haveCtx: ctx.Ctx{}, haveOpts: ctx.ResolveOptions{SkipSanitize: true}, want: "Skip Sanitize"},
		{name: "skip wildcard", haveCtx: ctx.Ctx{}, haveOpts: ctx.ResolveOptions{SkipWildcard: true}, want: "Skip Wildcard"},
		{name: "skip validation", haveCtx: ctx.Ctx{}, haveOpts: ctx.ResolveOptions{SkipValidation: true}, want: "Skip Validation"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buffer := new(bytes.Buffer)
			console.Output = buffer

			service := NewService(&test.haveCtx)
			service.PrintWithResolveOptions(&test.haveOpts)

			assert.Truef(t, strings.Contains(buffer.String(), test.want), "%s not found in output", test.want)
		})
	}
}
