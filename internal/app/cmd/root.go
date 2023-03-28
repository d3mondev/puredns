package cmd

import (
	"io/ioutil"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/internal/pkg/console"
	"github.com/spf13/cobra"
)

var (
	context *ctx.Ctx
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func newCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   context.ProgramName,
		Short: context.ProgramTagline,
		Long: context.ProgramName + " " + context.ProgramVersion + `

A subdomain bruteforce tool that wraps around massdns to quickly resolve
a massive number of DNS queries. Using its heuristic algorithm, it can filter out
wildcard subdomains and validate that the results are free of DNS poisoning
by using trusted resolvers.`,
		Example: `  puredns resolve domains.txt
  puredns bruteforce wordlist.txt domain.com --resolvers public.txt
  cat domains.txt | puredns resolve`,
		Version: context.ProgramVersion,
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&context.Options.Quiet, "quiet", "q", context.Options.Quiet, "quiet mode")
	rootCmd.PersistentFlags().BoolVar(&context.Options.Debug, "debug", context.Options.Debug, "keep intermediate files")
	rootCmd.Flags().SortFlags = false

	cmdResolve := newCmdResolve()
	cmdBruteforce := newCmdBruteforce()
	cmdSponsors := newCmdSponsors()
	rootCmd.AddCommand(cmdResolve)
	rootCmd.AddCommand(cmdBruteforce)
	rootCmd.AddCommand(cmdSponsors)

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SilenceErrors = true

	rootCmd.PersistentPreRun = preRun

	return rootCmd
}

func preRun(cmd *cobra.Command, args []string) {
	cmd.SilenceUsage = true

	if context.Options.Quiet {
		console.Output = ioutil.Discard
	}
}

// Execute executes the root command.
func Execute(ctx *ctx.Ctx) error {
	context = ctx
	cmdRoot := newCmdRoot()

	return cmdRoot.Execute()
}
