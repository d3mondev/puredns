package resolve

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/d3mondev/puredns/v2/internal/pkg/console"
	"github.com/d3mondev/puredns/v2/pkg/fileoperation"
	"github.com/d3mondev/puredns/v2/pkg/progressbar"
	"github.com/d3mondev/puredns/v2/pkg/wildcarder"
)

// WildcardFilterOptions defines the options for the Filter function.
type WildcardFilterOptions struct {
	// Input files
	CacheFilename string

	// Output files
	DomainOutputFilename string
	RootOutputFilename   string

	// Filtering parameters
	Resolvers        []string
	QueriesPerSecond int
	ThreadCount      int
	ResolveTestCount int
	BatchSize        int
}

// DefaultWildcardFilter implements the WildcardFilter interface used to filter wildcards.
type DefaultWildcardFilter struct {
	wc *wildcarder.Wildcarder
}

// NewDefaultWildcardFilter returns a new DefaultWildcardFilter object.
func NewDefaultWildcardFilter() *DefaultWildcardFilter {
	return &DefaultWildcardFilter{}
}

// Filter returns the number of domains that are not wildcards along with the wildcard roots found. It uses the massdns cache file
// to prepopulate a cache of DNS responses to optimize the number of DNS queries to perform. It saves the results to the specified filenames.
func (f *DefaultWildcardFilter) Filter(opt WildcardFilterOptions, totalCount int) (found int, roots []string, err error) {
	// Create the cache file reader
	cacheReader, err := createCacheReader(opt.CacheFilename)
	if err != nil {
		return 0, nil, err
	}
	defer cacheReader.Close()

	// Create wildcarder
	f.wc = createWildcarder(opt)

	// Create temporary file
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		return 0, nil, err
	}
	defer func() { tempFile.Close(); os.Remove(tempFile.Name()) }()

	// Start progress bar
	tmpl := "[ETA {{ eta }}] {{ bar }} {{ current }}/{{ total }} queries: {{ queries }} (time: {{ time }})"
	bar := progressbar.New(f.updateProgressBar, int64(totalCount), progressbar.WithTemplate(tmpl), progressbar.WithWriter(console.Output))
	bar.Start()

	// Process entries in batch to prevent the precache from taking too much memory on very large (70M+)
	// domain lists. The wildcard cache and DNS cache stay intact and keep growing between batches for now.
	rootMap := make(map[string]struct{})
	for {
		// Load precache batch
		precache, domainFile, count, err := prepareCache(cacheReader, tempFile.Name(), opt.BatchSize)
		if err != nil {
			return 0, nil, err
		}

		// Nothing to process, we're done!
		if count == 0 {
			break
		}

		// Set current precache
		f.wc.SetPreCache(precache)

		// Filter wildcards
		domains, roots := f.wc.Filter(domainFile)
		domainFile.Close()
		found += len(domains)

		// Save domains found
		if err := fileoperation.AppendLines(domains, opt.DomainOutputFilename); err != nil {
			return 0, nil, err
		}

		// Keep unique roots in map
		for _, root := range roots {
			rootMap[root] = struct{}{}
		}
	}

	// Save roots found
	var rootList []string
	for root := range rootMap {
		rootList = append(rootList, root)
	}

	if err := fileoperation.AppendLines(rootList, opt.RootOutputFilename); err != nil {
		return 0, nil, err
	}

	// Stop progress bar
	bar.Stop()

	return found, rootList, nil
}

// updateProgressBar is function called asynchronously to update the progress bar.
func (f *DefaultWildcardFilter) updateProgressBar(bar *progressbar.ProgressBar) {
	current := f.wc.Current()

	bar.SetCurrent(int64(current))
	bar.Set("queries", fmt.Sprintf("%d", f.wc.QueryCount()))
}

// createCacheReader creates a new cache reader. The reader needs to be closed
// by the caller in order to free the file.
func createCacheReader(filename string) (*CacheReader, error) {
	// Open the cache file
	cacheFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	cacheReader := NewCacheReader(cacheFile)

	return cacheReader, nil
}

// createWildcarder creates a new wildcarder.Wildcarder object from the options specified.
func createWildcarder(opt WildcardFilterOptions) *wildcarder.Wildcarder {
	// Convert global QPS to a QPS per resolver
	qps := qpsPerResolver(len(opt.Resolvers), opt.QueriesPerSecond)

	// Create a custom resolver
	resolver := wildcarder.NewClientDNS(opt.Resolvers, 10, qps, 100)

	// Create the wildcarder with the custom resolver
	wc := wildcarder.New(opt.ThreadCount, opt.ResolveTestCount, wildcarder.WithResolver(resolver))

	return wc
}

// prepareCache loads a massdns cache from a reader, saves the valid domains it contains to a file,
// and returns a populated wildcarder.DNSCache object along with the domain file created.
// The caller is responsible to close the domain file returned.
func prepareCache(cacheReader *CacheReader, tempFilename string, batchSize int) (*wildcarder.DNSCache, *os.File, int, error) {
	// Create the temporary file that will hold domains
	domainFile, err := os.Create(tempFilename)
	if err != nil {
		return nil, nil, 0, err
	}

	// Load cache and save found domains to file
	precache := wildcarder.NewDNSCache()
	totalCount, err := cacheReader.Read(domainFile, precache, batchSize)

	// Make sure domain data is written to disk and seek to the beginning of the file
	if err := domainFile.Sync(); err != nil {
		return nil, nil, 0, err
	}

	if _, err := domainFile.Seek(0, 0); err != nil {
		return nil, nil, 0, err
	}

	return precache, domainFile, totalCount, err
}

// qpsPerResolver transforms a global number of queries per second into a number of queries per second per resolver.
func qpsPerResolver(resolverCount, globalQPS int) int {
	if resolverCount == 0 {
		return 0
	}

	qps := globalQPS / resolverCount

	// Set a minimum of 1 query per second
	if qps == 0 {
		qps = 1
	}

	return qps
}
