package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/secretshield/cli/internal/rules"
)

// TableReporter outputs findings in a terminal-friendly table format with color coding.
type TableReporter struct{}

// ANSI color codes for terminal output.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
	colorGray   = "\033[90m"
)

// Report writes findings as a formatted table to the writer.
func (r *TableReporter) Report(w io.Writer, findings []rules.Finding) error {
	if len(findings) == 0 {
		fmt.Fprintln(w, colorGreen+"No secrets found. Your code is clean!"+colorReset)
		return nil
	}

	// Print header
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "%s%sSecretShield Scan Results%s%s\n", colorBold, colorCyan, colorReset, colorReset)
	fmt.Fprintf(w, "%s%s%s\n", colorGray, strings.Repeat("─", 90), colorReset)
	fmt.Fprintf(w, "%s%-6s %-8s %-8s %-12s %-10s %-30s%s\n",
		colorBold, "SEVERITY", "RULE ID", "LINE", "CATEGORY", "FILE", "MATCHED CONTENT", colorReset)
	fmt.Fprintf(w, "%s%s%s\n", colorGray, strings.Repeat("─", 90), colorReset)

	// Print each finding
	for _, f := range findings {
		severity := colorizeSeverity(f.Severity)
		fmt.Fprintf(w, "%-6s %-8s %-8d %-12s %-10s %s\n",
			severity,
			f.RuleID,
			f.Line,
			f.Category,
			truncatePath(f.File, 40),
			f.MatchedContent,
		)
	}

	// Print summary
	fmt.Fprintf(w, "%s%s%s\n", colorGray, strings.Repeat("─", 90), colorReset)
	fmt.Fprintf(w, "\n%sTotal findings: %d%s\n", colorBold, len(findings), colorReset)

	// Print severity breakdown
	critical := countBySeverity(findings, "critical")
	high := countBySeverity(findings, "high")
	medium := countBySeverity(findings, "medium")
	low := countBySeverity(findings, "low")

	if critical > 0 {
		fmt.Fprintf(w, "  %sCritical: %d%s\n", colorRed, critical, colorReset)
	}
	if high > 0 {
		fmt.Fprintf(w, "  %sHigh: %d%s\n", colorRed, high, colorReset)
	}
	if medium > 0 {
		fmt.Fprintf(w, "  %sMedium: %d%s\n", colorYellow, medium, colorReset)
	}
	if low > 0 {
		fmt.Fprintf(w, "  %sLow: %d%s\n", colorGreen, low, colorReset)
	}

	fmt.Fprintln(w, "")
	return nil
}

// colorizeSeverity returns the severity string with appropriate color.
func colorizeSeverity(severity string) string {
	switch severity {
	case "critical":
		return colorRed + "CRIT" + colorReset
	case "high":
		return colorRed + "HIGH" + colorReset
	case "medium":
		return colorYellow + "MED" + colorReset
	case "low":
		return colorGreen + "LOW" + colorReset
	default:
		return severity
	}
}

// truncatePath truncates a file path to fit within the specified width.
func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
}

// countBySeverity counts findings by severity level.
func countBySeverity(findings []rules.Finding, severity string) int {
	count := 0
	for _, f := range findings {
		if f.Severity == severity {
			count++
		}
	}
	return count
}
