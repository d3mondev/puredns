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
	opt := ctx.NewResolveOptions()

	svc := NewService(context, opt)

	assert.NotNil(t, svc)
}

func TestInitialize(t *testing.T) {
	tests := []struct {
		name                      string
		haveRequirementError      error
		haveWorkfileError         error
		havePrepareResolversError bool
		wantErr                   bool
	}{
		{name: "ok"},
		{name: "requirements error handling", haveRequirementError: errors.New("error"), wantErr: true},
		{name: "workfiles error handling", haveWorkfileError: errors.New("error"), wantErr: true},
		{name: "prepareresolvers error handling", havePrepareResolversError: true, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, stubs := newStubService(t)
			stubs.spyRequirementChecker.returns = test.haveRequirementError
			stubs.fakeWorkfileCreator.err = test.haveWorkfileError

			if test.havePrepareResolversError {
				service.Options.ResolverFile = ""
			}

			got := service.Initialize()

			assert.Equal(t, test.wantErr, got != nil, got)
		})
	}
}

func TestClose(t *testing.T) {
	t.Run("without initialize", func(t *testing.T) {
		service, stubs := newStubService(t)
		service.Close()
		assert.Equal(t, 0, stubs.fakeWorkfileCreator.called)
	})

	t.Run("after initialize", func(t *testing.T) {
		service, stubs := newStubService(t)
		service.Initialize()
		service.Close()
		assert.Equal(t, 1, stubs.fakeWorkfileCreator.called)
	})
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

func TestCreateDomainReader_Stdin(t *testing.T) {
	service, _ := newStubService(t)
	service.Context.Stdin = filetest.CreateFile(t, "stdin")
	service.Options.DomainFile = filetest.CreateFile(t, "file").Name()

	reader, err := service.createDomainReader()
	assert.Nil(t, err)
	assert.Equal(t, 0, service.domainCount)

	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)

	assert.Equal(t, "stdin\n", string(buf[:n]), "should prioritize stdin")
}

func TestCreateDomainReader_File(t *testing.T) {
	service, _ := newStubService(t)
	service.LineCounter = fileoperation.CountLines
	service.Options.DomainFile = filetest.CreateFile(t, "file\n").Name()

	reader, err := service.createDomainReader()
	assert.Nil(t, err)
	assert.Equal(t, 1, service.domainCount)

	buf := make([]byte, 1024)
	n, _ := reader.Read(buf)

	assert.Equal(t, "file\n", string(buf[:n]))
}

func TestCreateDomainReader_WordlistError(t *testing.T) {
	service, _ := newStubService(t)
	service.Options.Mode = 1
	service.Options.Wordlist = ""

	_, err := service.createDomainReader()
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestCreateDomainReader_CountLinesError(t *testing.T) {
	err := errors.New("error")

	service, stubs := newStubService(t)
	stubs.stubLineCounter.err = err
	service.Options.DomainFile = filetest.CreateFile(t, "file\n").Name()

	_, err = service.createDomainReader()
	assert.ErrorIs(t, err, err)
}

func TestResolvePublic(t *testing.T) {
	publicResolverFile := filetest.CreateFile(t, "public")
	trustedResolverFile := filetest.CreateFile(t, "trusted")

	service, stubs := newStubService(t)
	service.Options.ResolverFile = publicResolverFile.Name()
	service.Options.ResolverTrustedFile = trustedResolverFile.Name()
	require.Nil(t, service.Initialize())

	domainReader := NewDomainReader(io.NopCloser(strings.NewReader("")), "", nil)

	err := service.resolvePublic(domainReader)
	used := filetest.ReadFile(t, stubs.spyMassResolver.resolvers)

	assert.Nil(t, err)
	assert.Equal(t, []string{"public"}, used)
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
			require.Nil(t, fileoperation.WriteLines([]string{"example.com A 127.0.0.1"}, service.workfiles.Massdns))
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
