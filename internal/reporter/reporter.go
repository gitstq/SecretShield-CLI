// Package reporter provides output formatting for SecretShield-CLI scan results.
// It supports multiple output formats including JSON, SARIF, table, and Chinese.
package reporter

import (
	"fmt"
	"io"
	"os"

	"github.com/secretshield/cli/internal/rules"
)

// Reporter defines the interface for output formatting.
type Reporter interface {
	// Report writes the findings to the given writer.
	Report(w io.Writer, findings []rules.Finding) error
}

// NewReporter creates a new Reporter based on the specified format.
func NewReporter(format string) (Reporter, error) {
	switch format {
	case "json":
		return &JSONReporter{}, nil
	case "sarif":
		return &SARIFReporter{}, nil
	case "table":
		return &TableReporter{}, nil
	case "chinese":
		return &ChineseReporter{}, nil
	default:
		return nil, fmt.Errorf("unknown output format: %s", format)
	}
}

// WriteReport creates a reporter and writes findings to the specified output.
// If reportPath is empty, it writes to stdout.
func WriteReport(format string, findings []rules.Finding, reportPath string) error {
	reporter, err := NewReporter(format)
	if err != nil {
		return err
	}

	if reportPath != "" {
		f, err := os.Create(reportPath)
		if err != nil {
			return fmt.Errorf("failed to create report file: %w", err)
		}
		defer f.Close()
		return reporter.Report(f, findings)
	}

	return reporter.Report(os.Stdout, findings)
}
