package scanner

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/secretshield/cli/internal/rules"
	"github.com/secretshield/cli/pkg/utils"
)

// scanGitRepo scans a git repository including historical commits.
// It uses `git log` to find all blobs and scans their content.
func (s *Scanner) scanGitRepo(repoPath string) ([]rules.Finding, error) {
	// Verify the path is a git repository
	if !isGitRepo(repoPath) {
		return nil, fmt.Errorf("%s is not a git repository", repoPath)
	}

	var findings []rules.Finding

	// Get all unique file paths from git history
	files, err := s.getGitFiles(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get git files: %w", err)
	}

	// Scan each file's content from git show
	for _, file := range files {
		// Skip excluded paths
		if utils.IsExcluded(file, s.config.Excludes) {
			continue
		}

		// Get file content from HEAD
		content, err := s.getGitFileContent(repoPath, file)
		if err != nil {
			continue // Skip files that cannot be read
		}

		// Check if content looks like binary
		if isBinaryContent(content) {
			continue
		}

		fileFindings := s.scanContent(content, file)
		findings = append(findings, fileFindings...)
	}

	// Also scan historical diffs for secrets that were removed
	diffFindings, err := s.scanGitDiffs(repoPath)
	if err == nil {
		findings = append(findings, diffFindings...)
	}

	return findings, nil
}

// isGitRepo checks if the given path is a git repository.
func isGitRepo(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// getGitFiles retrieves all tracked file paths from a git repository.
func (s *Scanner) getGitFiles(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var files []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		file := strings.TrimSpace(scanner.Text())
		if file != "" {
			files = append(files, file)
		}
	}
	return files, nil
}

// getGitFileContent retrieves the content of a file from git HEAD.
func (s *Scanner) getGitFileContent(repoPath, filePath string) (string, error) {
	cmd := exec.Command("git", "show", "HEAD:"+filePath)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// scanGitDiffs scans git diff history for potential secrets.
func (s *Scanner) scanGitDiffs(repoPath string) ([]rules.Finding, error) {
	cmd := exec.Command("git", "log", "-p", "--all", "--diff-filter=ACDMR", "--")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	content := string(output)
	return s.scanContent(content, "git-history"), nil
}

// scanContent scans a string content for secrets, associating findings with the given filePath.
func (s *Scanner) scanContent(content, filePath string) []rules.Finding {
	var findings []rules.Finding
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip diff metadata lines
		if isDiffMetadata(line) {
			continue
		}

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

	return findings
}

// isDiffMetadata checks if a line is git diff metadata (not actual content).
func isDiffMetadata(line string) bool {
	prefixes := []string{
		"commit ", "Author: ", "Date:   ", "diff --git ",
		"index ", "--- ", "+++ ", "@@ ", "new file mode ",
		"deleted file mode ", "similarity index ", "rename from ",
		"rename to ", "copy from ", "copy to ", "Binary files ",
		"dissimilarity index ",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

// isBinaryContent checks if content appears to be binary.
func isBinaryContent(content string) bool {
	if len(content) == 0 {
		return false
	}
	checkLen := 512
	if len(content) < checkLen {
		checkLen = len(content)
	}
	for i := 0; i < checkLen; i++ {
		if content[i] == 0 {
			return true
		}
	}
	return false
}

// getAbsPath returns the absolute path for a relative path within the repo.
func getAbsPath(repoPath, relPath string) string {
	return filepath.Join(repoPath, relPath)
}
