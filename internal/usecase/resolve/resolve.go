package resolve

import (
	"fmt"
	"io"
	"os"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/internal/pkg/console"
	"github.com/d3mondev/puredns/v2/pkg/fileoperation"
	"github.com/d3mondev/puredns/v2/pkg/shellexecutor"
)

// Service contains the interfaces required to operate the service.
type Service struct {
	Context *ctx.Ctx
	Options *ctx.ResolveOptions

	RequirementChecker RequirementChecker
	ResolverLoader     ResolverLoader
	WorkfileCreator    WorkfileCreator
	MassResolver       MassResolver
	ResultSaver        ResultSaver
	WildcardFilter     WildcardFilter

	workfiles   *Workfiles
	domainCount int
}

// RequirementChecker checks if the dependencies are present on the system.
type RequirementChecker interface {
	Check(opt *ctx.ResolveOptions) error
}

// ResolverLoader loads the resolvers from a file and put them into the application context.
type ResolverLoader interface {
	Load(ctx *ctx.Ctx, filename string) error
}

// WorkfileCreator creates a new set of workfiles used during the program execution.
type WorkfileCreator interface {
	Create() (*Workfiles, error)
}

// MassResolver resolves the domains contained in the input file using the resolvers present in the resolvers file
// and saves the results. Queries per second can be limited by setting the qps argument (0 for unlimited).
type MassResolver interface {
	Resolve(reader io.Reader, output string, total int, resolversFilename string, qps int) error
}

// ResultSaver saves the results as direction by the options specified.
type ResultSaver interface {
	Save(workfiles *Workfiles, opt *ctx.ResolveOptions) error
}

// WildcardFilter filters the wildcard subdomains from a list of domains.
type WildcardFilter interface {
	Filter(opt WildcardFilterOptions, totalCount int) (found int, roots []string, err error)
}

// NewService creates a new ResolveService object.
func NewService(ctx *ctx.Ctx, opt *ctx.ResolveOptions) *Service {
	service := Service{
		Context: ctx,
		Options: opt,

		RequirementChecker: NewDefaultRequirementChecker(shellexecutor.NewShellExecutor()),
		ResolverLoader:     NewDefaultResolverFileLoader(),
		WorkfileCreator:    NewDefaultWorkfileCreator(),
		MassResolver:       NewDefaultMassResolver(opt.BinPath),
		ResultSaver:        NewResultFileSaver(),
		WildcardFilter:     NewDefaultWildcardFilter(),
	}

	return &service
}

// Initialize makes sure that the required binaries can be run and that all the required files are present.
func (s *Service) Initialize() error {
	var err error

	// Check if required binaries are present
	if err := s.RequirementChecker.Check(s.Options); err != nil {
		return err
	}

	// Create the temporary workfiles
	if s.workfiles, err = s.WorkfileCreator.Create(); err != nil {
		return err
	}

	// Prepare resolvers
	if err := s.prepareResolvers(); err != nil {
		return err
	}

	return nil
}

// Resolve resolves domain names contained in a file and saves the output according to the program context.
func (s *Service) Resolve() error {
	var err error

	// Create the domain reader. The reader is responsible for constructing a list of subdomains to
	// resolve from a list of domains or a list of words that can either be read from a file or from stdin.
	// This reader implements the io.ReadCloser interface and can be passed directly to a program's stdin.
	domainReader, err := s.createDomainReader()
	if err != nil {
		return err
	}

	// Resolve the domains from the domain reader using public resolvers
	if err = s.resolvePublic(domainReader); err != nil {
		return err
	}

	// Filter out the wildcard domains using the DNS cache that was produced by the earlier call to massdns.
	// The cache can use quite a bit of memory when resolving millions of domains at once.
	if err = s.filterWildcards(); err != nil {
		return err
	}

	// Resolve the remaining domains using trusted resolvers to attempt to filter-out any DNS poisoning
	if err = s.resolveTrusted(); err != nil {
		return err
	}

	// Write the results to stdout and to files
	if err = s.writeResults(); err != nil {
		return err
	}

	return nil
}

// Close terminates the service.
func (s *Service) Close(debug bool) {
	if debug {
		console.Printf("\nDebug files kept in: %s\n", s.workfiles.TempDirectory)
	} else if s.workfiles != nil {
		s.workfiles.Close()
	}
}

func (s *Service) prepareResolvers() error {
	// Create a copy of the public resolvers in a temporary directory. This is to ensure that the resolvers specified
	// don't get changed by an external process during the operation.
	if !s.Options.TrustedOnly {
		if err := fileoperation.Copy(s.Options.ResolverFile, s.workfiles.PublicResolvers); err != nil {
			return fmt.Errorf("unable to load public resolvers: %w", err)
		}
	}

	// If custom trusted resolvers are specified, load them from a file
	if err := s.ResolverLoader.Load(s.Context, s.Options.ResolverTrustedFile); err != nil {
		return fmt.Errorf("unable to load trusted resolvers: %w", err)
	}

	// Create a copy of the trusted resolvers in a temporary directory. This allows saving the default resolvers to
	// disk to be used by massdns.
	if err := fileoperation.WriteLines(s.Context.Options.TrustedResolvers, s.workfiles.TrustedResolvers); err != nil {
		return fmt.Errorf("unable to write trusted resolvers to temporary directory: %w", err)
	}

	return nil
}

func (s *Service) createDomainReader() (*DomainReader, error) {
	// Create a reader for the source words or domains, depending on the mode
	sourceReader, err := s.createDomainReaderSource()
	if err != nil {
		return nil, err
	}

	// Create a list of domains to process in bruteforce mode, otherwise it's nil
	var domains []string
	if s.Options.Mode == ctx.Bruteforce {
		if domains, err = s.createDomainReaderDomainList(); err != nil {
			return nil, err
		}
	}

	// Create a sanitizer if needed, otherwise it's nil
	var sanitizer DomainSanitizer
	if !s.Options.SkipSanitize {
		sanitizer = DefaultSanitizer
	}

	r := NewDomainReader(sourceReader, domains, sanitizer)

	return r, nil
}

func (s *Service) createDomainReaderSource() (io.ReadCloser, error) {
	var sourceReader io.ReadCloser

	if s.Context.Stdin != nil {
		// Use stdin first if present
		sourceReader = s.Context.Stdin
	} else {
		// Open the filename containing the source data
		var filename string
		if s.Options.Mode == ctx.Resolve {
			filename = s.Options.DomainFile
		} else {
			filename = s.Options.Wordlist
		}

		// Count the number of lines to get a total for the progress bar
		count, err := fileoperation.CountLines(filename)
		if err != nil {
			return nil, err
		}

		// File will be closed by the DomainReader when it has finished reading
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		sourceReader = file
		s.domainCount = count
	}

	return sourceReader, nil
}

func (s *Service) createDomainReaderDomainList() ([]string, error) {
	var domains []string

	if s.Options.DomainFile != "" {
		// Read domains from file
		var err error
		domains, err = fileoperation.ReadLines(s.Options.DomainFile)
		if err != nil {
			return nil, err
		}
	} else {
		// Use domain from options
		domains = []string{s.Options.Domain}
	}

	// Multiply the count by the actual number of domains to test
	s.domainCount = s.domainCount * len(domains)

	return domains, nil
}

func (s *Service) resolvePublic(reader *DomainReader) error {
	resolvers := s.workfiles.PublicResolvers
	ratelimit := s.Options.RateLimit
	resolverString := "public"

	if s.Options.TrustedOnly {
		resolvers = s.workfiles.TrustedResolvers
		ratelimit = s.Options.RateLimitTrusted
		resolverString = "trusted"
	}

	console.Printf("%sResolving domains with %s resolvers%s\n", console.ColorBrightWhite, resolverString, console.ColorReset)

	err := s.MassResolver.Resolve(
		reader,
		s.workfiles.MassdnsPublic,
		s.domainCount,
		resolvers,
		ratelimit,
	)

	if err != nil {
		return fmt.Errorf("error resolving domains: %w", err)
	}

	console.Printf("\n")

	return err
}

func (s *Service) filterWildcards() error {
	// If we're skipping wildcard filtering, we still need to produce a list of valid
	// domain names for the next step and update the domain count
	if s.Options.SkipWildcard {
		return s.parseCache(s.workfiles.MassdnsPublic, s.workfiles.Domains)
	}

	// Parsing the cache without a filename only updates the domain count for the progress bar
	if err := s.parseCache(s.workfiles.MassdnsPublic, ""); err != nil {
		return err
	}

	console.Printf("%sDetecting wildcard root subdomains%s\n", console.ColorBrightWhite, console.ColorReset)

	opt := WildcardFilterOptions{
		CacheFilename:        s.workfiles.MassdnsPublic,
		DomainOutputFilename: s.workfiles.Domains,
		RootOutputFilename:   s.workfiles.WildcardRoots,
		Resolvers:            s.Context.Options.TrustedResolvers,
		QueriesPerSecond:     s.Options.RateLimitTrusted,
		ThreadCount:          s.Options.WildcardThreads,
		ResolveTestCount:     s.Options.WildcardTests,
		BatchSize:            s.Options.WildcardBatchSize,
	}

	found, roots, err := s.WildcardFilter.Filter(opt, s.domainCount)

	if err != nil {
		return fmt.Errorf("unable to filter wildcard domains: %w", err)
	}

	if len(roots) > 0 {
		console.Printf("\n%sFound %s%d%s wildcard roots:%s\n",
			console.ColorBrightWhite,
			console.ColorBrightGreen,
			len(roots),
			console.ColorBrightWhite,
			console.ColorReset,
		)

		for _, root := range roots {
			console.Printf("*.%s\n", root)
		}
	}

	s.domainCount = found

	console.Printf("\n")

	return nil
}

// parseCache parses the massdns cache file to count valid domains and save them to a file
// if the filename is specified.
func (s *Service) parseCache(cacheFilename string, domainFilename string) error {
	var domainFile *os.File
	cacheFile, err := os.Open(cacheFilename)
	if err != nil {
		return err
	}

	// If we're skipping wildcard detection, we also need to extract domains
	// for the validation step
	if domainFilename != "" {
		if domainFile, err = os.Create(domainFilename); err != nil {
			return err
		}
	}

	cacheReader := NewCacheReader(cacheFile)

	if s.domainCount, err = cacheReader.Read(domainFile, nil, 0); err != nil {
		return err
	}

	// If we're also saving domains, make sure to sync the file to disk
	if domainFilename != "" {
		if err := domainFile.Sync(); err != nil {
			return err
		}
		cacheReader.Close()
		return domainFile.Close()
	}

	return cacheReader.Close()
}

func (s *Service) resolveTrusted() error {
	if s.Options.SkipValidation {
		return nil
	}

	domainFile, err := os.Open(s.workfiles.Domains)
	if err != nil {
		return nil
	}
	defer domainFile.Close()

	console.Printf("%sValidating domains against trusted resolvers%s\n", console.ColorBrightWhite, console.ColorReset)

	err = s.MassResolver.Resolve(
		domainFile,
		s.workfiles.MassdnsTrusted,
		s.domainCount,
		s.workfiles.TrustedResolvers,
		s.Options.RateLimitTrusted,
	)

	if err != nil {
		return fmt.Errorf("error resolving domains: %w", err)
	}

	console.Printf("\n")

	return s.parseCache(s.workfiles.MassdnsTrusted, s.workfiles.Domains)
}

func (s *Service) writeResults() error {
	if s.domainCount > 0 {
		console.Printf("%sFound %s%d%s valid domains:%s\n",
			console.ColorBrightWhite,
			console.ColorBrightGreen,
			s.domainCount,
			console.ColorBrightWhite,
			console.ColorReset)
	} else {
		console.Printf("\n%sNo valid domains remaining.%s\n", console.ColorBrightWhite, console.ColorReset)
	}

	if err := fileoperation.Cat([]string{s.workfiles.Domains}, os.Stdout); err != nil {
		return fmt.Errorf("unable to read domain file: %w", err)
	}

	if err := s.ResultSaver.Save(s.workfiles, s.Options); err != nil {
		return fmt.Errorf("unable to save results: %w", err)
	}

	return nil
}
