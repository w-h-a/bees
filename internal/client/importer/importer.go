package importer

import (
	"io"

	"github.com/w-h-a/bees/v2/internal/domain"
)

type Importer interface {
	Parse(r io.Reader) ([]domain.Issue, error)
}
