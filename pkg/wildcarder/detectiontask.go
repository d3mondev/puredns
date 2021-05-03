package wildcarder

import (
	"strings"
)

type detectionTaskContext struct {
	results *result

	resolver      Resolver
	wildcardCache *answerCache
	preCache      cache
	dnsCache      cache

	randomSubs []string
	queryCount int
}

type detectionTask struct {
	domain string
	ctx    detectionTaskContext
}

type cache interface {
	Add(question string, answers []DNSAnswer)
	Find(question string) []AnswerHash
}

func newDetectionTask(ctx detectionTaskContext, domain string) *detectionTask {
	return &detectionTask{
		domain: domain,
		ctx:    ctx,
	}
}

func (t *detectionTask) Run() {
	// Use the precache data to see if we can early out without any DNS queries
	if t.checkPrecache(t.domain) {
		// Wildcard, early out
		return
	}

	// Test if there is a wildcard at the current level
	root := t.testWildcard(t.domain)
	if root == "" {
		// Found domain, add to results
		t.addDomain(t.domain)
		return
	}

	// Now that we filled wildcard cache, check precache again to optimize the number of queries
	if t.checkPrecache(t.domain) {
		// Wildcard, early out
		return
	}

	// Resolve using a trusted DNS resolver and see if answers match wildcard answers
	if t.checkResolve(t.domain) {
		// Found wildcard after resolving, add to cache and early out
		t.ctx.wildcardCache.addHash(root, t.ctx.preCache.Find(t.domain))
		return
	}

	// Found domain, add to results
	t.addDomain(t.domain)
}

func (t *detectionTask) checkPrecache(domain string) bool {
	answers := t.ctx.preCache.Find(domain)

	return t.domainIsWildcard(domain, answers)
}

// testWildcard checks if there's a wildcard present at the subdomain level, and if so ensures
// the wildcard cache is filled with the proper root.
func (t *detectionTask) testWildcard(domain string) string {
	answers := t.resolveRandomSubdomains(domain)

	if len(answers) == 0 {
		return ""
	}

	root, answers := t.findWildcardRoot(domain, answers)
	t.ctx.wildcardCache.addHash(root, answers)

	return root
}

// findWildcardRoot tests the domain's parent to see if it's a wildcard root. If the parent is a wildcard, it tests the parent's parent.
// It accumulates the wildcard DNS answers found until it finds a wildcard root.
func (t *detectionTask) findWildcardRoot(domain string, answers []AnswerHash) (string, []AnswerHash) {
	parent := getParent(domain)
	if parent == "" {
		return domain, answers
	}

	// If parent is a wildcard, check the parent's parent
	parentAnswers := t.resolveWithCache(parent)
	parentRandomAnswers := t.resolveRandomSubdomains(parent)

	if answerMatch(parentAnswers, parentRandomAnswers) {
		answers = appendUnique(answers, parentAnswers...)
		answers = appendUnique(answers, parentRandomAnswers...)
		return t.findWildcardRoot(parent, answers)
	}

	return parent, answers
}

// answerMatch checks if there is at least one match between two set of answers.
func answerMatch(A []AnswerHash, B []AnswerHash) bool {
	for _, a := range A {
		for _, b := range B {
			if a == b {
				return true
			}
		}
	}

	return false
}

// checkResolve resolves the current domain with trusted resolvers. It then checks to see if the answers
// correspond to a wildcard root that was previously found, otherwise the domain's answers are not wildcard answers.
func (t *detectionTask) checkResolve(domain string) bool {
	answers := t.resolveWithCache(domain)

	return t.domainIsWildcard(domain, answers)
}

// resolveRandomSubdomains tests multiple random subdomains at the current domain level to detect wildcard answers.
// Multiple tests are performed to try to get all the possible answers from any DNS load-balancing occuring.
// The function returns the DNS records that were found.
func (t *detectionTask) resolveRandomSubdomains(subdomain string) []AnswerHash {
	// We can't make test subdomains when we're already at the topmost subdomain
	testSubdomains := t.makeTestSubdomains(subdomain)
	if testSubdomains == nil {
		return nil
	}

	// Early out if we have already resolved the first test domain
	if found := t.ctx.dnsCache.Find(testSubdomains[0]); found != nil {
		return found
	}

	// Resolve the first test domain and check if it returns results
	first := t.ctx.resolver.Resolve(testSubdomains[:1])
	t.ctx.dnsCache.Add(testSubdomains[0], first)

	if len(first) == 0 {
		return nil
	}

	// Resolve the rest of the test domains to populate results with load balancing
	rest := t.ctx.resolver.Resolve(testSubdomains[1:])
	t.ctx.dnsCache.Add(testSubdomains[0], rest)

	// Fetch the results back from the cache because they're deduplicated
	return t.ctx.dnsCache.Find(testSubdomains[0])
}

// makeTestSubdomains makes a list of random subdomains that should be inexistent.
func (t *detectionTask) makeTestSubdomains(domain string) []string {
	parent := getParent(domain)
	if parent == "" {
		return nil
	}

	testQueries := make([]string, len(t.ctx.randomSubs))
	for i := range testQueries {
		testQueries[i] = t.ctx.randomSubs[i] + "." + parent
	}

	return testQueries
}

// getParent returns the parent of a subdomain. Returns an empty string if the parent is a TLD.
func getParent(domain string) string {
	if strings.Count(domain, ".") <= 1 {
		return ""
	}

	parts := strings.SplitN(domain, ".", 2)

	return parts[1]
}

// domainIsWildcard checks the wildcardCache for its answers and a matching root domain.
func (t *detectionTask) domainIsWildcard(domain string, answers []AnswerHash) bool {
	roots := t.ctx.wildcardCache.findHash(answers)
	for _, root := range roots {
		if strings.HasSuffix(domain, root) {
			return true
		}
	}

	return false
}

// resolveWithCache checks the DNS cache for an answer, otherwise it performs a number of DNS queries
// to get all the possible answers and add them to the cache.
func (t *detectionTask) resolveWithCache(domain string) []AnswerHash {
	if answers := t.ctx.dnsCache.Find(domain); answers != nil {
		return answers
	}

	first := t.ctx.resolver.Resolve([]string{domain})
	t.ctx.dnsCache.Add(domain, first)

	if len(first) == 0 {
		return nil
	}

	for i := 1; i < t.ctx.queryCount; i++ {
		answers := t.ctx.resolver.Resolve([]string{domain})
		t.ctx.dnsCache.Add(domain, answers)
	}

	return t.ctx.dnsCache.Find(domain)
}

// addDomain adds a domain to the results.
func (t *detectionTask) addDomain(domain string) {
	t.ctx.results.mu.Lock()
	defer t.ctx.results.mu.Unlock()

	t.ctx.results.domains = append(t.ctx.results.domains, domain)
}

// appendUnique append unique answers to the list. Not optimized for complexity due to expected small number of elements.
func appendUnique(list []AnswerHash, elem ...AnswerHash) []AnswerHash {
	for _, e := range elem {
		var found bool

		for _, l := range list {
			if e == l {
				found = true
				break
			}
		}

		if !found {
			list = append(list, e)
		}
	}

	return list
}
