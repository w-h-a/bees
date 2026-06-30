package source

import (
	"context"
	"errors"

	"github.com/w-h-a/bees/v2/internal/domain"
)

var (
	ErrNotFound = errors.New("source db not found")
)

// Reader reads a v1 bees store. Implementations may not mutate the
// source.
type Reader interface {
	Read(ctx context.Context, dbPath string) ([]domain.Issue, error)
}
