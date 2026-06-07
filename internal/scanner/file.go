package scanner

import (
	"fmt"
	"os"

	"github.com/secretshield/cli/internal/rules"
)

// ScanSingleFile is a convenience function that scans a single file
// using the provided rules and returns findings.
func ScanSingleFile(filePath string, ruleList []rules.Rule) ([]rules.Finding, error) {
	// Verify the file exists
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot access file %s: %w", filePath, err)
	}

	// Verify it is not a directory
	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory, not a file", filePath)
	}

	// Create a minimal scanner just for file scanning
	s := &Scanner{
		rules: ruleList,
	}

	return s.scanFile(filePath)
}

// GetSupportedExtensions returns common file extensions that may contain secrets.
func GetSupportedExtensions() []string {
	return []string{
		".go", ".py", ".js", ".ts", ".jsx", ".tsx", ".java", ".kt",
		".rb", ".php", ".c", ".cpp", ".h", ".hpp", ".cs", ".swift",
		".rs", ".scala", ".r", ".lua", ".pl", ".pm", ".ex", ".exs",
		".erl", ".hrl", ".hs", ".ml", ".mli", ".fs", ".fsx",
		".sh", ".bash", ".zsh", ".fish", ".ps1", ".bat", ".cmd",
		".yml", ".yaml", ".json", ".toml", ".ini", ".cfg", ".conf",
		".env", ".properties", ".xml", ".html", ".htm", ".css",
		".sql", ".graphql", ".tf", ".tfvars", ".hcl",
		".md", ".rst", ".txt", ".csv",
		".dockerfile", ".makefile",
		".gradle", ".pom",
	}
}
