package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newHandoffCmd() *cobra.Command {
	var (
		done      string
		remaining string
		decisions string
		uncertain string
	)

	cmd := &cobra.Command{
		Use:   "handoff <id>",
		Short: "Record a structured handoff for an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handoff, err := svc.AddHandoff(cmd.Context(), args[0], done, remaining, decisions, uncertain)
			if err != nil {
				return err
			}

			if !jsonOutput {
				ts := handoff.CreatedAt.Format("2006-01-02 15:04")
				fmt.Printf("%s %s\n", dimStyle.Render(ts), headerStyle.Render("Handoff recorded"))
				printHandoffInline(*handoff)
				return nil
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", " ")

			return enc.Encode(handoff)
		},
	}

	cmd.Flags().StringVar(&done, "done", "", "What was completed")
	cmd.Flags().StringVar(&remaining, "remaining", "", "What remains to be done")
	cmd.Flags().StringVar(&decisions, "decisions", "", "Decisions made")
	cmd.Flags().StringVar(&uncertain, "uncertain", "", "Open questions or uncertainties")

	return cmd
}
