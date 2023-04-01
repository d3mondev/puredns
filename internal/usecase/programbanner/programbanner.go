package programbanner

import (
	"fmt"
	"strings"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/internal/pkg/console"
)

// Service prints the program banner and version number.
type Service struct {
	ctx *ctx.Ctx
}

// NewService returns a new Service object.
func NewService(ctx *ctx.Ctx) Service {
	return Service{
		ctx: ctx,
	}
}

// Print prints the program logo along with its name, tagline and version information.
func (s Service) Print() {
	version := s.ctx.ProgramVersion
	if s.ctx.GitBranch != "" {
		version = fmt.Sprintf("%s-%s", s.ctx.GitBranch, s.ctx.GitRevision)
	}

	padding := strings.Repeat(" ", 34-len(version)-len(s.ctx.ProgramName))

	console.Printf(console.ColorBrightBlue)
	console.Printf("                          _           \n")
	console.Printf("                         | |          \n")
	console.Printf(" _ __  _   _ _ __ ___  __| |_ __  ___ \n")
	console.Printf("| '_ \\| | | | '__/ _ \\/ _` | '_ \\/ __|\n")
	console.Printf("| |_) | |_| | | |  __/ (_| | | | \\__ \\\n")
	console.Printf("| .__/ \\__,_|_|  \\___|\\__,_|_| |_|___/\n")
	console.Printf("| |                                   \n")
	console.Printf("|_|%s%s%s %s%s\n", padding, console.ColorBrightCyan, s.ctx.ProgramName, console.ColorBrightBlue, version)
	console.Printf("\n")
	console.Printf("%sFast and accurate DNS resolving and bruteforcing\n", console.ColorBrightWhite)
	console.Printf("\n")
	console.Printf("%sCrafted with %s<3%s by @d3mondev\n", console.ColorBrightWhite, console.ColorBrightRed, console.ColorBrightWhite)
	console.Printf("https://github.com/sponsors/d3mondev\n")
	console.Printf(console.ColorReset + "\n")
}

// PrintWithResolveOptions prints the program's logo, along with the options selected
// for the resolve command.
func (s Service) PrintWithResolveOptions(opts *ctx.ResolveOptions) {
	s.Print()
	console.Printf(console.ColorBrightWhite + "------------------------------------------------------------\n" + console.ColorReset)

	defaultOptions := ctx.DefaultResolveOptions()

	var file string
	if s.ctx.Stdin != nil {
		file = "stdin"
	} else {
		if opts.Mode == 1 {
			file = opts.Wordlist
		} else {
			file = opts.DomainFile
		}
	}

	colorOptionLabel := console.ColorBrightWhite
	colorOptionSkipLabel := console.ColorBrightYellow
	colorOptionValue := console.ColorWhite
	colorOptionTick := console.ColorBrightBlue
	colorOptionTickWrite := console.ColorBrightGreen

	tickSymbol := fmt.Sprintf("%s[%s+%s]", colorOptionLabel, colorOptionTick, colorOptionLabel)
	tickSymbolWrite := fmt.Sprintf("%s[%s+%s]", colorOptionLabel, colorOptionTickWrite, colorOptionLabel)

	if opts.Mode == 1 {
		console.Printf("%s Mode                 :%s bruteforce\n", tickSymbol, colorOptionValue)

		if opts.DomainFile != "" {
			console.Printf("%s Domains              :%s %s\n", tickSymbol, colorOptionValue, opts.DomainFile)
		} else {
			console.Printf("%s Domain               :%s %s\n", tickSymbol, colorOptionValue, opts.Domain)
		}

		console.Printf("%s Wordlist             :%s %s\n", tickSymbol, colorOptionValue, file)
	} else {
		console.Printf("%s Mode                 :%s resolve\n", tickSymbol, colorOptionValue)
		console.Printf("%s File                 :%s %s\n", tickSymbol, colorOptionValue, file)
	}

	if opts.TrustedOnly {
		console.Printf("%s Trusted Only         :%s true\n", tickSymbol, colorOptionValue)
	}

	if !opts.TrustedOnly {
		console.Printf("%s Resolvers            :%s %s\n", tickSymbol, colorOptionValue, opts.ResolverFile)
	}

	if opts.ResolverTrustedFile != "" {
		console.Printf("%s Trusted Resolvers    :%s %s\n", tickSymbol, colorOptionValue, opts.ResolverTrustedFile)
	}

	if !opts.TrustedOnly {
		rate := "unlimited"
		if opts.RateLimit != 0 {
			rate = fmt.Sprintf("%d qps", opts.RateLimit)
		}
		console.Printf("%s Rate Limit           :%s %s\n", tickSymbol, colorOptionValue, rate)
	}

	console.Printf("%s Rate Limit (Trusted) :%s %d qps\n", tickSymbol, colorOptionValue, opts.RateLimitTrusted)
	console.Printf("%s Wildcard Threads     :%s %d\n", tickSymbol, colorOptionValue, opts.WildcardThreads)
	console.Printf("%s Wildcard Tests       :%s %d\n", tickSymbol, colorOptionValue, opts.WildcardTests)

	if opts.WildcardBatchSize != defaultOptions.WildcardBatchSize {
		console.Printf("%s Wildcard Batch Size  :%s %d\n", tickSymbol, colorOptionValue, opts.WildcardBatchSize)
	}

	if opts.WriteDomainsFile != "" {
		console.Printf("%s Write Domains        :%s %s\n", tickSymbolWrite, colorOptionValue, opts.WriteDomainsFile)
	}

	if opts.WriteMassdnsFile != "" {
		console.Printf("%s Write Massdns        :%s %s\n", tickSymbolWrite, colorOptionValue, opts.WriteMassdnsFile)
	}

	if opts.WriteWildcardsFile != "" {
		console.Printf("%s Write Wildcards      :%s %s\n", tickSymbolWrite, colorOptionValue, opts.WriteWildcardsFile)
	}

	if opts.SkipSanitize {
		console.Printf("%s[+] Skip Sanitize\n", colorOptionSkipLabel)
	}

	if opts.SkipWildcard {
		console.Printf("%s[+] Skip Wildcard Detection\n", colorOptionSkipLabel)
	}

	if !opts.TrustedOnly {
		if opts.SkipValidation {
			console.Printf("%s[+] Skip Validation\n", colorOptionSkipLabel)
		}
	}

	console.Printf(console.ColorBrightWhite + "------------------------------------------------------------\n" + console.ColorReset)
	console.Printf("\n")
}
