package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newCommentCmd() *cobra.Command {
	var author string

	cmd := &cobra.Command{
		Use:   "comment <id> <text>",
		Short: "Add a comment to an issue",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			comment, err := svc.AddComment(cmd.Context(), args[0], author, args[1])
			if err != nil {
				return err
			}

			if !jsonOutput {
				name := comment.Author
				if name == "" {
					name = "anonymous"
				}
				ts := comment.CreatedAt.Format("2006-01-02 15:04")
				fmt.Printf("%s %s\n", dimStyle.Render(ts), headerStyle.Render(name))
				for line := range strings.SplitSeq(comment.Body, "\n") {
					fmt.Printf("  %s\n", line)
				}
				return nil
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")

			return enc.Encode(comment)
		},
	}

	cmd.Flags().StringVar(&author, "author", "", "Comment author")

	return cmd
}
