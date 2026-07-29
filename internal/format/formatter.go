package format

import (
	"fmt"
	"io"

	"github.com/Fauxmen4/deps-check/internal/analyzer"
)

// Format is the output format identifier.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

// Formatter renders and prints dependency report.
type Formatter interface {
	Write(w io.Writer, report []analyzer.DependencyReport) error
}

func New(f Format) (Formatter, error) {
	switch f {
	case FormatTable:
		return &TableFormatter{}, nil
	case FormatJSON:
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unknown format: %s", f)
	}
}
