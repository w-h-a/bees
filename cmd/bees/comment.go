package main

import (
	"os"

	"github.com/spf13/cobra"
)

func newCommentCmd() *cobra.Command {
	var author string

	cmd := &cobra.Command{
		Use:   "comment <id> <text>",
		Short: "Add a comment to an issue",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if author == "" {
				author = os.Getenv("BEES_AUTHOR")
			}
			return h.Comment(cmd.Context(), cmd.OutOrStdout(), args[0], author, args[1])
		},
	}

	cmd.Flags().StringVar(&author, "author", "", "Comment author")

	return cmd
}
