package resolve

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/pkg/fileoperation"
	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	context := ctx.NewCtx()
	opt := ctx.DefaultResolveOptions()

	svc := NewService(context, opt)

	assert.NotNil(t, svc)
}

func TestInitialize_OK(t *testing.T) {
	service, _ := newStubService(t)
	err := service.Initialize()
	assert.Nil(t, err)
}

func TestInitialize_RequirementError(t *testing.T) {
	service, stubs := newStubService(t)
	stubs.spyRequirementChecker.returns = errors.New("error")

	err := service.Initialize()

	assert.ErrorIs(t, err, stubs.spyRequirementChecker.returns)
}

func TestInitialize_WorkfilesError(t *testing.T) {
	service, stubs := newStubService(t)
	stubs.fakeWorkfileCreator.err = errors.New("error")

	err := service.Initialize()

	assert.ErrorIs(t, err, stubs.fakeWorkfileCreator.err)
}

func TestInitialize_PrepareResolversError(t *testing.T) {
	service, _ := newStubService(t)
	service.Options.ResolverFile = ""

	err := service.Initialize()
	assert.NotNil(t, err)

	service.Options.TrustedOnly = true
	err = service.Initialize()
	assert.Nil(t, err, "should not cause error when skipping public resolvers")
}

func TestClose(t *testing.T) {
	t.Run("without initialize", func(t *testing.T) {
		service, stubs := newStubService(t)
		service.Close(false)
		assert.Equal(t, 0, stubs.fakeWorkfileCreator.called)
	})

	t.Run("after initialize", func(t *testing.T) {
		service, stubs := newStubService(t)
		service.Initialize()
		service.Close(false)
		assert.Equal(t, 1, stubs.fakeWorkfileCreator.called)
	})
}

func TestResolve(t *testing.T) {
	context := ctx.NewCtx()
	opt := ctx.DefaultResolveOptions()
	opt.Mode = 0
	opt.TrustedOnly = true
	opt.DomainFile = filetest.CreateFile(t, "").Name()

	service := NewService(context, opt)
	require.Nil(t, service.Initialize())

	err := service.Resolve()
	assert.Nil(t, err)
}

func TestPrepareResolvers(t *testing.T) {
	service, _ := newStubService(t)
	service.workfiles = &Workfiles{}
	service.workfiles.PublicResolvers = filetest.CreateFile(t, "").Name()
	service.workfiles.TrustedResolvers = filetest.CreateFile(t, "").Name()

	service.Context.Options.TrustedResolvers = []string{"trusted"}
	service.Options.ResolverFile = filetest.CreateFile(t, "public").Name()

	require.Nil(t, service.prepareResolvers())

	gotPublic := filetest.ReadFile(t, service.workfiles.PublicResolvers)
	gotTrusted := filetest.ReadFile(t, service.workfiles.TrustedResolvers)

	assert.Equal(t, []string{"public"}, gotPublic, "public resolvers file should be populated")
	assert.Equal(t, []string{"trusted"}, gotTrusted, "trusted resolvers file should be populated")
}

func TestCreateDomainReaderSource_Stdin(t *testing.T) {
	service, _ := newStubService(t)
	service.Options.Mode = ctx.Resolve
	service.Context.Stdin = filetest.CreateFile(t, "stdin")
	service.Options.DomainFile = filetest.CreateFile(t, "file").Name()

	reader, err := service.createDomainReaderSource()
	assert.Nil(t, err)
	assert.Equal(t, 0, service.domainCount)

	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)

	assert.Equal(t, "stdin", string(buf[:n]), "should prioritize stdin")
}

func TestCreateDomainReaderSource_DomainFile(t *testing.T) {
	service, _ := newStubService(t)
	service.Options.Mode = ctx.Resolve
	service.Options.DomainFile = filetest.CreateFile(t, "example.com\n").Name()

	reader, err := service.createDomainReaderSource()
	assert.Nil(t, err)
	assert.Equal(t, 1, service.domainCount)

	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)

	assert.Equal(t, "example.com\n", string(buf[:n]))
}

func TestCreateDomainReaderSource_WordlistFile(t *testing.T) {
	service, _ := newStubService(t)
	service.Options.Mode = ctx.Bruteforce
	service.Options.Wordlist = filetest.CreateFile(t, "word\n").Name()

	reader, err := service.createDomainReaderSource()
	assert.Nil(t, err)
	assert.Equal(t, 1, service.domainCount)

	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)

	assert.Equal(t, "word\n", string(buf[:n]))
}

func TestCreateDomainReaderSource_FileError(t *testing.T) {
	service, _ := newStubService(t)
	service.Options.Mode = ctx.Bruteforce
	service.Options.Wordlist = ""

	_, err := service.createDomainReaderSource()
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestCreateDomainReader(t *testing.T) {
	service, _ := newStubService(t)
	service.Options.Mode = ctx.Bruteforce
	service.Options.Domain = "example.com"
	service.Options.Wordlist = filetest.CreateFile(t, "word\n").Name()

	reader, err := service.createDomainReader()
	assert.Nil(t, err)
	assert.Equal(t, 1, service.domainCount)

	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)

	assert.Equal(t, "word.example.com\n", string(buf[:n]))
}

func TestCreateDomainReader_MultipleDomains(t *testing.T) {
	service, _ := newStubService(t)
	service.Options.Mode = ctx.Bruteforce
	service.Options.DomainFile = filetest.CreateFile(t, "example.com\nexample.org").Name()
	service.Options.Wordlist = filetest.CreateFile(t, "word\n").Name()

	reader, err := service.createDomainReader()
	assert.Nil(t, err)
	assert.Equal(t, 2, service.domainCount)

	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)

	assert.Equal(t, "word.example.com\nword.example.org\n", string(buf[:n]))
}

func TestResolvePublic(t *testing.T) {
	defOpts := ctx.DefaultResolveOptions()
	publicResolverFile := filetest.CreateFile(t, "public")
	trustedResolverFile := filetest.CreateFile(t, "trusted")

	tests := []struct {
		name          string
		haveNoPublic  bool
		wantResolvers string
		wantRateLimit int
	}{
		{name: "ok", wantResolvers: "public", wantRateLimit: defOpts.RateLimit},
		{name: "nopublic option", haveNoPublic: true, wantResolvers: "trusted", wantRateLimit: defOpts.RateLimitTrusted},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, stubs := newStubService(t)
			service.ResolverLoader = NewDefaultResolverFileLoader()
			service.Options.ResolverFile = publicResolverFile.Name()
			service.Options.ResolverTrustedFile = trustedResolverFile.Name()
			service.Options.TrustedOnly = test.haveNoPublic
			require.Nil(t, service.Initialize())

			domainReader := NewDomainReader(io.NopCloser(strings.NewReader("")), nil, nil)

			err := service.resolvePublic(domainReader)
			gotResolvers := filetest.ReadFile(t, stubs.spyMassResolver.resolvers)

			assert.Nil(t, err)
			assert.Equal(t, []string{test.wantResolvers}, gotResolvers)
			assert.Equal(t, test.wantRateLimit, stubs.spyMassResolver.ratelimit)
		})
	}
}

func TestResolveTrusted(t *testing.T) {
	publicResolverFile := filetest.CreateFile(t, "public")

	tests := []struct {
		name               string
		haveSkipValidation bool
		wantResolvers      []string
	}{
		{name: "skip validation", haveSkipValidation: true, wantResolvers: []string{}},
		{name: "correct resolvers used", haveSkipValidation: false, wantResolvers: []string{"trusted"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, stubs := newStubService(t)
			service.Options.SkipValidation = test.haveSkipValidation
			service.Context.Options.TrustedResolvers = []string{"trusted"}
			service.Options.ResolverFile = publicResolverFile.Name()

			require.Nil(t, service.Initialize())
			err := service.resolveTrusted()
			require.Nil(t, err)

			content := filetest.ReadFile(t, stubs.spyMassResolver.resolvers)

			assert.Equal(t, test.wantResolvers, content)
		})
	}
}

func TestFilterWildcards(t *testing.T) {
	tests := []struct {
		name             string
		haveSkipWildcard bool
		wantCalled       int
	}{
		{name: "with wildcard filtering", wantCalled: 1},
		{name: "no wildcard filtering", haveSkipWildcard: true, wantCalled: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, stubs := newStubService(t)
			service.Options.SkipWildcard = test.haveSkipWildcard

			stubs.stubWildcardFilter.roots = []string{"root"}
			stubs.stubWildcardFilter.domains = []string{"example.com"}

			require.Nil(t, service.Initialize())
			require.Nil(t, fileoperation.WriteLines([]string{"example.com A 127.0.0.1"}, service.workfiles.MassdnsPublic))
			err := service.filterWildcards()
			require.Nil(t, err)

			assert.Equal(t, test.wantCalled, stubs.stubWildcardFilter.called)
			assert.Equal(t, 1, service.domainCount, "should be 1 in all cases")
		})
	}
}

func TestWriteResults(t *testing.T) {
	tests := []struct {
		name                 string
		haveCatError         bool
		haveResultSaverError error
		wantErr              bool
	}{
		{name: "ok", wantErr: false},
		{name: "filecat error handling", haveCatError: true, wantErr: true},
		{name: "resultsaver error handling", haveResultSaverError: errors.New("error"), wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, stubs := newStubService(t)
			stubs.stubResultSaver.returns = test.haveResultSaverError

			require.Nil(t, service.Initialize())

			if test.haveCatError {
				stubs.fakeWorkfileCreator.workfiles.Domains = ""
			}

			gotErr := service.writeResults()

			assert.Equal(t, test.wantErr, gotErr != nil, gotErr)
		})
	}
}
