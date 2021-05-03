package massdns

import (
	"fmt"
	"os"
	"strings"
)

type state int

const (
	stateNewAnswerSection state = iota
	stateSaveAnswer
	stateSkip
)

// DefaultWriteCallback is a callback that can save massdns results to files.
// It can save the valid domains found and the massdns results that gave valid domains.
type DefaultWriteCallback struct {
	massdnsFile *os.File
	domainFile  *os.File

	curState    state
	curDomain   string
	domainSaved bool

	found int
}

var _ Callback = (*DefaultWriteCallback)(nil)

// NewDefaultWriteCallback creates a new DefaultWriteCallback.
// The file names can be empty to disable saving to a file.
func NewDefaultWriteCallback(massdnsFilename string, domainFilename string) (*DefaultWriteCallback, error) {
	cb := &DefaultWriteCallback{}

	// Create a writer that writes massdns answers
	if massdnsFilename != "" {
		file, err := os.Create(massdnsFilename)
		if err != nil {
			return nil, err
		}

		cb.massdnsFile = file
	}

	// Create a writer that writes valid domains found
	if domainFilename != "" {
		file, err := os.Create(domainFilename)
		if err != nil {
			return nil, err
		}

		cb.domainFile = file
	}

	return cb, nil
}

// Callback reads a line from the massdns stdout handler, parses the output and
// saves the relevant data.
func (c *DefaultWriteCallback) Callback(line string) error {
	// Don't parse JSON if we're not saving anything
	if c.domainFile == nil && c.massdnsFile == nil {
		return nil
	}

	// If we receive an empty line, it's the start of a new answer
	if line == "" {
		c.curState = stateNewAnswerSection
		return nil
	}

	switch c.curState {
	// We're at the beginning of a new answer section, look for the domain name
	case stateNewAnswerSection:
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			c.curState = stateSkip
			return nil
		}

		domain := strings.TrimSuffix(parts[0], ".")
		if domain == "" {
			c.curState = stateSkip
			return nil
		}

		c.curDomain = domain
		c.curState = stateSaveAnswer
		c.domainSaved = false
		fallthrough

	// Save the answer record found
	case stateSaveAnswer:
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			c.curState = stateSkip
			return nil
		}

		domain := c.curDomain
		rrType := parts[1]
		answer := strings.TrimSuffix(parts[2], ".")

		// Only look for A, AAAA, and CNAME records
		if rrType != "A" && rrType != "AAAA" && rrType != "CNAME" {
			return nil
		}

		// If we haven't saved the domain yet, save it
		if !c.domainSaved {
			c.saveDomain(c.curDomain)
			c.domainSaved = true
			c.found++
		}

		// Valid record found, save it
		return c.saveLine(fmt.Sprintf("%s %s %s", domain, rrType, answer))

	// Answer was invalid, skip until we receive a new answer section
	case stateSkip:
		return nil
	}

	return nil
}

// saveLine saves a line to the massdns file.
func (c *DefaultWriteCallback) saveLine(line string) error {
	if c.massdnsFile != nil {
		_, err := c.massdnsFile.WriteString(line + "\n")
		return err
	}

	return nil
}

// saveDomain saves a domain to the domain file.
func (c *DefaultWriteCallback) saveDomain(domain string) error {
	if c.domainFile != nil {
		_, err := c.domainFile.WriteString(domain + "\n")
		return err
	}

	return nil
}

// Close closes the writers.
func (c *DefaultWriteCallback) Close() {
	if c.massdnsFile != nil {
		c.massdnsFile.Sync()
		c.massdnsFile.Close()
	}

	if c.domainFile != nil {
		c.domainFile.Sync()
		c.domainFile.Close()
	}
}
