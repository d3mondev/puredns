package wildcarder

import "github.com/d3mondev/resolvermt"

// DNSAnswer represents a DNS answer without the question.
type DNSAnswer struct {
	Type   resolvermt.RRtype
	Answer string
}

// ClientDNS is a DNS client that implements the Resolver interface.
type ClientDNS struct {
	client resolver
}

type resolver interface {
	Resolve(domains []string, rrtype resolvermt.RRtype) []resolvermt.Record
	QueryCount() int
}

// NewClientDNS creates a new ResolverDNS object to use with a Wildcarder object.
func NewClientDNS(resolvers []string, retryCount int, qps int, concurrency int) *ClientDNS {
	return &ClientDNS{
		client: resolvermt.New(resolvers, retryCount, qps, concurrency),
	}
}

// Resolve resolves A records from a list of domain names and returns the answers.
func (r *ClientDNS) Resolve(domains []string) []DNSAnswer {
	records := r.client.Resolve(domains, resolvermt.TypeA)

	// Removed AAAA records as those are not being handled by massdns right now
	//	records = append(records, r.client.Resolve(domains, resolvermt.TypeAAAA)...)

	answers := []DNSAnswer{}
	for _, record := range records {
		answers = append(answers, DNSAnswer{Type: record.Type, Answer: record.Answer})
	}

	return answers
}

// QueryCount returns the number of DNS queries really performed.
func (r *ClientDNS) QueryCount() int {
	return r.client.QueryCount()
}
