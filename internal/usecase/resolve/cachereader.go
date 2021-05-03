package resolve

import (
	"bufio"
	"io"
	"strings"

	"github.com/d3mondev/puredns/v2/pkg/wildcarder"
	"github.com/d3mondev/resolvermt"
)

// CacheReader reads a DNS cache from a file and can fill a wildcarder.DNSCache object,
// save valid domains to a file, and count valid domains. The number of items processed can be
// limited to a specific number, and subsequent calls to Read will resume without starting over.
type CacheReader struct {
	reader  io.ReadCloser
	scanner *bufio.Scanner
}

// NewCacheReader returns a new CacheReader.
func NewCacheReader(r io.ReadCloser) *CacheReader {
	return &CacheReader{
		reader:  r,
		scanner: bufio.NewScanner(r),
	}
}

// Read reads a massdns cache from a file (created with -o Snl), can save the valid domains to a writer,
// fill a wildcarder.DNSCache object, and return the number of valid domains in the cache.
// Subsequent calls to Read will resume without starting over.
func (c CacheReader) Read(w io.Writer, cache *wildcarder.DNSCache, maxCount int) (count int, err error) {
	type state int
	const (
		stateNewAnswerSection state = iota
		stateSaveAnswer
		stateSkip
	)

	var curDomain string
	var curState state
	var domainSaved bool
	var found int

	for c.scanner.Scan() {
		line := c.scanner.Text()

		// If we receive an empty line, it's the beginning of a new answer
		if line == "" {
			curState = stateNewAnswerSection

			// Break from the loop if we have reached the maximum number of elements to process
			if maxCount > 0 && found == maxCount {
				break
			}

			continue
		}

		switch curState {
		// We're at the beginning of a new answer section, look for the domain name
		case stateNewAnswerSection:
			// Records should be in the form "domain RRTYPE answer"
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				curState = stateSkip
				continue
			}

			domain := strings.TrimSuffix(parts[0], ".")
			if domain == "" {
				curState = stateSkip
				continue
			}

			curDomain = domain
			domainSaved = false
			curState = stateSaveAnswer

			fallthrough

		// Save the answer record found
		case stateSaveAnswer:
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				curState = stateSkip
				continue
			}

			domain := curDomain
			rrtypeStr := parts[1]
			answer := parts[2]

			var rrtype resolvermt.RRtype
			switch rrtypeStr {
			case "A":
				rrtype = resolvermt.TypeA
			case "AAAA":
				rrtype = resolvermt.TypeAAAA
			case "CNAME":
				answer = strings.TrimSuffix(answer, ".")
				rrtype = resolvermt.TypeCNAME
			default:
				continue
			}

			// Save valid domain just once
			if !domainSaved {
				found++
				domainSaved = true

				if w != nil {
					w.Write([]byte(domain + "\n"))
				}
			}

			// Valid record found, add it to the cache
			if cache != nil {
				cacheAnswer := wildcarder.DNSAnswer{
					Type:   rrtype,
					Answer: answer,
				}
				cache.Add(domain, []wildcarder.DNSAnswer{cacheAnswer})
			}

			// If we're just counting valid domains, we can skip the rest of the records
			if cache == nil && w == nil {
				curState = stateSkip
			}

		// Answer was invalid, skip until we receive a new answer section
		case stateSkip:
			continue
		}
	}

	return found, c.scanner.Err()
}

// Close closes the input reader.
func (c CacheReader) Close() error {
	return c.reader.Close()
}
