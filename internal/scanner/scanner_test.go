package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/secretshield/cli/internal/config"
	"github.com/secretshield/cli/internal/rules"
)

// TestNewScanner verifies that a new Scanner is created with default rules.
func TestNewScanner(t *testing.T) {
	cfg := config.DefaultConfig()
	s := NewScanner(cfg)

	if s == nil {
		t.Fatal("NewScanner returned nil")
	}

	if s.GetRuleCount() == 0 {
		t.Error("Expected non-zero rule count after initialization")
	}
}

// TestScanFileWithAWSSecret tests detection of AWS Access Key ID in a file.
func TestScanFileWithAWSSecret(t *testing.T) {
	// Create a temporary file with an AWS key
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "config.go")

	content := `package main

const awsAccessKey = "AKIAIOSFODNN7EXAMPLE"
const awsSecret = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := &config.ScanConfig{
		Target:   testFile,
		FileMode: true,
		Output:   "json",
		Excludes: []string{},
	}

	s := NewScanner(cfg)
	findings, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(findings) == 0 {
		t.Error("Expected to find AWS Access Key ID in test file")
	}

	// Check that at least one finding is for AWS
	foundAWS := false
	for _, f := range findings {
		if f.RuleID == "GEN-001" {
			foundAWS = true
			break
		}
	}
	if !foundAWS {
		t.Error("Expected to find GEN-001 (AWS Access Key ID) rule match")
	}
}

// TestScanFileWithAlibabaCloud tests detection of Alibaba Cloud AccessKey.
func TestScanFileWithAlibabaCloud(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "settings.py")

	content := `# Alibaba Cloud credentials
ALIBABA_ACCESS_KEY_ID = "LTAI5t7example1234567"
ALIBABA_ACCESS_KEY_SECRET = "abcdef1234567890abcdef12345678"
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := &config.ScanConfig{
		Target:   testFile,
		FileMode: true,
		Output:   "json",
		Excludes: []string{},
	}

	s := NewScanner(cfg)
	findings, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(findings) == 0 {
		t.Error("Expected to find Alibaba Cloud AccessKey in test file")
	}

	foundAlibaba := false
	for _, f := range findings {
		if f.RuleID == "CN-001" {
			foundAlibaba = true
			break
		}
	}
	if !foundAlibaba {
		t.Error("Expected to find CN-001 (Alibaba Cloud AccessKey ID) rule match")
	}
}

// TestScanFileWithSSHKey tests detection of private SSH keys.
func TestScanFileWithSSHKey(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "id_rsa")

	content := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0Z3VS5JJcds3xfn/ygWyF8PbnGy0AHB7M2nNjYp5z
-----END RSA PRIVATE KEY-----
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := &config.ScanConfig{
		Target:   testFile,
		FileMode: true,
		Output:   "json",
		Excludes: []string{},
	}

	s := NewScanner(cfg)
	findings, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(findings) == 0 {
		t.Error("Expected to find RSA private key in test file")
	}

	foundSSHKey := false
	for _, f := range findings {
		if f.RuleID == "GEN-019" {
			foundSSHKey = true
			break
		}
	}
	if !foundSSHKey {
		t.Error("Expected to find GEN-019 (RSA Private Key) rule match")
	}
}

// TestScanCleanFile tests that a clean file produces no findings.
func TestScanCleanFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "clean.go")

	content := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := &config.ScanConfig{
		Target:   testFile,
		FileMode: true,
		Output:   "json",
		Excludes: []string{},
	}

	s := NewScanner(cfg)
	findings, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(findings) != 0 {
		t.Errorf("Expected no findings in clean file, got %d", len(findings))
	}
}

// TestScanDirectoryWithExclusions tests that excluded directories are skipped.
func TestScanDirectoryWithExclusions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file in an excluded directory
	vendorDir := filepath.Join(tmpDir, "vendor")
	os.MkdirAll(vendorDir, 0755)
	vendorFile := filepath.Join(vendorDir, "deps.go")
	os.WriteFile(vendorFile, []byte(`const key = "AKIAIOSFODNN7EXAMPLE"`), 0644)

	// Create a file in a non-excluded directory
	srcDir := filepath.Join(tmpDir, "src")
	os.MkdirAll(srcDir, 0755)
	srcFile := filepath.Join(srcDir, "main.go")
	os.WriteFile(srcFile, []byte(`const key = "AKIAIOSFODNN7EXAMPLE"`), 0644)

	cfg := &config.ScanConfig{
		Target:   tmpDir,
		Output:   "json",
		Excludes: []string{"vendor"},
	}

	s := NewScanner(cfg)
	findings, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Should only find one match (from src/main.go, not vendor/deps.go)
	if len(findings) != 1 {
		t.Errorf("Expected 1 finding (vendor excluded), got %d", len(findings))
	}
}

// TestSeverityFilter tests that severity filtering works correctly.
func TestSeverityFilter(t *testing.T) {
	// Use GetRulesBySeverity directly to verify the rules package filter logic
	filtered := rules.GetRulesBySeverity([]string{"critical"})

	if len(filtered) == 0 {
		t.Error("Expected at least one critical rule")
	}

	// After severity filter, only critical rules should remain
	for _, r := range filtered {
		if r.Severity != rules.SeverityCritical {
			t.Errorf("Expected only critical rules after filter, found %s (%s)", r.ID, r.Severity)
		}
	}
}

// TestRulesLoaded verifies that both generic and china rules are loaded.
func TestRulesLoaded(t *testing.T) {
	genericRules := rules.GetRulesByCategory(rules.CategoryGeneric)
	chinaRules := rules.GetRulesByCategory(rules.CategoryChina)

	if len(genericRules) < 15 {
		t.Errorf("Expected at least 15 generic rules, got %d", len(genericRules))
	}

	if len(chinaRules) < 10 {
		t.Errorf("Expected at least 10 china rules, got %d", len(chinaRules))
	}
}
