package noop

import (
	"fmt"
	"io"

	"github.com/w-h-a/bees/internal/client/importer"
	"github.com/w-h-a/bees/internal/domain"
)

type noopImporter struct {
	options importer.Options
}

func (i *noopImporter) Parse(r io.Reader) ([]domain.Issue, error) {
	return nil, fmt.Errorf("no importer configured")
}

func NewImporter(opts ...importer.Option) (importer.Importer, error) {
	options := importer.NewOptions(opts...)
	return &noopImporter{
		options: options,
	}, nil
}
