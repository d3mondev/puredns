package resolve

import (
	"io"
	"strings"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/d3mondev/puredns/v2/pkg/wildcarder"
	"github.com/d3mondev/resolvermt"
	"github.com/stretchr/testify/assert"
)

func TestCacheReaderRead(t *testing.T) {
	type cacheEntry struct {
		question string
		answers  []wildcarder.AnswerHash
	}

	tests := []struct {
		name       string
		haveData   string
		wantCache  []cacheEntry
		wantDomain []string
		wantErr    error
	}{
		{
			name:     "single record",
			haveData: `example.com. A 127.0.0.1`,
			wantCache: []cacheEntry{
				{
					question: "example.com",
					answers: []wildcarder.AnswerHash{
						wildcarder.HashAnswer(wildcarder.DNSAnswer{Type: resolvermt.TypeA, Answer: "127.0.0.1"}),
					},
				},
			},
			wantDomain: []string{
				"example.com",
			},
		},
		{
			name: "multiple record",
			haveData: `www.example.com. CNAME example.com.
example.com. A 127.0.0.1
example.com. AAAA ::1`,
			wantCache: []cacheEntry{
				{
					question: "www.example.com",
					answers: []wildcarder.AnswerHash{
						wildcarder.HashAnswer(wildcarder.DNSAnswer{Type: resolvermt.TypeCNAME, Answer: "example.com"}),
						wildcarder.HashAnswer(wildcarder.DNSAnswer{Type: resolvermt.TypeA, Answer: "127.0.0.1"}),
						wildcarder.HashAnswer(wildcarder.DNSAnswer{Type: resolvermt.TypeAAAA, Answer: "::1"}),
					},
				},
			},
			wantDomain: []string{
				"www.example.com",
			},
		},
		{
			name:       "invalid record type",
			haveData:   `example.com. NS ns.example.com.`,
			wantCache:  []cacheEntry{},
			wantDomain: []string{},
		},
		{
			name: "save domain after valid record is found",
			haveData: `example.com. NS ns.example.com.
example.com. AAAA ::1`,
			wantCache: []cacheEntry{
				{
					question: "example.com",
					answers: []wildcarder.AnswerHash{
						wildcarder.HashAnswer(wildcarder.DNSAnswer{Type: resolvermt.TypeAAAA, Answer: "::1"}),
					},
				},
			},
			wantDomain: []string{
				"example.com",
			},
		},
		{
			name: "multiple answer sections",
			haveData: `
example.com. A 127.0.0.1

www.test.com. CNAME test.com.
test.com. A 127.0.0.1
test.com. AAAA ::1
`,
			wantCache: []cacheEntry{
				{
					question: "example.com",
					answers: []wildcarder.AnswerHash{
						wildcarder.HashAnswer(wildcarder.DNSAnswer{Type: resolvermt.TypeA, Answer: "127.0.0.1"}),
					},
				},
				{
					question: "www.test.com",
					answers: []wildcarder.AnswerHash{
						wildcarder.HashAnswer(wildcarder.DNSAnswer{Type: resolvermt.TypeCNAME, Answer: "test.com"}),
						wildcarder.HashAnswer(wildcarder.DNSAnswer{Type: resolvermt.TypeA, Answer: "127.0.0.1"}),
						wildcarder.HashAnswer(wildcarder.DNSAnswer{Type: resolvermt.TypeAAAA, Answer: "::1"}),
					},
				},
			},
			wantDomain: []string{
				"example.com",
				"www.test.com",
			},
		},
		{
			name: "skip if domain name can't be parsed",
			haveData: `garbage
example.com. A 127.0.0.1`,
			wantCache:  []cacheEntry{},
			wantDomain: []string{},
		},
		{
			name: "skip answer section containing bad data",
			haveData: `example.com. A 127.0.0.1
garbage`,
			wantCache: []cacheEntry{
				{
					question: "example.com",
					answers: []wildcarder.AnswerHash{
						wildcarder.HashAnswer(wildcarder.DNSAnswer{Type: resolvermt.TypeA, Answer: "127.0.0.1"}),
					},
				},
			},
			wantDomain: []string{"example.com"},
		},
		{
			name:       "empty domain",
			haveData:   `. A 127.0.0.1`,
			wantCache:  []cacheEntry{},
			wantDomain: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			domainFile := filetest.CreateFile(t, "")
			cache := wildcarder.NewDNSCache()

			loader := NewCacheReader(io.NopCloser(strings.NewReader(test.haveData)))
			count, err := loader.Read(domainFile, cache, 0)
			assert.ErrorIs(t, err, test.wantErr)
			assert.Equal(t, len(test.wantDomain), count)

			gotDomain := filetest.ReadFile(t, domainFile.Name())

			assert.Equal(t, test.wantDomain, gotDomain)

			for _, cacheTest := range test.wantCache {
				got := cache.Find(cacheTest.question)
				assert.ElementsMatch(t, cacheTest.answers, got)
			}
		})
	}
}

func TestCacheReaderRead_WithMax(t *testing.T) {
	domainFile := filetest.CreateFile(t, "")
	data := `
example.com. A 127.0.0.1

example.net. AAAA ::1

example.org. CNAME example.net.
example.net. AAAA ::1`

	r := io.NopCloser(strings.NewReader(data))
	loader := NewCacheReader(r)

	_, err := loader.Read(domainFile, nil, 2)
	gotDomain := filetest.ReadFile(t, domainFile.Name())

	assert.Nil(t, err)
	assert.Equal(t, []string{"example.com", "example.net"}, gotDomain)

	_, err = loader.Read(domainFile, nil, 2)
	gotDomain = filetest.ReadFile(t, domainFile.Name())

	assert.Nil(t, err)
	assert.Equal(t, []string{"example.com", "example.net", "example.org"}, gotDomain)
}

func TestCacheReaderRead_CountOnly(t *testing.T) {
	data := `
example.com. A 127.0.0.1

example.net. AAAA ::1

example.org. CNAME example.net.
example.net. AAAA ::1
`

	r := io.NopCloser(strings.NewReader(data))
	loader := NewCacheReader(r)

	count, _ := loader.Read(nil, nil, 0)
	assert.Equal(t, 3, count)
}
