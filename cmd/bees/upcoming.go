package main

import (
	"github.com/spf13/cobra"
)

func newUpcomingCmd() *cobra.Command {
	var (
		days     int
		assignee string
	)

	cmd := &cobra.Command{
		Use:   "upcoming",
		Short: "Show issues scheduled for the coming days",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Upcoming(cmd.Context(), cmd.OutOrStdout(), days, assignee)
		},
	}

	cmd.Flags().IntVar(&days, "days", 0, "Lookahead window in days (default 15)")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Filter by assignee")

	return cmd
}
