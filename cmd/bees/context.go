package main

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

func newContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Show current repo context",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			summary, err := svc.Context(cmd.Context())
			if err != nil {
				return err
			}

			if !jsonOutput {
				printContextSummary(summary)
				return nil
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", " ")

			return enc.Encode(summary)
		},
	}

	return cmd
}
