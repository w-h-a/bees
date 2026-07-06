package main

import (
	"github.com/spf13/cobra"
)

func newContextCmd() *cobra.Command {
	var (
		prefix string
	)

	cmd := &cobra.Command{
		Use:   "context",
		Short: "Show the current bees context",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Context(cmd.Context(), cmd.OutOrStdout(), prefix)
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", `Filter by prefix`)

	return cmd
}
