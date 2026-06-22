package main

import (
	"github.com/spf13/cobra"
)

func newHandoffCmd() *cobra.Command {
	var (
		done      string
		remaining string
		decisions string
		uncertain string
	)

	cmd := &cobra.Command{
		Use:   "handoff <id>",
		Short: "Record a structured handoff for an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Handoff(cmd.Context(), cmd.OutOrStdout(), args[0], done, remaining, decisions, uncertain)
		},
	}

	cmd.Flags().StringVar(&done, "done", "", "What was completed")
	cmd.Flags().StringVar(&remaining, "remaining", "", "What remains to be done")
	cmd.Flags().StringVar(&decisions, "decisions", "", "Decisions made")
	cmd.Flags().StringVar(&uncertain, "uncertain", "", "Open questions or uncertainties")

	return cmd
}
