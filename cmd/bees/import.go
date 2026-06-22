package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import <file.jsonl>",
		Short: "Import issues from a JSONL file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer f.Close()

			return h.Import(cmd.Context(), cmd.OutOrStdout(), f)
		},
	}
}
