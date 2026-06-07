// Package scanner provides the core scanning engine for SecretShield-CLI.
// It scans files, directories, and git repositories for leaked secrets.
package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/secretshield/cli/internal/config"
	"github.com/secretshield/cli/internal/rules"
	"github.com/secretshield/cli/pkg/utils"
)

// Scanner is the main scanning engine that detects secrets in files.
type Scanner struct {
	config *config.ScanConfig
	rules  []rules.Rule
}

// NewScanner creates a new Scanner with the given configuration.
func NewScanner(cfg *config.ScanConfig) *Scanner {
	return &Scanner{
		config: cfg,
		rules:  rules.AllRules,
	}
}

// Scan performs the scan based on the configuration.
// It returns a list of findings and any error encountered.
func (s *Scanner) Scan() ([]rules.Finding, error) {
	// Apply severity filter if specified
	if len(s.config.Severities) > 0 {
		s.rules = rules.GetRulesBySeverity(s.config.Severities)
	}

	// Load custom rules if specified
	if s.config.RulesFile != "" {
		customRules, err := config.LoadCustomRules(s.config.RulesFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load custom rules: %w", err)
		}
		for _, cr := range customRules {
			s.rules = append(s.rules, rules.Rule{
				ID:          cr.ID,
				Name:        cr.Name,
				Description: fmt.Sprintf("Custom rule: %s", cr.Name),
				Severity:    rules.Severity(cr.Severity),
				Category:    rules.Category(cr.Category),
				Keywords:    []string{},
				Pattern:     cr.Pattern,
				ChineseName:  cr.Name,
			})
		}
	}

	var findings []rules.Finding
	var err error

	switch {
	case s.config.GitMode:
		findings, err = s.scanGitRepo(s.config.Target)
	case s.config.FileMode:
		findings, err = s.scanFile(s.config.Target)
	default:
		findings, err = s.scanDirectory(s.config.Target)
	}

	if err != nil {
		return nil, err
	}

	return findings, nil
}

// scanDirectory recursively scans a directory for secrets.
func (s *Scanner) scanDirectory(dirPath string) ([]rules.Finding, error) {
	var findings []rules.Finding

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip files we cannot access
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if path is excluded
		if utils.IsExcluded(path, s.config.Excludes) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip binary files
		if utils.IsBinaryFile(path) {
			return nil
		}

		// Scan the file
		fileFindings, err := s.scanFile(path)
		if err != nil {
			return nil // Continue scanning other files
		}

		findings = append(findings, fileFindings...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return findings, nil
}

// scanFile scans a single file for secrets.
func (s *Scanner) scanFile(filePath string) ([]rules.Finding, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer f.Close()

	var findings []rules.Finding
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		for _, rule := range s.rules {
			matched, matchedContent := rules.MatchRule(line, rule)
			if matched {
				column := findColumnPosition(line, rule.Pattern)
				findings = append(findings, rules.Finding{
					File:           filePath,
					Line:           lineNum,
					Column:         column,
					RuleID:         rule.ID,
					RuleName:       rule.Name,
					Severity:       string(rule.Severity),
					Category:       string(rule.Category),
					MatchedContent: utils.MaskSecret(matchedContent),
					Description:    rule.Description,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return findings, nil
}

// findColumnPosition finds the 1-based column position of a regex match in a line.
func findColumnPosition(line, pattern string) int {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0
	}
	loc := re.FindStringIndex(line)
	if loc == nil {
		return 0
	}
	return loc[0] + 1
}

// GetRuleCount returns the number of active rules.
func (s *Scanner) GetRuleCount() int {
	return len(s.rules)
}
