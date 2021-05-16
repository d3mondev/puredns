package wildcarder

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/d3mondev/puredns/v2/pkg/threadpool"
)

var defaultResolvers []string = []string{
	"8.8.8.8",
	"8.8.4.4",
}

// Wildcarder filters out wildcard subdomains from a list.
type Wildcarder struct {
	resolver    Resolver
	threadCount int

	answerCache *answerCache
	preCache    *DNSCache
	dnsCache    *DNSCache

	tpool      *threadpool.ThreadPool
	tpoolMutex sync.Mutex
	total      int

	randomSubdomains []string
}

// Resolver resolves domain names A and AAAA records and returns the DNS answers found.
type Resolver interface {
	Resolve(domains []string) []DNSAnswer
	QueryCount() int
}

type result struct {
	mu      sync.Mutex
	domains []string
}

// New returns a Wildcarder object used to filter out wildcards.
func New(threadCount int, testCount int, options ...Option) *Wildcarder {
	config := buildConfig(options)

	resolver := config.resolver
	if resolver == nil {
		resolver = NewClientDNS(defaultResolvers, 3, 100, 10)
	}

	precache := config.precache
	if precache == nil {
		precache = NewDNSCache()
	}

	wc := &Wildcarder{
		threadCount: threadCount,
		resolver:    resolver,

		answerCache: newAnswerCache(),
		preCache:    precache,
		dnsCache:    NewDNSCache(),

		randomSubdomains: newRandomSubdomains(testCount),
	}

	return wc
}

// Filter reads subdomains from a reader and returns a list of domains that are not wildcards,
// along with the wildcard subdomain roots found.
func (wc *Wildcarder) Filter(r io.Reader) (domains, roots []string) {
	// Mutex used because a progress bar could be trying to access wc.tpool through wc.Current(),
	// creating a benign race condition that can make tests fail
	wc.tpoolMutex.Lock()
	if wc.tpool != nil {
		panic("concurrent executions of Filter on the same Wildcarder object is not supported")
	}
	wc.tpool = threadpool.NewThreadPool(wc.threadCount, 1000)
	wc.tpoolMutex.Unlock()

	results := &result{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain == "" {
			continue
		}

		ctx := detectionTaskContext{
			results: results,

			resolver:      wc.resolver,
			wildcardCache: wc.answerCache,
			preCache:      wc.preCache,
			dnsCache:      wc.dnsCache,

			randomSubs: wc.randomSubdomains,
			queryCount: len(wc.randomSubdomains),
		}

		task := newDetectionTask(ctx, domain)
		wc.tpool.Execute(task)
	}

	wc.tpool.Wait()

	wc.tpoolMutex.Lock()
	wc.total += wc.tpool.CurrentCount()
	wc.tpool.Close()
	wc.tpool = nil
	wc.tpoolMutex.Unlock()

	domains = results.domains
	roots = gatherRoots(wc.answerCache)

	return domains, roots
}

// QueryCount returns the total number of DNS queries made so far to detect wildcards.
func (wc *Wildcarder) QueryCount() int {
	return wc.resolver.QueryCount()
}

// Current returns the current number of domains that have been processed.
func (wc *Wildcarder) Current() int {
	wc.tpoolMutex.Lock()
	defer wc.tpoolMutex.Unlock()

	if wc.tpool == nil {
		return wc.total
	}

	return wc.total + wc.tpool.CurrentCount()
}

// SetPreCache sets the precache after the Wildcarder object has been created.
func (wc *Wildcarder) SetPreCache(precache *DNSCache) {
	wc.preCache = precache
}

// Option configures a wildcarder.
type Option interface {
	apply(c *config)
}

// WithPreCache returns an option that provides a pre-populated DNS cache used to
// optimize the number of DNS queries made during the wildcard detection phase.
// This DNS cache is not trusted, and the results will be validated as needed using trusted resolvers.
func WithPreCache(cache *DNSCache) Option {
	return precacheOption{precache: cache}
}

type precacheOption struct {
	precache *DNSCache
}

func (o precacheOption) apply(c *config) {
	c.precache = o.precache
}

// WithResolver returns an option that provides a custom resolver to use while performing wildcard detection.
func WithResolver(resolver Resolver) Option {
	return resolverOption{resolver: resolver}
}

type resolverOption struct {
	resolver Resolver
}

func (o resolverOption) apply(c *config) {
	c.resolver = o.resolver
}

type config struct {
	precache *DNSCache
	resolver Resolver
}

func buildConfig(options []Option) config {
	config := config{}

	for _, opt := range options {
		opt.apply(&config)
	}

	return config
}
