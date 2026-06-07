package reporter

import (
	"encoding/json"
	"io"

	"github.com/secretshield/cli/internal/rules"
)

// ChineseReporter outputs findings in a Chinese-localized JSON format.
type ChineseReporter struct{}

// chineseFinding is a Chinese-localized version of a finding for display.
type chineseFinding struct {
	File           string `json:"文件"`
	Line           int    `json:"行号"`
	Column         int    `json:"列号"`
	RuleID         string `json:"规则ID"`
	RuleName       string `json:"规则名称"`
	Severity       string `json:"严重级别"`
	SeverityCN     string `json:"严重级别(中文)"`
	Category       string `json:"类别"`
	CategoryCN     string `json:"类别(中文)"`
	MatchedContent string `json:"匹配内容(脱敏)"`
	Description    string `json:"描述"`
	DescriptionCN  string `json:"描述(中文)"`
}

// chineseReport is the top-level Chinese report structure.
type chineseReport struct {
	ToolName    string             `json:"工具名称"`
	Version     string             `json:"版本"`
	TotalCount  int                `json:"发现总数"`
	Critical    int                `json:"严重数量"`
	High        int                `json:"高危数量"`
	Medium      int                `json:"中危数量"`
	Low         int                `json:"低危数量"`
	Findings    []chineseFinding   `json:"发现详情"`
}

// Report writes findings as a Chinese-localized JSON report to the writer.
func (r *ChineseReporter) Report(w io.Writer, findings []rules.Finding) error {
	report := chineseReport{
		ToolName:   "SecretShield-CLI 密钥泄露扫描引擎",
		Version:    "1.0.0",
		TotalCount: len(findings),
		Critical:   0,
		High:       0,
		Medium:     0,
		Low:        0,
	}

	report.Findings = make([]chineseFinding, 0, len(findings))

	for _, f := range findings {
		switch f.Severity {
		case "critical":
			report.Critical++
		case "high":
			report.High++
		case "medium":
			report.Medium++
		case "low":
			report.Low++
		}

		report.Findings = append(report.Findings, chineseFinding{
			File:           f.File,
			Line:           f.Line,
			Column:         f.Column,
			RuleID:         f.RuleID,
			RuleName:       f.RuleName,
			Severity:       f.Severity,
			SeverityCN:     severityToChinese(f.Severity),
			Category:       f.Category,
			CategoryCN:     categoryToChinese(f.Category),
			MatchedContent: f.MatchedContent,
			Description:    f.Description,
			DescriptionCN:  getChineseDescription(f),
		})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

// severityToChinese converts severity level to Chinese.
func severityToChinese(severity string) string {
	switch severity {
	case "critical":
		return "严重"
	case "high":
		return "高危"
	case "medium":
		return "中危"
	case "low":
		return "低危"
	default:
		return "未知"
	}
}

// categoryToChinese converts category to Chinese.
func categoryToChinese(category string) string {
	switch category {
	case "generic":
		return "通用规则"
	case "china":
		return "国内云服务商"
	default:
		return category
	}
}

// getChineseDescription returns a Chinese description for a finding.
func getChineseDescription(f rules.Finding) string {
	descMap := map[string]string{
		"GEN-001": "检测到 AWS Access Key ID（以 AKIA 开头）",
		"GEN-002": "检测到 AWS Secret Access Key（40字符 Base64 密钥）",
		"GEN-003": "检测到 Google API Key（以 AIza 开头）",
		"GEN-004": "检测到 Google OAuth 访问令牌",
		"GEN-005": "检测到 Azure 客户端密钥",
		"GEN-006": "检测到 Azure 存储/数据库连接字符串",
		"GEN-007": "检测到 Stripe 密钥 API Key（生产环境）",
		"GEN-008": "检测到 Stripe 可发布密钥（生产环境）",
		"GEN-009": "检测到 GitHub 个人访问令牌",
		"GEN-010": "检测到 GitLab 个人访问令牌",
		"GEN-011": "检测到 Slack 机器人/用户令牌",
		"GEN-012": "检测到 Slack Webhook URL",
		"GEN-013": "检测到 Twilio API Key",
		"GEN-014": "检测到 SendGrid API Key",
		"GEN-015": "检测到 Mailgun API Key",
		"GEN-016": "检测到 PostgreSQL 数据库连接字符串（含凭据）",
		"GEN-017": "检测到 MongoDB 数据库连接 URI（含凭据）",
		"GEN-018": "检测到 Redis 连接 URL（含密码）",
		"GEN-019": "检测到 RSA 私钥",
		"GEN-020": "检测到 EC 私钥",
		"GEN-021": "检测到 OpenSSH 私钥",
		"GEN-022": "检测到 DSA 私钥",
		"GEN-023": "检测到 PKCS8 私钥",
		"GEN-024": "检测到 JSON Web Token (JWT)",
		"CN-001":  "检测到阿里云 AccessKey ID（以 LTAI 开头）",
		"CN-002":  "检测到阿里云 AccessKey Secret",
		"CN-003":  "检测到腾讯云 SecretId（以 AKID 开头）",
		"CN-004":  "检测到腾讯云 SecretKey",
		"CN-005":  "检测到华为云 Access Key ID",
		"CN-006":  "检测到华为云 Secret Access Key",
		"CN-007":  "检测到百度云 API Key",
		"CN-008":  "检测到百度云 Secret Key",
		"CN-009":  "检测到又拍云 Operator",
		"CN-010":  "检测到又拍云 Password",
		"CN-011":  "检测到七牛云 AccessKey",
		"CN-012":  "检测到七牛云 SecretKey",
		"CN-013":  "检测到新浪云 Access Key",
		"CN-014":  "检测到京东云 Access Key",
		"CN-015":  "检测到京东云 Secret Key",
		"CN-016":  "检测到金山云 Access Key",
		"CN-017":  "检测到金山云 Secret Key",
		"CN-018":  "检测到 MiniMax API Key",
	}

	if desc, ok := descMap[f.RuleID]; ok {
		return desc
	}
	return "检测到潜在密钥泄露"
}
