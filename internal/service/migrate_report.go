package service

import "github.com/w-h-a/bees/v2/internal/domain"

// MigrateReport is the dry-run result the handler renders: per-source import
// counts, the source paths skipped because they had no store, and any full-ID
// collisions that force a refusal.
type MigrateReport struct {
	Sources    []domain.SourceCounts
	Skipped    []string
	Collisions []string
}
