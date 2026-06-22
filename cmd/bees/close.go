package main

import (
	"github.com/spf13/cobra"
)

func newCloseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "close <id>",
		Short: "Close an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Close(cmd.Context(), cmd.OutOrStdout(), args[0])
		},
	}
}
