package main

import (
	"github.com/spf13/cobra"
)

func newDepCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dep",
		Short: "Manage and graph issue dependencies",
	}

	cmd.AddCommand(newDepAddCmd())
	cmd.AddCommand(newDepRemoveCmd())
	cmd.AddCommand(newDepGraphCmd())

	return cmd
}

func newDepAddCmd() *cobra.Command {
	var blockedID string

	cmd := &cobra.Command{
		Use:   "add <blocker-id> --blocks <blocked-id>",
		Short: "Add a blocking dependency",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.DepAdd(cmd.Context(), cmd.OutOrStdout(), args[0], blockedID)
		},
	}

	cmd.Flags().StringVar(&blockedID, "blocks", "", "ID of the issue being blocked")
	cmd.MarkFlagRequired("blocks")

	return cmd
}

func newDepRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <blocker-id> <blocked-id>",
		Short: "Remove a blocking dependency",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.DepRemove(cmd.Context(), cmd.OutOrStdout(), args[0], args[1])
		},
	}
}

func newDepGraphCmd() *cobra.Command {
	var status string

	cmd := &cobra.Command{
		Use:   "graph [<id>]",
		Short: "Show the dependency graph",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var id *string
			if len(args) == 1 {
				id = &args[0]
			}
			return h.DepGraph(cmd.Context(), cmd.OutOrStdout(), id, status)
		},
	}

	cmd.Flags().StringVar(&status, "status", "", `Filter by status (open, in_progress, closed, all) (default excludes closed)`)

	return cmd
}
