package ctx

import (
	"errors"
	"os/user"
	"path/filepath"

	"github.com/d3mondev/puredns/v2/internal/app"
	"github.com/d3mondev/puredns/v2/pkg/fileoperation"
)

// ResolveMode is the resolve mode.
type ResolveMode int

const (
	// Resolve resolves domains.
	Resolve ResolveMode = iota

	// Bruteforce bruteforces subdomains.
	Bruteforce
)

var (
	// ErrNoDomain no domain specified.
	ErrNoDomain error = errors.New("no domain specified")

	// ErrNoWordlist no wordlist specified.
	ErrNoWordlist error = errors.New("no wordlist specified")
)

// GlobalOptions contains the program's global options.
type GlobalOptions struct {
	TrustedResolvers []string

	Quiet bool
	Debug bool
}

// DefaultGlobalOptions creates a new GlobalOptions struct with default values.
func DefaultGlobalOptions() *GlobalOptions {
	return &GlobalOptions{
		TrustedResolvers: []string{
			"8.8.8.8",
			"8.8.4.4",
		},

		Quiet: false,
		Debug: false,
	}
}

// ResolveOptions contains a resolve command's options.
type ResolveOptions struct {
	BinPath string

	ResolverFile        string
	ResolverTrustedFile string
	TrustedOnly         bool

	RateLimit        int
	RateLimitTrusted int

	WildcardThreads   int
	WildcardTests     int
	WildcardBatchSize int

	SkipSanitize   bool
	SkipWildcard   bool
	SkipValidation bool

	WriteDomainsFile   string
	WriteMassdnsFile   string
	WriteWildcardsFile string

	Mode       ResolveMode
	Domain     string
	Wordlist   string
	DomainFile string
}

// DefaultResolveOptions creates a new ResolveOptions struct with default values.
func DefaultResolveOptions() *ResolveOptions {
	resolversPath := "resolvers.txt"
	trustedResolversPath := ""

	if !fileoperation.FileExists(resolversPath) {
		usr, err := user.Current()
		if err == nil {
			resolversPath = filepath.Join(usr.HomeDir, ".config", "puredns", "resolvers.txt")
			trustedResolversPath = filepath.Join(usr.HomeDir, ".config", "puredns", "resolvers-trusted.txt")

			if !fileoperation.FileExists(trustedResolversPath) {
				trustedResolversPath = ""
			}
		}
	}

	return &ResolveOptions{
		BinPath: "massdns",

		ResolverFile:        resolversPath,
		ResolverTrustedFile: trustedResolversPath,
		TrustedOnly:         false,

		RateLimit:        0,
		RateLimitTrusted: 500,

		WildcardThreads:   100,
		WildcardTests:     3,
		WildcardBatchSize: 0,

		SkipSanitize:   false,
		SkipWildcard:   false,
		SkipValidation: false,

		WriteDomainsFile:   "",
		WriteMassdnsFile:   "",
		WriteWildcardsFile: "",

		Mode:       Resolve,
		Domain:     "",
		Wordlist:   "",
		DomainFile: "",
	}
}

// Validate validates the options.
func (o *ResolveOptions) Validate() error {
	// Enforce --skip-validation when --trusted-only is set
	if o.TrustedOnly {
		o.SkipValidation = true
	}

	// Validate that a wordlist and a domain are present in bruteforce mode
	if o.Mode == Bruteforce {
		if o.Domain == "" && o.DomainFile == "" {
			return ErrNoDomain
		}

		if o.Wordlist == "" && !app.HasStdin() {
			return ErrNoWordlist
		}
	}

	return nil
}
