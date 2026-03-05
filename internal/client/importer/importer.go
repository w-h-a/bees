package importer

import (
	"io"

	"github.com/w-h-a/bees/internal/domain"
)

type Importer interface {
	Parse(r io.Reader) ([]domain.Issue, error)
}
