package main

import (
	"github.com/spf13/cobra"
)

func newMigrateCmd() *cobra.Command {
	var commit bool

	cmd := &cobra.Command{
		Use:   "migrate <repo>...",
		Short: "Migrate per-repo v1 stores into the device-global store (dry run by default)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if commit {
				return h.MigrateCommit(cmd.Context(), cmd.OutOrStdout(), args)
			}

			return h.Migrate(cmd.Context(), cmd.OutOrStdout(), migrateTargetDB, args)
		},
	}

	cmd.Flags().BoolVar(&commit, "commit", false, "Perform the migration (default is a dry run)")

	return cmd
}
