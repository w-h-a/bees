package main

import (
	"github.com/spf13/cobra"
)

func newReopenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reopen <id>",
		Short: "Reopen a closed issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Reopen(cmd.Context(), cmd.OutOrStdout(), args[0])
		},
	}
}
