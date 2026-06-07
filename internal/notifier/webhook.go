// Package notifier provides webhook notification capabilities for SecretShield-CLI.
// It supports sending scan results to Feishu (Lark) and DingTalk webhooks.
package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/secretshield/cli/internal/rules"
)

// Notifier sends scan results via webhook.
type Notifier struct {
	WebhookURL string
	HTTPClient *http.Client
}

// NewNotifier creates a new Notifier with the given webhook URL.
func NewNotifier(webhookURL string) *Notifier {
	return &Notifier{
		WebhookURL: webhookURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Notify sends the scan results to the configured webhook.
// It auto-detects the webhook type (Feishu or DingTalk) based on the URL.
func (n *Notifier) Notify(findings []rules.Finding) error {
	if n.WebhookURL == "" {
		return nil
	}

	switch {
	case strings.Contains(n.WebhookURL, "open.feishu.cn") || strings.Contains(n.WebhookURL, "open.larksuite.com"):
		return n.sendFeishu(findings)
	case strings.Contains(n.WebhookURL, "oapi.dingtalk.com"):
		return n.sendDingTalk(findings)
	default:
		return n.sendGeneric(findings)
	}
}

// sendFeishu sends results to a Feishu (Lark) robot webhook.
func (n *Notifier) sendFeishu(findings []rules.Finding) error {
	// Build the text content
	content := buildNotificationText(findings)

	payload := feishuPayload{
		MsgType: "interactive",
		Card: feishuCard{
			Header: feishuHeader{
				Title: feishuText{
					Tag:     "plain_text",
					Content: "SecretShield 密钥泄露扫描报告",
				},
				Template: "red",
			},
			Elements: []interface{}{
				feishuMarkdown{
					Tag:     "markdown",
					Content: content,
				},
				feishuNote{
					Elements: []feishuNoteElement{
						{
							Tag:     "plain_text",
							Content: fmt.Sprintf("共发现 %d 个潜在密钥泄露 | SecretShield-CLI v1.0.0", len(findings)),
						},
					},
				},
			},
		},
	}

	return n.sendWebhook(payload)
}

// sendDingTalk sends results to a DingTalk robot webhook.
func (n *Notifier) sendDingTalk(findings []rules.Finding) error {
	content := buildNotificationText(findings)

	payload := dingTalkPayload{
		MsgType: "markdown",
		Markdown: dingTalkMarkdown{
			Title: "SecretShield 密钥泄露扫描报告",
			Text:  content,
		},
	}

	return n.sendWebhook(payload)
}

// sendGeneric sends a generic JSON payload to the webhook.
func (n *Notifier) sendGeneric(findings []rules.Finding) error {
	payload := map[string]interface{}{
		"tool":    "SecretShield-CLI",
		"version": "1.0.0",
		"total":   len(findings),
		"findings": findings,
	}

	return n.sendWebhook(payload)
}

// sendWebhook sends a JSON payload to the webhook URL.
func (n *Notifier) sendWebhook(payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	resp, err := n.HTTPClient.Post(n.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status code: %d", resp.StatusCode)
	}

	return nil
}

// buildNotificationText builds a markdown text summary of findings.
func buildNotificationText(findings []rules.Finding) string {
	if len(findings) == 0 {
		return "**扫描完成：未发现密钥泄露** \n\n所有文件检查通过，代码安全。"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**扫描完成：发现 %d 个潜在密钥泄露**\n\n", len(findings)))

	critical := 0
	high := 0
	medium := 0
	low := 0

	for _, f := range findings {
		switch f.Severity {
		case "critical":
			critical++
		case "high":
			high++
		case "medium":
			medium++
		case "low":
			low++
		}
	}

	sb.WriteString(fmt.Sprintf("- 严重: **%d**\n", critical))
	sb.WriteString(fmt.Sprintf("- 高危: **%d**\n", high))
	sb.WriteString(fmt.Sprintf("- 中危: **%d**\n", medium))
	sb.WriteString(fmt.Sprintf("- 低危: **%d**\n\n", low))

	// Show top 10 findings
	maxShow := 10
	if len(findings) < maxShow {
		maxShow = len(findings)
	}

	sb.WriteString("**详细发现：**\n")
	for i := 0; i < maxShow; i++ {
		f := findings[i]
		sb.WriteString(fmt.Sprintf("%d. **[%s]** `%s` - %s:%d\n",
			i+1, f.RuleID, f.RuleName, f.File, f.Line))
	}

	if len(findings) > maxShow {
		sb.WriteString(fmt.Sprintf("\n... 还有 %d 个发现未显示\n", len(findings)-maxShow))
	}

	return sb.String()
}

// Feishu (Lark) webhook payload structures.
type feishuPayload struct {
	MsgType string     `json:"msg_type"`
	Card    feishuCard `json:"card"`
}

type feishuCard struct {
	Header   feishuHeader    `json:"header"`
	Elements []interface{}   `json:"elements"`
}

type feishuHeader struct {
	Title    feishuText `json:"title"`
	Template string    `json:"template"`
}

type feishuText struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type feishuMarkdown struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type feishuNote struct {
	Elements []feishuNoteElement `json:"elements"`
}

type feishuNoteElement struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

// DingTalk webhook payload structures.
type dingTalkPayload struct {
	MsgType  string            `json:"msgtype"`
	Markdown dingTalkMarkdown  `json:"markdown"`
}

type dingTalkMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
