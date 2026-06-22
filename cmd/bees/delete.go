package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/w-h-a/bees/internal/domain"
	"github.com/w-h-a/bees/internal/util/duration"
)

func newDeleteCmd() *cobra.Command {
	var (
		closedBefore string
		yes          bool
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete closed issues in bulk",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if closedBefore == "" {
				return fmt.Errorf("--closed-before is required")
			}

			cutoff, err := duration.Parse(closedBefore)
			if err != nil {
				return fmt.Errorf("failed to parse --closed-before value %q: %w", closedBefore, err)
			}

			filter := domain.DeleteFilter{ClosedBefore: cutoff}

			return h.Delete(cmd.Context(), cmd.OutOrStdout(), filter, yes)
		},
	}

	cmd.Flags().StringVar(&closedBefore, "closed-before", "", "Delete issues closed before duration (e.g. 12mo, 1y)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm deletion (without this flag, only previews)")

	return cmd
}
