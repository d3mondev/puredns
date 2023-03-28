package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/d3mondev/puredns/v2/internal/app"
	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/internal/usecase/programbanner"
	"github.com/d3mondev/puredns/v2/internal/usecase/resolve"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	resolveFlags   *pflag.FlagSet
	resolveOptions *ctx.ResolveOptions
)

func newCmdResolve() *cobra.Command {
	resolveOptions = ctx.DefaultResolveOptions()

	cmdResolve := &cobra.Command{
		Use:   "resolve <file> [flags]",
		Short: "Resolve a list of domains",
		Long: `Resolve takes a file containing a list of domains and performs DNS queries
to resolve each domain. It will invoke massdns using public resolvers for
a quick first pass, then attempt to filter out any wildcard subdomains found.
Finally, it will ensure the results are free of DNS poisoning by resolving
the remaining domains using trusted resolvers.

The <file> argument can be omitted if the domains to resolve are read from stdin.`,
		Args: cobra.MinimumNArgs(0),
		RunE: runResolve,
	}

	resolveFlags = pflag.NewFlagSet("resolve", pflag.ExitOnError)
	resolveFlags.StringVarP(&resolveOptions.BinPath, "bin", "b", resolveOptions.BinPath, "path to massdns binary file")
	resolveFlags.IntVarP(&resolveOptions.RateLimit, "rate-limit", "l", resolveOptions.RateLimit, "limit total queries per second for public resolvers (0 = unlimited) (default unlimited)")
	resolveFlags.IntVar(&resolveOptions.RateLimitTrusted, "rate-limit-trusted", resolveOptions.RateLimitTrusted, "limit total queries per second for trusted resolvers (0 = unlimited)")
	resolveFlags.StringVarP(&resolveOptions.ResolverFile, "resolvers", "r", resolveOptions.ResolverFile, "text file containing public resolvers")
	resolveFlags.StringVar(&resolveOptions.ResolverTrustedFile, "resolvers-trusted", resolveOptions.ResolverTrustedFile, "text file containing trusted resolvers")
	resolveFlags.IntVarP(&resolveOptions.WildcardThreads, "threads", "t", resolveOptions.WildcardThreads, "number of threads to use while filtering wildcards")
	resolveFlags.BoolVar(&resolveOptions.TrustedOnly, "trusted-only", resolveOptions.TrustedOnly, "use only trusted resolvers (implies --skip-validation)")
	resolveFlags.IntVarP(&resolveOptions.WildcardTests, "wildcard-tests", "n", resolveOptions.WildcardTests, "number of tests to perform to detect DNS load balancing")
	resolveFlags.IntVar(&resolveOptions.WildcardBatchSize, "wildcard-batch", resolveOptions.WildcardBatchSize, "number of subdomains to test for wildcards in a single batch (0 = unlimited) (default unlimited)")
	resolveFlags.StringVarP(&resolveOptions.WriteDomainsFile, "write", "w", resolveOptions.WriteDomainsFile, "write found domains to a file")
	resolveFlags.StringVar(&resolveOptions.WriteMassdnsFile, "write-massdns", resolveOptions.WriteMassdnsFile, "write massdns database to a file (-o Snl format)")
	resolveFlags.StringVar(&resolveOptions.WriteWildcardsFile, "write-wildcards", resolveOptions.WriteWildcardsFile, "write wildcard subdomain roots to a file")
	resolveFlags.BoolVar(&resolveOptions.SkipSanitize, "skip-sanitize", resolveOptions.SkipSanitize, "do not sanitize the list of domains to test")
	resolveFlags.BoolVar(&resolveOptions.SkipWildcard, "skip-wildcard-filter", resolveOptions.SkipWildcard, "do not perform wildcard detection and filtering")
	resolveFlags.BoolVar(&resolveOptions.SkipValidation, "skip-validation", resolveOptions.SkipValidation, "do not validate results with trusted resolvers")

	must(cobra.MarkFlagFilename(resolveFlags, "bin"))
	must(cobra.MarkFlagFilename(resolveFlags, "resolvers"))
	must(cobra.MarkFlagFilename(resolveFlags, "resolvers-trusted"))
	must(cobra.MarkFlagFilename(resolveFlags, "write"))
	must(cobra.MarkFlagFilename(resolveFlags, "write-massdns"))
	must(cobra.MarkFlagFilename(resolveFlags, "write-wildcards"))

	cmdResolve.Flags().AddFlagSet(resolveFlags)
	cmdResolve.Flags().SortFlags = false

	return cmdResolve
}

func runResolve(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		if !app.HasStdin() {
			fmt.Println(cmd.UsageString())
			return errors.New("requires a list of domains to resolve")
		}
		context.Stdin = os.Stdin
	} else {
		resolveOptions.DomainFile = args[0]
	}

	resolveOptions.Mode = 0

	if err := resolveOptions.Validate(); err != nil {
		return err
	}

	bannerService := programbanner.NewService(context)
	resolveService := resolve.NewService(context, resolveOptions)

	err := resolveService.Initialize()
	if err != nil {
		return err
	}
	defer resolveService.Close(context.Options.Debug)

	bannerService.PrintWithResolveOptions(resolveOptions)

	return resolveService.Resolve()
}
