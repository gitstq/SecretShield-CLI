// Package rules defines detection rules for secret/credential scanning.
// It provides rule definitions, loading, and matching capabilities.
package rules

import "strings"

// Severity represents the severity level of a finding.
type Severity string

const (
	// SeverityCritical indicates a critical finding (e.g., active private keys).
	SeverityCritical Severity = "critical"
	// SeverityHigh indicates a high severity finding (e.g., API keys).
	SeverityHigh Severity = "high"
	// SeverityMedium indicates a medium severity finding.
	SeverityMedium Severity = "medium"
	// SeverityLow indicates a low severity finding.
	SeverityLow Severity = "low"
)

// Category represents the category of a detection rule.
type Category string

const (
	// CategoryGeneric represents generic/international cloud provider rules.
	CategoryGeneric Category = "generic"
	// CategoryChina represents Chinese cloud provider rules.
	CategoryChina Category = "china"
)

// Rule defines a single detection rule for secret scanning.
type Rule struct {
	ID          string   // Unique rule identifier, e.g. "GEN-001"
	Name        string   // Human-readable rule name
	Description string   // Description of what this rule detects
	Severity    Severity // Severity level
	Category    Category // Rule category
	Keywords    []string // Keywords to pre-filter lines before regex matching
	Pattern     string   // Regular expression pattern to match
	ChineseName string   // Chinese name for Chinese report output
}

// Finding represents a single secret detection result.
type Finding struct {
	File           string `json:"file"`
	Line           int    `json:"line"`
	Column         int    `json:"column"`
	RuleID         string `json:"rule_id"`
	RuleName       string `json:"rule_name"`
	Severity       string `json:"severity"`
	Category       string `json:"category"`
	MatchedContent string `json:"matched_content"`
	Description    string `json:"description"`
}

// AllRules contains all built-in detection rules.
var AllRules []Rule

func init() {
	AllRules = append(AllRules, GenericRules()...)
	AllRules = append(AllRules, ChinaRules()...)
}

// GetRulesByCategory returns rules filtered by category.
func GetRulesByCategory(cat Category) []Rule {
	var result []Rule
	for _, r := range AllRules {
		if r.Category == cat {
			result = append(result, r)
		}
	}
	return result
}

// GetRulesBySeverity returns rules filtered by severity levels.
func GetRulesBySeverity(severities []string) []Rule {
	severitySet := make(map[string]bool)
	for _, s := range severities {
		severitySet[strings.ToLower(s)] = true
	}

	var result []Rule
	for _, r := range AllRules {
		if severitySet[string(r.Severity)] {
			result = append(result, r)
		}
	}
	return result
}

// MatchRule checks if a line of text matches the given rule.
// It first checks for keywords as a fast pre-filter, then applies the regex pattern.
func MatchRule(line string, rule Rule) (bool, string) {
	// Fast pre-filter: check if any keyword exists in the line
	if len(rule.Keywords) > 0 {
		found := false
		for _, kw := range rule.Keywords {
			if strings.Contains(line, kw) {
				found = true
				break
			}
		}
		if !found {
			return false, ""
		}
	}

	// Apply regex pattern matching
	matched := rulePatternMatch(line, rule.Pattern)
	if matched != "" {
		return true, matched
	}
	return false, ""
}
