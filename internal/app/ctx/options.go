package ctx

// GlobalOptions contains the program's global options.
type GlobalOptions struct {
	TrustedResolvers []string
	Quiet            bool
}

// DefaultGlobalOptions creates a new GlobalOptions struct with default values.
func DefaultGlobalOptions() *GlobalOptions {
	return &GlobalOptions{
		TrustedResolvers: []string{
			"8.8.8.8",
			"8.8.4.4",
		},

		Quiet: false,
	}
}

// ResolveOptions contains a resolve command's options.
type ResolveOptions struct {
	BinPath string

	ResolverFile        string
	ResolverTrustedFile string
	NoPublicResolvers   bool

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

	Mode       int
	Domain     string
	Wordlist   string
	DomainFile string
}

// DefaultResolveOptions creates a new ResolveOptions struct with default values.
func DefaultResolveOptions() *ResolveOptions {
	return &ResolveOptions{
		BinPath: "massdns",

		ResolverFile:        "resolvers.txt",
		ResolverTrustedFile: "",
		NoPublicResolvers:   false,

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

		Mode:       0,
		Domain:     "",
		Wordlist:   "",
		DomainFile: "",
	}
}

// Validate validates the options.
func (o *ResolveOptions) Validate() error {
	// Enforce --skip-validation when --no-public is set
	if o.NoPublicResolvers {
		o.SkipValidation = true
	}

	return nil
}
