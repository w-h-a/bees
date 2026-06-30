package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/w-h-a/bees/v2/internal/client/exporter/jsonl"
	noopexporter "github.com/w-h-a/bees/v2/internal/client/exporter/noop"
	"github.com/w-h-a/bees/v2/internal/client/importer/bees"
	noopimporter "github.com/w-h-a/bees/v2/internal/client/importer/noop"
	"github.com/w-h-a/bees/v2/internal/client/repo"
	"github.com/w-h-a/bees/v2/internal/client/repo/sqlite"
	sqlitesource "github.com/w-h-a/bees/v2/internal/client/source/sqlite"
	"github.com/w-h-a/bees/v2/internal/handler/cli"
	"github.com/w-h-a/bees/v2/internal/service"
	"github.com/w-h-a/bees/v2/internal/util/home"
	"github.com/w-h-a/bees/v2/internal/util/prefix"
)

var (
	jsonOutput      bool
	verbose         bool
	migrateTargetDB string
	svc             *service.Service
	h               *cli.Handler
	dbCloser        func() error
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bees",
		Short: "An alternative to a sea of .md files for developers who pair with agentic navigators.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if verbose || os.Getenv("BEES_DEBUG") == "1" {
				opts := &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}
				var handler slog.Handler
				if jsonOutput {
					handler = slog.NewJSONHandler(os.Stderr, opts)
				} else {
					handler = slog.NewTextHandler(os.Stderr, opts)
				}
				slog.SetDefault(slog.New(handler))
			}

			if cmd.Name() == "migrate" {
				reader, err := sqlitesource.NewReader()
				if err != nil {
					return fmt.Errorf("failed to initialize reader: %w", err)
				}

				i, _ := noopimporter.NewImporter()
				e, _ := noopexporter.NewExporter()

				userHome, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to resolve user home: %w", err)
				}

				resolved, err := home.Resolve(os.Getenv("BEES_HOME"), userHome)
				if err != nil {
					return err
				}
				migrateTargetDB = resolved.DB

				var r repo.Repo
				if commit, _ := cmd.Flags().GetBool("commit"); commit {
					if err := os.MkdirAll(resolved.Home, 0o755); err != nil {
						return fmt.Errorf("failed to create bees home: %w", err)
					}

					gr, err := sqlite.NewRepo(repo.WithLocation(resolved.DB))
					if err != nil {
						return fmt.Errorf("failed to open global store: %w", err)
					}

					r = gr
					dbCloser = gr.Close
				}

				svc = service.NewService(r, i, e, reader, "")
				h = cli.New(svc, jsonOutput)

				return nil
			}

			beesHome := os.Getenv("BEES_HOME")

			var userHome string
			if beesHome == "" {
				var err error
				userHome, err = os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to resolve user home: %w", err)
				}
			}

			resolved, err := home.Resolve(beesHome, userHome)
			if err != nil {
				return err
			}

			if !cmd.Flags().Changed("json") {
				if v := os.Getenv("BEES_JSON"); v == "true" || v == "1" {
					jsonOutput = true
				}
			}

			prefixFlag, _ := cmd.Flags().GetString("prefix")
			resolvedPrefix := prefix.Resolve(prefixFlag, os.Getenv("BEES_PREFIX"), "")

			slog.Debug("bees home resolved", "home", resolved.Home, "db", resolved.DB, "prefix", resolvedPrefix)

			slog.Debug("command path", "path", cmd.CommandPath(), "name", cmd.Name())

			needsDB := map[string]bool{
				"bees import":     true,
				"bees export":     true,
				"bees create":     true,
				"bees show":       true,
				"bees list":       true,
				"bees search":     true,
				"bees update":     true,
				"bees close":      true,
				"bees reopen":     true,
				"bees delete":     true,
				"bees context":    true,
				"bees ready":      true,
				"bees upcoming":   true,
				"bees dep add":    true,
				"bees dep remove": true,
				"bees dep graph":  true,
				"bees comment":    true,
				"bees handoff":    true,
			}
			if needsDB[cmd.CommandPath()] {
				if err := os.MkdirAll(resolved.Home, 0o755); err != nil {
					return fmt.Errorf("failed to create bees home: %w", err)
				}

				r, err := sqlite.NewRepo(repo.WithLocation(resolved.DB))
				if err != nil {
					return fmt.Errorf("failed to open database: %w", err)
				}
				dbCloser = r.Close

				i, _ := noopimporter.NewImporter()
				if cmd.CommandPath() == "bees import" {
					i, err = bees.NewImporter()
					if err != nil {
						return fmt.Errorf("failed to initialize importer: %w", err)
					}
				}

				e, _ := noopexporter.NewExporter()
				if cmd.CommandPath() == "bees export" {
					e, err = jsonl.NewExporter()
					if err != nil {
						return fmt.Errorf("failed to initialize exporter: %w", err)
					}
				}

				svc = service.NewService(r, i, e, nil, resolvedPrefix)
			}

			h = cli.New(svc, jsonOutput)

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if dbCloser == nil {
				return nil
			}
			return dbCloser()
		},
	}

	cmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	cmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable debug logging")

	cmd.AddCommand(newImportCmd())
	cmd.AddCommand(newExportCmd())
	cmd.AddCommand(newMigrateCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newShowCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newSearchCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newCloseCmd())
	cmd.AddCommand(newReopenCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newContextCmd())
	cmd.AddCommand(newReadyCmd())
	cmd.AddCommand(newUpcomingCmd())
	cmd.AddCommand(newDepCmd())
	cmd.AddCommand(newCommentCmd())
	cmd.AddCommand(newHandoffCmd())
	cmd.AddCommand(newVersionCmd())

	return cmd
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
