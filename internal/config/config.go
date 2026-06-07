// Package config manages configuration for SecretShield-CLI.
// It handles loading and parsing of scan options and settings.
package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ScanConfig holds the configuration for a scan operation.
type ScanConfig struct {
	Target    string   // Target path to scan (directory, git repo, or file)
	GitMode   bool     // Whether to scan git history
	FileMode  bool     // Whether to scan a single file
	Output    string   // Output format: json, sarif, table, chinese
	Report    string   // Output file path (empty means stdout)
	RulesFile string   // Custom rules file path
	Excludes  []string // Excluded directory patterns
	Severities []string // Severity filter
	Webhook   string   // Webhook URL for notifications
}

// DefaultConfig returns a ScanConfig with default values.
func DefaultConfig() *ScanConfig {
	return &ScanConfig{
		Target:   ".",
		Output:   "table",
		Excludes: []string{"vendor", "node_modules", ".git", ".svn", "__pycache__", ".idea", ".vscode"},
	}
}

// Validate checks the configuration for validity and returns any errors.
func (c *ScanConfig) Validate() error {
	if c.FileMode && c.GitMode {
		return fmt.Errorf("cannot use both --file and --git flags simultaneously")
	}

	if c.Target == "" {
		return fmt.Errorf("target path cannot be empty")
	}

	validFormats := map[string]bool{
		"json":    true,
		"sarif":   true,
		"table":   true,
		"chinese": true,
	}
	if !validFormats[c.Output] {
		return fmt.Errorf("invalid output format: %s (valid: json, sarif, table, chinese)", c.Output)
	}

	return nil
}

// LoadCustomRules loads custom rules from a file.
// Each line should be in the format: "rule_id|rule_name|pattern|severity|category"
func LoadCustomRules(path string) ([]CustomRule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open custom rules file: %w", err)
	}
	defer f.Close()

	var rules []CustomRule
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			return nil, fmt.Errorf("invalid rule format at line %d: expected 'id|name|pattern|severity|category'", lineNum)
		}

		rules = append(rules, CustomRule{
			ID:       parts[0],
			Name:     parts[1],
			Pattern:  parts[2],
			Severity: parts[3],
			Category: parts[4],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading custom rules file: %w", err)
	}

	return rules, nil
}

// CustomRule represents a user-defined detection rule loaded from a file.
type CustomRule struct {
	ID       string
	Name     string
	Pattern  string
	Severity string
	Category string
}
