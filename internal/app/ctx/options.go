package ctx

// GlobalOptions contains the program's global options.
type GlobalOptions struct {
	TrustedResolvers []string
	Quiet            bool
}

// NewGlobalOptions creates a new GlobalOptions struct with default values.
func NewGlobalOptions() *GlobalOptions {
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

// NewResolveOptions creates a new ResolveOptions struct with default values.
func NewResolveOptions() *ResolveOptions {
	return &ResolveOptions{
		BinPath: "massdns",

		ResolverFile:        "resolvers.txt",
		ResolverTrustedFile: "",

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
