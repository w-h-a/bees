package main

import (
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search issues by title or description",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Search(cmd.Context(), cmd.OutOrStdout(), args[0], limit)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 0, "Maximum number of results (default is 50)")

	return cmd
}
