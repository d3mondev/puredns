package wildcarder

import (
	"testing"

	"github.com/d3mondev/resolvermt"
	"github.com/stretchr/testify/assert"
)

type cacheData struct {
	root    string
	answers []DNSAnswer
}

func TestAnswerCache(t *testing.T) {
	singleAnswer := []DNSAnswer{
		{
			Type:   resolvermt.TypeA,
			Answer: "127.0.0.1",
		},
	}

	singleCNAME := []DNSAnswer{
		{
			Type:   resolvermt.TypeCNAME,
			Answer: "127.0.0.1",
		},
	}

	multipleAnswers := []DNSAnswer{
		{
			Type:   resolvermt.TypeA,
			Answer: "127.0.0.1",
		},
		{
			Type:   resolvermt.TypeAAAA,
			Answer: "::1",
		},
		{
			Type:   resolvermt.TypeCNAME,
			Answer: "cname",
		},
	}

	tests := []struct {
		name          string
		haveCacheData []cacheData
		haveSearch    []DNSAnswer
		want          []string
	}{
		{name: "empty cache", haveCacheData: nil, haveSearch: singleAnswer, want: []string{}},
		{name: "empty search", haveCacheData: []cacheData{{root: "root", answers: singleAnswer}}, haveSearch: nil, want: []string{}},
		{name: "find single record", haveCacheData: []cacheData{{root: "root", answers: multipleAnswers}}, haveSearch: singleAnswer, want: []string{"root"}},
		{name: "same root", haveCacheData: []cacheData{{root: "root", answers: singleAnswer}, {root: "root", answers: singleCNAME}}, haveSearch: singleAnswer, want: []string{"root"}},
		{name: "duplicate answer", haveCacheData: []cacheData{{root: "root", answers: singleAnswer}, {root: "root", answers: singleAnswer}}, haveSearch: singleAnswer, want: []string{"root"}},
		{name: "multiple roots",
			haveCacheData: []cacheData{
				{
					root:    "root A",
					answers: singleAnswer,
				},
				{
					root:    "root B",
					answers: singleAnswer,
				},
			},
			haveSearch: multipleAnswers,
			want:       []string{"root A", "root B"},
		},
		{name: "different types", haveCacheData: []cacheData{{root: "root", answers: multipleAnswers}}, haveSearch: singleCNAME, want: []string{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := newAnswerCache()

			for _, data := range test.haveCacheData {
				cache.add(data.root, data.answers)
			}

			got := cache.find(test.haveSearch)
			assert.ElementsMatch(t, test.want, got)
		})
	}
}

func TestAnswerCacheCount(t *testing.T) {
	cache := newAnswerCache()
	assert.Equal(t, 0, cache.count(), "empty cache count is 0")

	cache.add("root", []DNSAnswer{})
	assert.Equal(t, 0, cache.count(), "empty record cache count is 0")

	cache.add("root", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})
	assert.Equal(t, 1, cache.count(), "add single answer cache count is 1")

	cache.add("root", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "192.168.0.1"}})
	assert.Equal(t, 2, cache.count(), "add new answer cache count is 2")
}
