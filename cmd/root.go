// Package cmd contains the CLI command definitions for SecretShield-CLI.
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/secretshield/cli/internal/config"
	"github.com/secretshield/cli/internal/notifier"
	"github.com/secretshield/cli/internal/reporter"
	"github.com/secretshield/cli/internal/rules"
	"github.com/secretshield/cli/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	// Scan flags
	flagGit       bool
	flagFile      bool
	flagOutput    string
	flagReport    string
	flagRules     string
	flagExclude   string
	flagSeverity  string
	flagWebhook   string

	// Rules command flags
	flagCategory string
)

// Version holds the application version string.
var Version = "1.0.0"

// rootCmd is the base command for SecretShield-CLI.
var rootCmd = &cobra.Command{
	Use:   "secretshield",
	Short: "SecretShield-CLI - Lightweight terminal secret leak scanner",
	Long: `SecretShield-CLI is a lightweight terminal-based secret leak scanning engine.

It scans Git repositories, directories, and files for leaked credentials,
API keys, and other sensitive information. It includes built-in rules for
international cloud providers (AWS, GCP, Azure, etc.) and Chinese cloud
providers (Alibaba Cloud, Tencent Cloud, Huawei Cloud, etc.).`,
	SilenceUsage: true,
}

// scanCmd represents the scan command.
var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan a directory, git repository, or file for leaked secrets",
	Long: `Scan a directory, git repository, or file for leaked secrets.

Examples:
  secretshield scan                    # Scan current directory
  secretshield scan /path/to/repo       # Scan a specific directory
  secretshield scan --git /path/to/repo # Scan git repository with history
  secretshield scan --file myfile.py   # Scan a single file
  secretshield scan --output json      # Output in JSON format
  secretshield scan --output sarif     # Output in SARIF format
  secretshield scan --output chinese   # Output in Chinese format
  secretshield scan --exclude vendor,node_modules
  secretshield scan --severity high,critical`,
	Args: cobra.MaximumNArgs(1),
	RunE: runScan,
}

// rulesCmd represents the rules command.
var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "List all built-in detection rules",
	Long: `List all built-in detection rules.

Examples:
  secretshield rules              # List all rules
  secretshield rules --category china    # List Chinese cloud provider rules
  secretshield rules --category generic  # List generic/international rules`,
	RunE: runRules,
}

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("SecretShield-CLI v%s\n", Version)
		fmt.Printf("Built with Go | Zero external dependencies (except Cobra)\n")
	},
}

func init() {
	// Register subcommands
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(rulesCmd)
	rootCmd.AddCommand(versionCmd)

	// Scan command flags
	scanCmd.Flags().BoolVarP(&flagGit, "git", "g", false, "Scan git repository including commit history")
	scanCmd.Flags().BoolVarP(&flagFile, "file", "f", false, "Scan a single file instead of a directory")
	scanCmd.Flags().StringVarP(&flagOutput, "output", "o", "table", "Output format: json, sarif, table, chinese")
	scanCmd.Flags().StringVarP(&flagReport, "report", "r", "", "Write report to file instead of stdout")
	scanCmd.Flags().StringVarP(&flagRules, "rules", "", "", "Path to custom rules file")
	scanCmd.Flags().StringVarP(&flagExclude, "exclude", "e", "", "Comma-separated list of directories to exclude")
	scanCmd.Flags().StringVarP(&flagSeverity, "severity", "s", "", "Comma-separated severity filter: critical,high,medium,low")
	scanCmd.Flags().StringVarP(&flagWebhook, "webhook", "w", "", "Webhook URL for notifications (Feishu/DingTalk)")

	// Rules command flags
	rulesCmd.Flags().StringVarP(&flagCategory, "category", "c", "", "Filter rules by category: china, generic")
}

// runScan executes the scan command.
func runScan(cmd *cobra.Command, args []string) error {
	startTime := time.Now()

	// Determine target path
	target := "."
	if len(args) > 0 {
		target = args[0]
	}

	// Build configuration
	cfg := config.DefaultConfig()
	cfg.Target = target
	cfg.GitMode = flagGit
	cfg.FileMode = flagFile
	cfg.Output = flagOutput
	cfg.Report = flagReport
	cfg.RulesFile = flagRules
	cfg.Webhook = flagWebhook

	if flagExclude != "" {
		cfg.Excludes = strings.Split(flagExclude, ",")
	}

	if flagSeverity != "" {
		cfg.Severities = strings.Split(flagSeverity, ",")
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Create scanner
	s := scanner.NewScanner(cfg)

	// Print scan header
	fmt.Fprintf(os.Stderr, "SecretShield-CLI v%s\n", Version)
	fmt.Fprintf(os.Stderr, "Scanning: %s\n", target)
	fmt.Fprintf(os.Stderr, "Rules loaded: %d\n", s.GetRuleCount())
	fmt.Fprintf(os.Stderr, "Output format: %s\n", cfg.Output)
	fmt.Fprintf(os.Stderr, "\n")

	// Run scan
	findings, err := s.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// Write report
	if err := reporter.WriteReport(cfg.Output, findings, cfg.Report); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	// Print scan footer to stderr
	elapsed := time.Since(startTime)
	fmt.Fprintf(os.Stderr, "\nScan completed in %s\n", elapsed)

	// Send webhook notification if configured
	if cfg.Webhook != "" {
		fmt.Fprintf(os.Stderr, "Sending webhook notification...\n")
		n := notifier.NewNotifier(cfg.Webhook)
		if err := n.Notify(findings); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: webhook notification failed: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Webhook notification sent successfully.\n")
		}
	}

	// Exit with non-zero code if findings exist
	if len(findings) > 0 {
		os.Exit(1)
	}

	return nil
}

// runRules executes the rules command.
func runRules(cmd *cobra.Command, args []string) error {
	var ruleList []rules.Rule

	if flagCategory != "" {
		switch strings.ToLower(flagCategory) {
		case "china":
			ruleList = rules.GetRulesByCategory(rules.CategoryChina)
		case "generic":
			ruleList = rules.GetRulesByCategory(rules.CategoryGeneric)
		default:
			return fmt.Errorf("unknown category: %s (valid: china, generic)", flagCategory)
		}
	} else {
		ruleList = rules.AllRules
	}

	fmt.Printf("SecretShield-CLI v%s - Detection Rules\n", Version)
	fmt.Printf("Total rules: %d\n\n", len(ruleList))

	// Group by category
	genericCount := 0
	chinaCount := 0
	for _, r := range ruleList {
		switch r.Category {
		case rules.CategoryGeneric:
			genericCount++
		case rules.CategoryChina:
			chinaCount++
		}
	}

	fmt.Printf("Generic rules: %d\n", genericCount)
	fmt.Printf("China rules:   %d\n\n", chinaCount)

	// Print rules table
	fmt.Printf("%-8s %-35s %-10s %-10s %s\n", "ID", "NAME", "SEVERITY", "CATEGORY", "DESCRIPTION")
	fmt.Printf("%s\n", strings.Repeat("-", 90))

	for _, r := range ruleList {
		fmt.Printf("%-8s %-35s %-10s %-10s %s\n",
			r.ID,
			truncate(r.Name, 33),
			r.Severity,
			r.Category,
			truncate(r.Description, 50),
		)
	}

	fmt.Printf("\nTotal: %d rules\n", len(ruleList))
	return nil
}

// truncate truncates a string to the specified length.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
