package wildcarder

import (
	"testing"

	"github.com/d3mondev/resolvermt"
	"github.com/stretchr/testify/assert"
)

func TestGatherRoots(t *testing.T) {
	emptyCache := newAnswerCache()

	rootCache := newAnswerCache()
	rootCache.add("root A", []DNSAnswer{{Type: resolvermt.TypeA}})
	rootCache.add("root B", []DNSAnswer{{Type: resolvermt.TypeA}})
	rootCache.add("root C", []DNSAnswer{{Type: resolvermt.TypeA}})

	tests := []struct {
		name      string
		haveCache *answerCache
		want      []string
	}{
		{name: "empty cache", haveCache: emptyCache, want: []string{}},
		{name: "root detected", haveCache: rootCache, want: []string{"root A", "root B", "root C"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := gatherRoots(test.haveCache)

			assert.ElementsMatch(t, test.want, got)
		})
	}
}
