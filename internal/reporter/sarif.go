package reporter

import (
	"encoding/json"
	"io"
	"time"

	"github.com/secretshield/cli/internal/rules"
)

// SARIFReporter outputs findings in Microsoft SARIF v2.1.0 format.
type SARIFReporter struct{}

// sarifDocument represents the top-level SARIF document structure.
type sarifDocument struct {
	Schema  string       `json:"$schema"`
	Version string       `json:"version"`
	Runs    []sarifRun   `json:"runs"`
}

// sarifRun represents a single run in a SARIF document.
type sarifRun struct {
	Tool    sarifTool    `json:"tool"`
	Results []sarifResult `json:"results"`
}

// sarifTool represents the tool information in a SARIF run.
type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

// sarifDriver represents the tool driver in SARIF.
type sarifDriver struct {
	Name           string         `json:"name"`
	Version        string         `json:"version"`
	InformationURI string         `json:"informationUri"`
	Rules          []sarifRuleDef `json:"rules"`
}

// sarifRuleDef represents a rule definition in SARIF.
type sarifRuleDef struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ShortDescription struct {
		Text string `json:"text"`
	} `json:"shortDescription"`
	FullDescription struct {
		Text string `json:"text"`
	} `json:"fullDescription"`
	HelpURI string `json:"helpUri,omitempty"`
}

// sarifResult represents a single result (finding) in SARIF.
type sarifResult struct {
	RuleID    string         `json:"ruleId"`
	RuleIndex int            `json:"ruleIndex"`
	Level     string         `json:"level"`
	Message   sarifMessage   `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

// sarifMessage represents a message in SARIF.
type sarifMessage struct {
	Text string `json:"text"`
}

// sarifLocation represents a physical location in SARIF.
type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

// sarifPhysicalLocation represents a physical location in SARIF.
type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region"`
}

// sarifArtifactLocation represents an artifact location in SARIF.
type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

// sarifRegion represents a region within an artifact in SARIF.
type sarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
}

// Report writes findings in SARIF v2.1.0 format to the writer.
func (r *SARIFReporter) Report(w io.Writer, findings []rules.Finding) error {
	// Build unique rule definitions
	ruleMap := make(map[string]rules.Finding)
	ruleOrder := make([]string, 0)
	for _, f := range findings {
		if _, exists := ruleMap[f.RuleID]; !exists {
			ruleMap[f.RuleID] = f
			ruleOrder = append(ruleOrder, f.RuleID)
		}
	}

	// Build rule definitions
	ruleDefs := make([]sarifRuleDef, 0, len(ruleMap))
	for _, id := range ruleOrder {
		f := ruleMap[id]
		ruleDefs = append(ruleDefs, sarifRuleDef{
			ID:   f.RuleID,
			Name: f.RuleName,
			ShortDescription: struct {
				Text string `json:"text"`
			}{Text: f.Description},
			FullDescription: struct {
				Text string `json:"text"`
			}{Text: f.Description},
		})
	}

	// Build results
	ruleIndexMap := make(map[string]int)
	for i, rd := range ruleDefs {
		ruleIndexMap[rd.ID] = i
	}

	results := make([]sarifResult, 0, len(findings))
	for _, f := range findings {
		level := severityToSARIFLevel(f.Severity)
		results = append(results, sarifResult{
			RuleID:    f.RuleID,
			RuleIndex: ruleIndexMap[f.RuleID],
			Level:     level,
			Message: sarifMessage{
				Text: f.Description,
			},
			Locations: []sarifLocation{
				{
					PhysicalLocation: sarifPhysicalLocation{
						ArtifactLocation: sarifArtifactLocation{
							URI: f.File,
						},
						Region: sarifRegion{
							StartLine:   f.Line,
							StartColumn: f.Column,
						},
					},
				},
			},
		})
	}

	doc := sarifDocument{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:           "SecretShield-CLI",
						Version:        "1.0.0",
						InformationURI: "https://github.com/secretshield/cli",
						Rules:          ruleDefs,
					},
				},
				Results: results,
			},
		},
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(doc)
}

// severityToSARIFLevel maps SecretShield severity to SARIF level.
func severityToSARIFLevel(severity string) string {
	switch severity {
	case "critical":
		return "error"
	case "high":
		return "error"
	case "medium":
		return "warning"
	case "low":
		return "note"
	default:
		return "none"
	}
}

// init ensures time package is used (suppresses unused import if needed).
var _ = time.Now
