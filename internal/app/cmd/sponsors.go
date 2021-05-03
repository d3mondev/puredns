package cmd

import (
	"github.com/d3mondev/puredns/v2/internal/app"
	"github.com/d3mondev/puredns/v2/internal/usecase/sponsors"
	"github.com/spf13/cobra"
)

func newCmdSponsors() *cobra.Command {
	cmdSponsors := &cobra.Command{
		Use:   "sponsors",
		Short: "Show the active sponsors <3",
		Long: `Show the very kind-hearted people who support my work as sponsors.

This software is made by me, @d3mondev. I'm on a mission to make free and open-souce
software for the bug bounty community and infosec professionals.

As you know, free doesn't help pay the bills. If my work is earning you money,
consider becoming a sponsor! It would mean A WHOLE LOT as it would allow me to continue
working for free for the community: https://github.com/sponsors/d3mondev`,
		RunE: runSponsors,
	}

	return cmdSponsors
}

func runSponsors(cmd *cobra.Command, args []string) error {
	service := sponsors.NewService()

	return service.Show(app.AppSponsorsURL)
}
