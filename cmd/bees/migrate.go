package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/w-h-a/bees/internal/util/home"
)

func newMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate <repo>...",
		Short: "Preview migrating per-repo v1 stores into the device-global store (dry run)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			userHome, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to resolve user home: %w", err)
			}

			resolved, err := home.Resolve(os.Getenv("BEES_HOME"), userHome)
			if err != nil {
				return err
			}

			return h.Migrate(cmd.Context(), cmd.OutOrStdout(), resolved.DB, args)
		},
	}
}
