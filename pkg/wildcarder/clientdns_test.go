package wildcarder

import (
	"testing"

	"github.com/d3mondev/resolvermt"
	"github.com/stretchr/testify/assert"
)

type stubClient struct {
	empty   bool
	queries int
}

func (r *stubClient) Resolve(domains []string, rrtype resolvermt.RRtype) []resolvermt.Record {
	if r.empty {
		return []resolvermt.Record{}
	}

	records := []resolvermt.Record{}
	for _, domain := range domains {
		var answer string

		if rrtype == resolvermt.TypeA {
			answer = "127.0.0.1"
		}

		records = append(records, resolvermt.Record{Question: domain, Type: rrtype, Answer: answer})
		r.queries++
	}

	return records
}

func (r *stubClient) QueryCount() int {
	return r.queries
}

func TestResolverDNS(t *testing.T) {
	tests := []struct {
		name        string
		haveRecords bool
		want        []DNSAnswer
	}{
		{name: "empty answer", haveRecords: false, want: []DNSAnswer{}},
		{name: "non-empty answer", haveRecords: true, want: []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}}},
	}

	for _, test := range tests {
		client := &stubClient{empty: !test.haveRecords}
		resolver := ClientDNS{}
		resolver.client = client

		got := resolver.Resolve([]string{"test"})

		assert.ElementsMatch(t, test.want, got)
	}
}

func TestResolverDNSQueryCount(t *testing.T) {
	client := &stubClient{empty: false}

	resolver := ClientDNS{}
	resolver.client = client

	got := resolver.QueryCount()
	assert.Equal(t, 0, got, "initial query count should be 0")

	resolver.Resolve([]string{"test A", "test B"})
	got = resolver.QueryCount()
	assert.Equal(t, 2, got, "query count should increment by 1 for each query")
}
