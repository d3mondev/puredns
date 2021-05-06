package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/d3mondev/puredns/v2/internal/usecase/programbanner"
	"github.com/d3mondev/puredns/v2/internal/usecase/resolve"
	"github.com/spf13/cobra"
)

func newCmdBruteforce() *cobra.Command {
	cmdBruteforce := &cobra.Command{
		Use:   "bruteforce <wordlist> domain [flags]",
		Short: "Bruteforce subdomains using a wordlist",
		Long: `Bruteforce takes a file containing words to test as subdomains against the
domain specified. It will invoke massdns using public resolvers for
a quick first pass, then attempt to filter out any wildcard subdomains found.
Finally, it will ensure the results are free of DNS poisoning by resolving
the remaining domains using trusted resolvers.

The <wordlist> argument can be omitted if the wordlist is read from stdin.`,
		Args: cobra.MinimumNArgs(1),
		RunE: runBruteforce,
	}

	cmdBruteforce.Flags().AddFlagSet(resolveFlags)
	cmdBruteforce.Flags().SortFlags = false

	return cmdBruteforce
}

func runBruteforce(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		if !hasStdin() {
			fmt.Println(cmd.UsageString())
			return errors.New("requires a wordlist")
		}
		context.Stdin = os.Stdin
		resolveOptions.Domain = args[0]
	} else {
		resolveOptions.Wordlist = args[0]
		resolveOptions.Domain = args[1]
	}

	resolveOptions.Mode = 1

	if err := resolveOptions.Validate(); err != nil {
		return err
	}

	bannerService := programbanner.NewService(context)
	resolveService := resolve.NewService(context, resolveOptions)

	err := resolveService.Initialize()
	if err != nil {
		return err
	}
	defer resolveService.Close()

	bannerService.PrintWithResolveOptions(resolveOptions)

	return resolveService.Resolve()
}
