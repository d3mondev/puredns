package cmd

import (
	"os"

	"github.com/d3mondev/puredns/v2/internal/app"
	"github.com/d3mondev/puredns/v2/internal/usecase/programbanner"
	"github.com/d3mondev/puredns/v2/internal/usecase/resolve"
	"github.com/spf13/cobra"
)

func newCmdBruteforce() *cobra.Command {
	cmdBruteforce := &cobra.Command{
		Use:   "bruteforce <wordlist> domain [flags]\n  puredns bruteforce <wordlist> -d domains.txt [flags]",
		Short: "Bruteforce subdomains using a wordlist",
		Long: `Bruteforce takes a file containing words to test as subdomains against the
domain specified. It will invoke massdns using public resolvers for
a quick first pass, then attempt to filter out any wildcard subdomains found.
Finally, it will ensure the results are free of DNS poisoning by resolving
the remaining domains using trusted resolvers.

The <wordlist> argument can be omitted if the wordlist is read from stdin.`,
		RunE: runBruteforce,
	}

	cmdBruteforce.Flags().StringVarP(&resolveOptions.DomainFile, "domains", "d", resolveOptions.DomainFile, "text file containing domains to bruteforce")

	cmdBruteforce.Flags().AddFlagSet(resolveFlags)
	cmdBruteforce.Flags().SortFlags = false

	return cmdBruteforce
}

func runBruteforce(cmd *cobra.Command, args []string) error {
	parseBruteforceArgs(args)

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

func parseBruteforceArgs(args []string) error {
	if app.HasStdin() {
		context.Stdin = os.Stdin

		if len(args) >= 1 {
			if resolveOptions.DomainFile == "" {
				resolveOptions.Domain = args[0]
			}
		}
	} else {
		if len(args) == 1 {
			resolveOptions.Wordlist = args[0]
		} else if len(args) >= 2 {
			resolveOptions.Wordlist = args[0]
			resolveOptions.Domain = args[1]
		}
	}

	resolveOptions.Mode = 1

	return nil
}
