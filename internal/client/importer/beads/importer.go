package beads

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/w-h-a/bees/internal/client/importer"
	"github.com/w-h-a/bees/internal/domain"
)

type beadsImporter struct {
	options importer.Options
}

func (b *beadsImporter) Parse(r io.Reader) ([]domain.Issue, error) {
	var issues []domain.Issue

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var bi beadsIssue
		if err := json.Unmarshal(line, &bi); err != nil {
			return nil, fmt.Errorf("failed to unmarshal line %d: %w", lineNum, err)
		}

		issues = append(issues, mapToIssue(bi))
	}

	return issues, scanner.Err()
}

func NewImporter(opts ...importer.Option) (importer.Importer, error) {
	options := importer.NewOptions(opts...)

	i := &beadsImporter{
		options: options,
	}

	return i, nil
}
