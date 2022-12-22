package resolve

import (
	"io"
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/pkg/filetest"
)

type stubs struct {
	spyRequirementChecker *spyRequirementChecker
	fakeWorkfileCreator   *fakeWorkfileCreator
	spyResolverLoader     *spyResolverLoader
	stubDomainSanitizer   *stubDomainSanitizer
	spyMassResolver       *spyMassResolver
	stubWildcardFilter    *stubWildcardFilter
	stubResultSaver       *stubResultSaver
}

func newStubService(t *testing.T) (*Service, stubs) {
	stubs := stubs{
		spyRequirementChecker: &spyRequirementChecker{},
		fakeWorkfileCreator:   newFakeWorkfileCreator(t),
		spyResolverLoader:     &spyResolverLoader{},
		stubDomainSanitizer:   &stubDomainSanitizer{},
		spyMassResolver:       &spyMassResolver{},
		stubWildcardFilter:    &stubWildcardFilter{},
		stubResultSaver:       &stubResultSaver{},
	}

	service := &Service{
		Context: ctx.NewCtx(),
		Options: ctx.DefaultResolveOptions(),

		RequirementChecker: stubs.spyRequirementChecker,
		WorkfileCreator:    stubs.fakeWorkfileCreator,
		ResolverLoader:     stubs.spyResolverLoader,
		MassResolver:       stubs.spyMassResolver,
		WildcardFilter:     stubs.stubWildcardFilter,
		ResultSaver:        stubs.stubResultSaver,
	}

	service.Options.ResolverFile = filetest.CreateFile(t, "8.8.8.8").Name()

	t.Cleanup(func() {
		service.Close(false)
	})

	return service, stubs
}

type spyRequirementChecker struct {
	called  int
	returns error
}

func (s *spyRequirementChecker) Check(opt *ctx.ResolveOptions) error {
	s.called++
	return s.returns
}

func newFakeWorkfileCreator(t *testing.T) *fakeWorkfileCreator {
	return &fakeWorkfileCreator{t: t}
}

type fakeWorkfileCreator struct {
	t *testing.T

	workfiles *Workfiles
	called    int

	err error
}

func (f *fakeWorkfileCreator) Create() (*Workfiles, error) {
	f.called++

	if f.err != nil {
		return nil, f.err
	}

	realCreator := NewDefaultWorkfileCreator()

	files, err := realCreator.Create()
	if err != nil {
		f.t.Fatal(err)
	}

	f.workfiles = files

	return f.workfiles, nil
}

type spyResolverLoader struct {
	called int
	err    error
}

func (s *spyResolverLoader) Load(*ctx.Ctx, string) error {
	s.called++
	return s.err
}

type spyMassResolver struct {
	called    int
	resolvers string
	ratelimit int
}

func (s *spyMassResolver) Resolve(r io.Reader, output string, total int, resolvers string, qps int) error {
	s.called++
	s.resolvers = resolvers
	s.ratelimit = qps
	return nil
}

func (s *spyMassResolver) Current() int {
	return 0
}

func (s *spyMassResolver) Rate() float64 {
	return 0.0
}

type stubWildcardFilter struct {
	called  int
	err     error
	domains []string
	roots   []string
}

func (s *stubWildcardFilter) Filter(WildcardFilterOptions, int) (found int, roots []string, err error) {
	s.called++
	return len(s.domains), s.roots, s.err
}

type stubDomainSanitizer struct {
	called  int
	returns error
}

func (s *stubDomainSanitizer) Sanitize(string, string) error {
	s.called++
	return s.returns
}

type stubResultSaver struct {
	called  int
	returns error
}

func (s *stubResultSaver) Save(*Workfiles, *ctx.ResolveOptions) error {
	s.called++
	return s.returns
}
