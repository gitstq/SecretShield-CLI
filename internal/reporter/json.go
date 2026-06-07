package reporter

import (
	"encoding/json"
	"io"

	"github.com/secretshield/cli/internal/rules"
)

// JSONReporter outputs findings in standard JSON format.
type JSONReporter struct{}

// Report writes findings as a JSON array to the writer.
func (r *JSONReporter) Report(w io.Writer, findings []rules.Finding) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(findings)
}
