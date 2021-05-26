package wildcarder

import (
	"strings"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/threadpool"
	"github.com/d3mondev/resolvermt"
	"github.com/stretchr/testify/assert"
)

func overrideWildcardTest(wc *Wildcarder) {
	for i := 0; i < len(wc.randomSubdomains); i++ {
		wc.randomSubdomains[i] = "random"
	}
}

type fakeResolver struct {
	records    map[string][]DNSAnswer
	queryCount int
}

func newFakeResolver() *fakeResolver {
	resolver := fakeResolver{}
	resolver.records = make(map[string][]DNSAnswer)

	return &resolver
}

func (r *fakeResolver) addAnswer(query string, answers []DNSAnswer) {
	if _, ok := r.records[query]; ok {
		r.records[query] = append(r.records[query], answers...)
	} else {
		r.records[query] = answers
	}
}

func (r *fakeResolver) Resolve(queries []string) []DNSAnswer {
	answers := []DNSAnswer{}

	for _, query := range queries {
		if answer, ok := r.records[query]; ok {
			answers = append(answers, answer...)
		}

		r.queryCount++
	}

	return answers
}

func (r *fakeResolver) QueryCount() int {
	return r.queryCount
}

func TestNew(t *testing.T) {
	wc := New(1, 1)
	assert.NotNil(t, wc)
}

func TestFilterEmptyDomain(t *testing.T) {
	resolver := newFakeResolver()

	wc := New(1, 1, WithResolver(resolver))
	overrideWildcardTest(wc)

	domains, roots := wc.Filter(strings.NewReader("     \n\n"))

	assert.ElementsMatch(t, []string{}, domains)
	assert.ElementsMatch(t, []string{}, roots)
}

func TestFilterConcurrent(t *testing.T) {
	wc := New(1, 1)
	wc.tpool = &threadpool.ThreadPool{}

	assert.Panics(t, func() { wc.Filter(strings.NewReader("test.com")) })
}

func TestCurrent_AfterFilter(t *testing.T) {
	resolver := newFakeResolver()
	wc := New(1, 1, WithResolver(resolver))
	wc.Filter(strings.NewReader("test.com"))

	got := wc.Current()

	assert.Equal(t, 1, got)
}

func TestCurrent_WithThreadPool(t *testing.T) {
	resolver := newFakeResolver()
	wc := New(1, 1, WithResolver(resolver))
	wc.tpool = threadpool.NewThreadPool(1, 1)
	wc.total = 100

	got := wc.Current()

	assert.Equal(t, 100, got)
}

func TestSetPreCache(t *testing.T) {
	preCache := NewDNSCache()

	resolver := newFakeResolver()
	wc := New(1, 1, WithResolver(resolver))
	assert.NotSame(t, wc.preCache, preCache)

	wc.SetPreCache(preCache)

	assert.Same(t, wc.preCache, preCache)
}

func TestFilterNotWildcard(t *testing.T) {
	resolver := newFakeResolver()
	resolver.addAnswer("test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("www.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("random.test.com", []DNSAnswer{})

	wc := New(1, 1, WithResolver(resolver))
	overrideWildcardTest(wc)

	domains, roots := wc.Filter(strings.NewReader("test.com"))

	assert.ElementsMatch(t, []string{"test.com"}, domains)
	assert.ElementsMatch(t, []string{}, roots)
}

func TestFilterSimpleWildcard(t *testing.T) {
	resolver := newFakeResolver()
	resolver.addAnswer("test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("www.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("random.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})

	wc := New(1, 1, WithResolver(resolver))
	overrideWildcardTest(wc)

	domains, roots := wc.Filter(strings.NewReader("www.test.com"))

	assert.ElementsMatch(t, []string{}, domains)
	assert.ElementsMatch(t, []string{"test.com"}, roots)
}

func TestFilterSimpleWildcardCNAME(t *testing.T) {
	resolver := newFakeResolver()
	resolver.addAnswer("test.com", []DNSAnswer{{Type: resolvermt.TypeCNAME, Answer: "example.com"}})
	resolver.addAnswer("www.test.com", []DNSAnswer{{Type: resolvermt.TypeCNAME, Answer: "example.com"}})
	resolver.addAnswer("random.test.com", []DNSAnswer{{Type: resolvermt.TypeCNAME, Answer: "example.com"}})

	wc := New(1, 1, WithResolver(resolver))
	overrideWildcardTest(wc)

	domains, roots := wc.Filter(strings.NewReader("www.test.com"))

	assert.ElementsMatch(t, []string{}, domains)
	assert.ElementsMatch(t, []string{"test.com"}, roots)
}

func TestFilterMultilevelWildcard(t *testing.T) {
	resolver := newFakeResolver()
	resolver.addAnswer("test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("www.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("random.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("store.www.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("random.www.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})

	wc := New(1, 1, WithResolver(resolver))
	overrideWildcardTest(wc)

	domains, roots := wc.Filter(strings.NewReader("store.www.api.test.com"))

	assert.ElementsMatch(t, []string{}, domains)
	assert.ElementsMatch(t, []string{"api.test.com"}, roots)
}

func TestFilterMultilevelNotWildcard(t *testing.T) {
	resolver := newFakeResolver()
	resolver.addAnswer("test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("www.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("random.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("store.www.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.1"}})
	resolver.addAnswer("random.www.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})

	wc := New(1, 1, WithResolver(resolver))
	overrideWildcardTest(wc)

	domains, roots := wc.Filter(strings.NewReader("store.www.api.test.com"))

	assert.ElementsMatch(t, []string{"store.www.api.test.com"}, domains)
	assert.ElementsMatch(t, []string{"api.test.com"}, roots)
}

func TestFilterMultilevelMultipleWildcards(t *testing.T) {
	resolver := newFakeResolver()
	resolver.addAnswer("test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("www.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.13"}})
	resolver.addAnswer("random.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.13"}})
	resolver.addAnswer("store.www.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.14"}})
	resolver.addAnswer("random.www.api.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.14"}})

	wc := New(1, 1, WithResolver(resolver))
	overrideWildcardTest(wc)

	domains, roots := wc.Filter(strings.NewReader("store.www.api.test.com"))

	assert.ElementsMatch(t, []string{}, domains)
	assert.ElementsMatch(t, []string{"api.test.com"}, roots)
}

func TestFilterWildcardParentSameWildcard(t *testing.T) {
	resolver := newFakeResolver()
	resolver.addAnswer("test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("random.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("custom.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})
	resolver.addAnswer("store.custom.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("random.custom.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})

	wc := New(1, 1, WithResolver(resolver))
	overrideWildcardTest(wc)

	domains, roots := wc.Filter(strings.NewReader("store.custom.test.com"))

	assert.ElementsMatch(t, []string{}, domains)
	assert.ElementsMatch(t, []string{"custom.test.com"}, roots)
}

func TestFilterSimpleWildcardWithPrecache(t *testing.T) {
	resolver := newFakeResolver()
	resolver.addAnswer("test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("www.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("random.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})

	precache := NewDNSCache()
	precache.Add("www.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})

	wc := New(1, 1, WithResolver(resolver), WithPreCache(precache))
	overrideWildcardTest(wc)

	domains, roots := wc.Filter(strings.NewReader("www.test.com"))
	count := wc.QueryCount()

	assert.ElementsMatch(t, []string{}, domains)
	assert.ElementsMatch(t, []string{"test.com"}, roots)
	assert.Equal(t, len(wc.randomSubdomains)*2, count, "once for the parent, once for the wildcard test")

	domains, roots = wc.Filter(strings.NewReader("www.test.com"))
	count = wc.QueryCount()

	assert.ElementsMatch(t, []string{}, domains)
	assert.ElementsMatch(t, []string{"test.com"}, roots)
	assert.Equal(t, len(wc.randomSubdomains)*2, count, "once for the parent, once for the wildcard test")
}

func TestFilterSimpleWildcardWithBadPrecache(t *testing.T) {
	resolver := newFakeResolver()
	resolver.addAnswer("test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("www.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})
	resolver.addAnswer("random.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.5"}})

	precache := NewDNSCache()
	precache.Add("www.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.1.1"}})

	wc := New(1, 1, WithResolver(resolver), WithPreCache(precache))
	overrideWildcardTest(wc)

	domains, roots := wc.Filter(strings.NewReader("www.test.com"))
	count := wc.QueryCount()

	assert.ElementsMatch(t, []string{}, domains)
	assert.ElementsMatch(t, []string{"test.com"}, roots)
	assert.Equal(t, len(wc.randomSubdomains)*2+1, count, "once for the parent, once for the wildcard test, plus one resolve for final check")

	domains, roots = wc.Filter(strings.NewReader("www.test.com"))
	count = wc.QueryCount()

	assert.ElementsMatch(t, []string{}, domains)
	assert.ElementsMatch(t, []string{"test.com"}, roots)
	assert.Equal(t, len(wc.randomSubdomains)*2+1, count, "should be same as previous test because of cache reuse")
}
