package resolve

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/d3mondev/puredns/v2/pkg/procreader"
)

// DomainReader implements an io.Reader interface that generates subdomains to resolve.
// It reads data line by line from a source scanner. This data is either words that will be
// prefixed to a domain to create subdomains, or a straight list of subdomains to resolve.
// The DomainReader will also discard any generated domains that do not pass the specified
// domain sanitizer filter if present.
type DomainReader struct {
	source          io.ReadCloser
	sourceScanner   *bufio.Scanner
	subdomainReader *procreader.ProcReader

	domains   []string
	sanitizer DomainSanitizer
}

var _ io.Reader = (*DomainReader)(nil)

// DomainSanitizer is a function that sanitizes a domain, typically removing invalid characters.
// If the domain cannot be sanitized or is invalid, an empty string is expected.
type DomainSanitizer func(domain string) string

// NewDomainReader creates a new DomainReader. If domains is not empty, the source
// reader is expected to contain words that will be prefixed to the domains to create subdomains.
func NewDomainReader(source io.ReadCloser, domains []string, sanitizer DomainSanitizer) *DomainReader {
	domainReader := &DomainReader{
		source:        source,
		sourceScanner: bufio.NewScanner(source),
		domains:       domains,
		sanitizer:     sanitizer,
	}

	domainReader.subdomainReader = procreader.New(domainReader.nextSubdomains)

	return domainReader
}

// Read creates and returns subdomains in the buffer specified.
func (r *DomainReader) Read(p []byte) (int, error) {
	return r.subdomainReader.Read(p)
}

// nextSubdomain is a callback used to generate the next subdomains.
func (r *DomainReader) nextSubdomains(size int) ([]byte, error) {
	if !r.sourceScanner.Scan() {
		// Make sure the close the source, discarding the error
		// as we want the error from the scanner
		r.source.Close()

		// Return the error from the scanner
		if err := r.sourceScanner.Err(); err != nil {
			return nil, err
		}

		// Return EOF
		return nil, io.EOF
	}

	var output bytes.Buffer
	word := r.sourceScanner.Text()

	if len(r.domains) == 0 {
		// Single domain was read from reader
		domain := word
		domain = r.processDomain(domain)
		output.WriteString(domain)
	} else {
		// Generate a subdomain from the word and the list of domains
		for _, domain := range r.domains {
			if strings.ContainsRune(domain, '*') {
				domain = strings.ReplaceAll(domain, "*", word)
			} else {
				domain = fmt.Sprintf("%s.%s", word, domain)
			}

			domain = r.processDomain(domain)
			output.WriteString(domain)
		}
	}

	return output.Bytes(), nil
}

// processDomain processes the domain data
func (r *DomainReader) processDomain(domain string) string {
	// Sanitize the domain
	if r.sanitizer != nil {
		domain = r.sanitizer(domain)
	}

	// Append newline even if we have empty domain for accurate progress bar
	domain = domain + "\n"

	return domain
}
