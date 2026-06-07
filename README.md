<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green.svg?style=flat-square" alt="License">
  <img src="https://img.shields.io/badge/Version-v1.0.0-blue.svg?style=flat-square" alt="Version">
  <img src="https://img.shields.io/badge/Cobra-CLI-EF5030?style=flat-square&logo=github" alt="Cobra CLI">
</p>

<p align="center">
  <a href="#简体中文">简体中文</a> |
  <a href="#繁體中文">繁體中文</a> |
  <a href="#english">English</a>
</p>

---

<h1 align="center">🛡️ SecretShield-CLI</h1>

<p align="center">
  <strong>轻量级终端密钥泄露扫描引擎</strong><br>
  专为国内开发者打造，42+ 内置检测规则，零依赖 Go 二进制分发
</p>

---

<a id="简体中文"></a>

## 🎉 项目介绍

**SecretShield-CLI** 是一款轻量级的终端密钥泄露扫描工具，旨在帮助开发者在代码提交前、CI/CD 流水线中快速检测潜在的密钥和凭据泄露。它支持扫描本地目录、单个文件以及 Git 仓库（含完整提交历史）。

### 核心价值

- **安全左移**：在代码进入仓库之前就发现密钥泄露，而非事后补救
- **降低风险**：覆盖国内外主流云服务商和开发平台的密钥格式
- **提升效率**：一键扫描，秒级输出结果，无缝融入开发工作流

### 自研差异化亮点

与 [TruffleHog](https://github.com/trufflesecurity/trufflehog)、[Gitleaks](https://github.com/gitleaks/gitleaks) 等主流工具相比，SecretShield-CLI 具有以下独特优势：

| 特性 | SecretShield-CLI | TruffleHog | Gitleaks |
|------|:-:|:-:|:-:|
| 国内云服务商原生支持 | ✅ | ❌ | ❌ |
| 中文报告输出 | ✅ | ❌ | ❌ |
| 飞书/钉钉 Webhook 告警 | ✅ | ❌ | ❌ |
| 零依赖单二进制分发 | ✅ | ❌ (Python) | ✅ |
| 自定义规则文件 | ✅ | ✅ | ✅ |
| SARIF 格式输出 | ✅ | ✅ | ✅ |

### 灵感来源

密钥泄露是当前云原生时代最常见的安全隐患之一。根据 GitHub 安全报告，每年有数百万条密钥在公开代码仓库中被泄露。现有开源工具大多以国际云服务商（AWS、GCP、Azure）为核心，对国内云服务商（阿里云、腾讯云、华为云等）的覆盖严重不足。SecretShield-CLI 正是在这一背景下诞生的——**让国内开发者也能拥有贴合自身技术栈的密钥泄露检测工具**。

---

## ✨ 核心特性

- 🔍 **42+ 内置检测规则**：涵盖 24 条通用规则（AWS/GCP/Azure/GitHub/GitLab/Stripe 等）和 18 条国内云服务商规则（阿里云/腾讯云/华为云/百度云/七牛云/京东云/金山云/又拍云/新浪云/MiniMax 等）
- 🇨🇳 **国内云服务商原生支持**：深度覆盖阿里云 AccessKey、腾讯云 SecretId/Key、华为云 AK/SK、百度云、七牛云、京东云、金山云等主流国内平台
- 📜 **Git 仓库历史扫描**：通过 `--git` 参数可扫描 Git 仓库的完整提交历史和 Diff 内容，发现已经被删除但仍存在于历史记录中的密钥
- 📊 **4 种输出格式**：支持 **JSON**（结构化数据）、**SARIF**（CI/CD 集成）、**Table**（终端表格）和 **中文报告**（全中文字段名，适合国内团队）
- 📢 **飞书/钉钉 Webhook 告警**：扫描完成后自动推送结果到飞书群机器人或钉钉群机器人，支持交互式卡片和 Markdown 格式
- ⚙️ **自定义规则 YAML 配置**：支持通过自定义规则文件扩展检测能力，灵活适配企业内部系统
- 💻 **跨平台单二进制分发**：编译后为单个可执行文件，无任何外部运行时依赖（除 Cobra CLI 框架外），支持 Linux / macOS / Windows
- 🔄 **CI/CD 集成友好**：输出 SARIF 格式可直接对接 GitHub Code Scanning、GitLab SAST 等平台；扫描发现问题时以非零退出码返回
- 🚀 **零运行时依赖**：无需安装 Python、Node.js 等运行时环境，下载即用

---

## 🚀 快速开始

### 环境要求

- **Go 1.21+**（仅从源码编译时需要）
- **Git**（使用 `--git` 历史扫描时需要）

### 安装方式

#### 方式一：go install（推荐 Go 用户）

```bash
go install github.com/secretshield/cli@latest
```

安装后，`secretshield` 命令将自动加入你的 `GOPATH/bin`（请确保该目录已在 `PATH` 中）。

#### 方式二：下载预编译二进制

从 [GitHub Releases](https://github.com/gitstq/SecretShield-CLI/releases) 页面下载对应平台的预编译二进制文件：

```bash
# Linux (amd64)
wget https://github.com/gitstq/SecretShield-CLI/releases/download/v1.0.0/secretshield-linux-amd64 -O secretshield
chmod +x secretshield
sudo mv secretshield /usr/local/bin/

# macOS (amd64)
wget https://github.com/gitstq/SecretShield-CLI/releases/download/v1.0.0/secretshield-darwin-amd64 -O secretshield
chmod +x secretshield
sudo mv secretshield /usr/local/bin/

# Windows (amd64)
# 下载 secretshield-windows-amd64.exe 并添加到 PATH
```

#### 方式三：从源码编译

```bash
git clone https://github.com/gitstq/SecretShield-CLI.git
cd SecretShield-CLI
go build -o bin/secretshield .
```

### 验证安装

```bash
secretshield version
# 输出：SecretShield-CLI v1.0.0
```

---

## 📖 详细使用指南

### scan 命令

`scan` 是 SecretShield-CLI 的核心命令，用于扫描目录、Git 仓库或单个文件中的密钥泄露。

#### 完整参数说明

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--git` | `-g` | `false` | 启用 Git 仓库历史扫描模式 |
| `--file` | `-f` | `false` | 扫描单个文件而非目录 |
| `--output` | `-o` | `table` | 输出格式：`json`、`sarif`、`table`、`chinese` |
| `--report` | `-r` | `""` | 将报告写入指定文件（默认输出到 stdout） |
| `--rules` | | `""` | 自定义规则文件路径 |
| `--exclude` | `-e` | `""` | 逗号分隔的排除目录列表 |
| `--severity` | `-s` | `""` | 逗号分隔的严重级别过滤：`critical`、`high`、`medium`、`low` |
| `--webhook` | `-w` | `""` | Webhook URL（支持飞书/钉钉自动识别） |

> **注意**：`--file` 和 `--git` 参数不能同时使用。

### rules 命令

列出所有内置检测规则，支持按类别筛选。

```bash
# 列出所有规则
secretshield rules

# 仅列出国内云服务商规则
secretshield rules --category china

# 仅列出通用/国际规则
secretshield rules --category generic
```

### version 命令

查看当前版本信息。

```bash
secretshield version
```

### 典型使用场景

#### 场景一：扫描当前目录

```bash
secretshield scan
```

扫描当前工作目录下的所有文件，使用默认的 `table` 格式输出结果。

#### 场景二：扫描指定项目目录

```bash
secretshield scan /path/to/your/project
```

#### 场景三：扫描 Git 仓库历史

```bash
secretshield scan --git /path/to/your/repo
```

启用 Git 历史扫描，会检查所有历史提交和 Diff 中的密钥泄露。

#### 场景四：输出 SARIF 格式用于 CI/CD

```bash
secretshield scan --output sarif --report results.sarif /path/to/project
```

生成的 SARIF 文件可直接上传至 GitHub Code Scanning 或 GitLab SAST。

#### 场景五：仅检测高危和严重级别

```bash
secretshield scan --severity high,critical /path/to/project
```

#### 场景六：排除特定目录

```bash
secretshield scan --exclude vendor,node_modules,dist,build /path/to/project
```

#### 场景七：扫描单个文件

```bash
secretshield scan --file config/production.yaml
```

#### 场景八：中文报告 + 飞书告警

```bash
secretshield scan --output chinese --report report.json --webhook "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxx" /path/to/project
```

#### 场景九：使用自定义规则

```bash
secretshield scan --rules ./my-custom-rules.txt /path/to/project
```

#### 场景十：CI/CD 中使用（GitHub Actions）

```yaml
- name: Secret Leak Scan
  run: |
    secretshield scan --output sarif --report results.sarif --severity high,critical .
  # 扫描到问题时 secretshield 会返回非零退出码，自动阻断流水线
```

### 自定义规则配置

自定义规则文件采用文本格式，每行一条规则，字段以 `|` 分隔：

```text
# 格式：规则ID|规则名称|正则表达式|严重级别|类别
# 严重级别：critical / high / medium / low
# 类别：generic / china

# 示例：检测内部 API Token
CUSTOM-001|Internal API Token|internal_token\s*[=:]\s*["']?[A-Za-z0-9]{32,}["']?|high|generic

# 示例：检测内部数据库密码
CUSTOM-002|Internal DB Password|db_password\s*[=:]\s*["']?[^\s"']{8,}["']?|critical|china
```

### Webhook 配置示例

#### 飞书群机器人

```bash
secretshield scan --webhook "https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-token" /path/to/project
```

飞书 Webhook 会自动发送**交互式卡片消息**，包含扫描结果摘要和详细发现列表。

#### 钉钉群机器人

```bash
secretshield scan --webhook "https://oapi.dingtalk.com/robot/send?access_token=your-access-token" /path/to/project
```

钉钉 Webhook 会自动发送 **Markdown 格式消息**，包含扫描结果摘要。

#### 通用 Webhook

如果 URL 不匹配飞书或钉钉的域名，将发送通用 JSON 格式的数据：

```json
{
  "tool": "SecretShield-CLI",
  "version": "1.0.0",
  "total": 3,
  "findings": [...]
}
```

---

## 💡 设计思路与迭代规划

### 设计理念

1. **简洁至上**：单二进制，零配置即可使用，降低使用门槛
2. **安全左移**：面向开发阶段的密钥检测，而非事后审计
3. **本土化优先**：原生支持国内云服务商和中文报告，贴合国内开发者需求
4. **可扩展性**：支持自定义规则和 Webhook，适配不同团队的安全工作流

### 技术选型原因

| 技术选择 | 原因 |
|----------|------|
| **Go 语言** | 编译为单二进制、跨平台、启动速度快、部署简单 |
| **Cobra CLI** | Go 生态中最成熟的 CLI 框架，提供参数解析、帮助文档生成等能力 |
| **正则表达式匹配** | 密钥格式检测的核心手段，配合关键词预过滤提升性能 |
| **关键词预过滤** | 先通过字符串包含检查快速排除不相关行，再执行正则匹配，大幅提升扫描速度 |

### 后续功能迭代计划

#### v1.1.0 — 增强检测能力

- 新增更多国内云服务商规则（如火山引擎、UCloud、青云等）
- 支持忽略文件 `.secretshieldignore`（类似 `.gitignore` 语法）
- 增加对 `.env` 文件的专门解析和检测

#### v1.2.0 — 团队协作

- 支持企业微信 Webhook 告警
- 增加扫描结果基线（baseline）功能，仅报告新增泄露
- 支持配置文件 `secretshield.yaml` 统一管理扫描参数

#### v2.0.0 — 架构升级

- 引入规则引擎插件化架构，支持外部规则包
- 增加 HTML 可视化报告输出
- 提供 Server 模式，支持 API 调用和定时扫描任务
- 支持增量扫描，仅扫描变更文件

### 社区贡献方向

我们欢迎以下方向的贡献：

- 🆕 新增检测规则（尤其是国内云服务商和开发平台）
- 🌐 多语言报告模板（日语、韩语等）
- 📦 新增 Webhook 通道（企业微信、Slack 等）
- 📖 文档改进和翻译
- 🐛 Bug 修复和性能优化

---

## 📦 打包与部署指南

### 从源码编译

```bash
git clone https://github.com/gitstq/SecretShield-CLI.git
cd SecretShield-CLI
go build -o bin/secretshield .
```

### 交叉编译

使用 Makefile 提供的 `release` 目标进行交叉编译：

```bash
make release
```

或手动指定目标平台：

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -o secretshield-linux-amd64 .

# Linux arm64
GOOS=linux GOARCH=arm64 go build -o secretshield-linux-arm64 .

# macOS amd64
GOOS=darwin GOARCH=amd64 go build -o secretshield-darwin-amd64 .

# macOS arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o secretshield-darwin-arm64 .

# Windows amd64
GOOS=windows GOARCH=amd64 go build -o secretshield-windows-amd64.exe .
```

### Release 下载

预编译二进制文件可在 [GitHub Releases](https://github.com/gitstq/SecretShield-CLI/releases) 页面获取。

### 兼容环境说明

| 平台 | 架构 | 支持情况 |
|------|------|----------|
| Linux | amd64 | ✅ |
| Linux | arm64 | ✅ |
| macOS | amd64 | ✅ |
| macOS | arm64 (Apple Silicon) | ✅ |
| Windows | amd64 | ✅ |

> **最低要求**：Go 1.21+ 编译环境（仅编译时需要，运行时无 Go 依赖）

---

## 🤝 贡献指南

我们非常欢迎社区贡献！以下是参与贡献的基本流程。

### PR 提交规范

1. **Fork** 本仓库并创建你的特性分支：`git checkout -b feature/your-feature-name`
2. **编写代码**，确保通过所有测试：`go test ./... -v`
3. **提交信息** 遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：
   - `feat: 新增火山引擎检测规则`
   - `fix: 修复正则表达式误报问题`
   - `docs: 更新 README 使用说明`
   - `refactor: 重构扫描引擎性能优化`
4. **推送** 到你的 Fork：`git push origin feature/your-feature-name`
5. 创建 **Pull Request**，详细描述变更内容

### Issue 反馈规则

提交 Issue 时，请包含以下信息：

- **问题描述**：清晰描述遇到的问题或期望的功能
- **复现步骤**：如果是 Bug，请提供完整的复现步骤
- **环境信息**：操作系统、Go 版本、SecretShield-CLI 版本
- **日志输出**：如有相关日志或截图，请一并提供

### 开发环境搭建

```bash
# 克隆仓库
git clone https://github.com/gitstq/SecretShield-CLI.git
cd SecretShield-CLI

# 安装依赖
go mod download

# 运行测试
go test ./... -v

# 代码检查
go vet ./...

# 本地构建
go build -o bin/secretshield .

# 运行
./bin/secretshield version
```

---

## 📄 开源协议说明

本项目基于 **MIT License** 开源。

```
MIT License

Copyright (c) 2024 SecretShield Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

这意味着你可以自由地使用、复制、修改、合并、发布、分发、再授权和/或销售本软件，唯一的要求是保留版权声明和许可声明。

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/gitstq/SecretShield-CLI">SecretShield Contributors</a>
</p>

---

<a id="繁體中文"></a>

## 🎉 專案介紹

**SecretShield-CLI** 是一款輕量級的終端密鑰洩漏掃描工具，旨在協助開發者在程式碼提交前、CI/CD 管線中快速偵測潛在的密鑰與憑證洩漏。它支援掃描本地目錄、單一檔案以及 Git 倉庫（含完整提交歷史）。

### 核心價值

- **安全左移**：在程式碼進入倉庫之前就發現密鑰洩漏，而非事後補救
- **降低風險**：涵蓋國內外主流雲端服務商和開發平台的密鑰格式
- **提升效率**：一鍵掃描，秒級輸出結果，無縫融入開發工作流

### 自研差異化亮點

與 [TruffleHog](https://github.com/trufflesecurity/trufflehog)、[Gitleaks](https://github.com/gitleaks/gitleaks) 等主流工具相比，SecretShield-CLI 具有以下獨特優勢：

| 特性 | SecretShield-CLI | TruffleHog | Gitleaks |
|------|:-:|:-:|:-:|
| 國內雲端服務商原生支援 | ✅ | ❌ | ❌ |
| 中文報告輸出 | ✅ | ❌ | ❌ |
| 飛書/釘釘 Webhook 告警 | ✅ | ❌ | ❌ |
| 零依賴單二進制分發 | ✅ | ❌ (Python) | ✅ |
| 自訂規則檔案 | ✅ | ✅ | ✅ |
| SARIF 格式輸出 | ✅ | ✅ | ✅ |

### 靈感來源

密鑰洩漏是當前雲原生時代最常見的安全隱患之一。根據 GitHub 安全報告，每年有數百萬條密鑰在公開程式碼倉庫中被洩漏。現有開源工具大多以國際雲端服務商（AWS、GCP、Azure）為核心，對國內雲端服務商（阿里雲、騰訊雲、華為雲等）的覆蓋嚴重不足。SecretShield-CLI 正是在這一背景下誕生的——**讓國內開發者也能擁有貼合自身技術棧的密鑰洩漏偵測工具**。

---

## ✨ 核心特性

- 🔍 **42+ 內建偵測規則**：涵蓋 24 條通用規則（AWS/GCP/Azure/GitHub/GitLab/Stripe 等）和 18 條國內雲端服務商規則（阿里雲/騰訊雲/華為雲/百度雲/七牛雲/京東雲/金山雲/又拍雲/新浪雲/MiniMax 等）
- 🇨🇳 **國內雲端服務商原生支援**：深度覆蓋阿里雲 AccessKey、騰訊雲 SecretId/Key、華為雲 AK/SK、百度雲、七牛雲、京東雲、金山雲等主流國內平台
- 📜 **Git 倉庫歷史掃描**：透過 `--git` 參數可掃描 Git 倉庫的完整提交歷史和 Diff 內容，發現已經被刪除但仍存在於歷史記錄中的密鑰
- 📊 **4 種輸出格式**：支援 **JSON**（結構化資料）、**SARIF**（CI/CD 整合）、**Table**（終端表格）和 **中文報告**（全中文字段名，適合國內團隊）
- 📢 **飛書/釘釘 Webhook 告警**：掃描完成後自動推送結果到飛書群機器人或釘釘群機器人，支援互動式卡片和 Markdown 格式
- ⚙️ **自訂規則配置**：支援透過自訂規則檔案擴展偵測能力，靈活適配企業內部系統
- 💻 **跨平台單二進制分發**：編譯後為單一可執行檔，無任何外部執行期依賴（除 Cobra CLI 框架外），支援 Linux / macOS / Windows
- 🔄 **CI/CD 整合友善**：輸出 SARIF 格式可直接對接 GitHub Code Scanning、GitLab SAST 等平台；掃描發現問題時以非零退出碼返回
- 🚀 **零執行期依賴**：無需安裝 Python、Node.js 等執行期環境，下載即用

---

## 🚀 快速開始

### 環境要求

- **Go 1.21+**（僅從原始碼編譯時需要）
- **Git**（使用 `--git` 歷史掃描時需要）

### 安裝方式

#### 方式一：go install（推薦 Go 使用者）

```bash
go install github.com/secretshield/cli@latest
```

安裝後，`secretshield` 指令將自動加入你的 `GOPATH/bin`（請確保該目錄已在 `PATH` 中）。

#### 方式二：下載預編譯二進制檔

從 [GitHub Releases](https://github.com/gitstq/SecretShield-CLI/releases) 頁面下載對應平台的預編譯二進制檔：

```bash
# Linux (amd64)
wget https://github.com/gitstq/SecretShield-CLI/releases/download/v1.0.0/secretshield-linux-amd64 -O secretshield
chmod +x secretshield
sudo mv secretshield /usr/local/bin/

# macOS (amd64)
wget https://github.com/gitstq/SecretShield-CLI/releases/download/v1.0.0/secretshield-darwin-amd64 -O secretshield
chmod +x secretshield
sudo mv secretshield /usr/local/bin/

# Windows (amd64)
# 下載 secretshield-windows-amd64.exe 並新增至 PATH
```

#### 方式三：從原始碼編譯

```bash
git clone https://github.com/gitstq/SecretShield-CLI.git
cd SecretShield-CLI
go build -o bin/secretshield .
```

### 驗證安裝

```bash
secretshield version
# 輸出：SecretShield-CLI v1.0.0
```

---

## 📖 詳細使用指南

### scan 指令

`scan` 是 SecretShield-CLI 的核心指令，用於掃描目錄、Git 倉庫或單一檔案中的密鑰洩漏。

#### 完整參數說明

| 參數 | 簡寫 | 預設值 | 說明 |
|------|------|--------|------|
| `--git` | `-g` | `false` | 啟用 Git 倉庫歷史掃描模式 |
| `--file` | `-f` | `false` | 掃描單一檔案而非目錄 |
| `--output` | `-o` | `table` | 輸出格式：`json`、`sarif`、`table`、`chinese` |
| `--report` | `-r` | `""` | 將報告寫入指定檔案（預設輸出到 stdout） |
| `--rules` | | `""` | 自訂規則檔案路徑 |
| `--exclude` | `-e` | `""` | 逗號分隔的排除目錄列表 |
| `--severity` | `-s` | `""` | 逗號分隔的嚴重級別過濾：`critical`、`high`、`medium`、`low` |
| `--webhook` | `-w` | `""` | Webhook URL（支援飛書/釘釘自動識別） |

> **注意**：`--file` 和 `--git` 參數不能同時使用。

### rules 指令

列出所有內建偵測規則，支援按類別篩選。

```bash
# 列出所有規則
secretshield rules

# 僅列出國內雲端服務商規則
secretshield rules --category china

# 僅列出通用/國際規則
secretshield rules --category generic
```

### version 指令

查看目前版本資訊。

```bash
secretshield version
```

### 典型使用場景

#### 場景一：掃描目前目錄

```bash
secretshield scan
```

掃描目前工作目錄下的所有檔案，使用預設的 `table` 格式輸出結果。

#### 場景二：掃描指定專案目錄

```bash
secretshield scan /path/to/your/project
```

#### 場景三：掃描 Git 倉庫歷史

```bash
secretshield scan --git /path/to/your/repo
```

啟用 Git 歷史掃描，會檢查所有歷史提交和 Diff 中的密鑰洩漏。

#### 場景四：輸出 SARIF 格式用於 CI/CD

```bash
secretshield scan --output sarif --report results.sarif /path/to/project
```

產生的 SARIF 檔案可直接上傳至 GitHub Code Scanning 或 GitLab SAST。

#### 場景五：僅偵測高危和嚴重級別

```bash
secretshield scan --severity high,critical /path/to/project
```

#### 場景六：排除特定目錄

```bash
secretshield scan --exclude vendor,node_modules,dist,build /path/to/project
```

#### 場景七：掃描單一檔案

```bash
secretshield scan --file config/production.yaml
```

#### 場景八：中文報告 + 飛書告警

```bash
secretshield scan --output chinese --report report.json --webhook "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxx" /path/to/project
```

#### 場景九：使用自訂規則

```bash
secretshield scan --rules ./my-custom-rules.txt /path/to/project
```

#### 場景十：CI/CD 中使用（GitHub Actions）

```yaml
- name: Secret Leak Scan
  run: |
    secretshield scan --output sarif --report results.sarif --severity high,critical .
  # 掃描到問題時 secretshield 會回傳非零退出碼，自動阻斷管線
```

### 自訂規則配置

自訂規則檔案採用文字格式，每行一條規則，欄位以 `|` 分隔：

```text
# 格式：規則ID|規則名稱|正規表示式|嚴重級別|類別
# 嚴重級別：critical / high / medium / low
# 類別：generic / china

# 範例：偵測內部 API Token
CUSTOM-001|Internal API Token|internal_token\s*[=:]\s*["']?[A-Za-z0-9]{32,}["']?|high|generic

# 範例：偵測內部資料庫密碼
CUSTOM-002|Internal DB Password|db_password\s*[=:]\s*["']?[^\s"']{8,}["']?|critical|china
```

### Webhook 配置範例

#### 飛書群機器人

```bash
secretshield scan --webhook "https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-token" /path/to/project
```

飛書 Webhook 會自動發送**互動式卡片訊息**，包含掃描結果摘要和詳細發現列表。

#### 釘釘群機器人

```bash
secretshield scan --webhook "https://oapi.dingtalk.com/robot/send?access_token=your-access-token" /path/to/project
```

釘釘 Webhook 會自動發送 **Markdown 格式訊息**，包含掃描結果摘要。

#### 通用 Webhook

如果 URL 不符合飛書或釘釘的網域，將發送通用 JSON 格式的資料：

```json
{
  "tool": "SecretShield-CLI",
  "version": "1.0.0",
  "total": 3,
  "findings": [...]
}
```

---

## 💡 設計思路與迭代規劃

### 設計理念

1. **簡潔至上**：單二進制，零配置即可使用，降低使用門檻
2. **安全左移**：面向開發階段的密鑰偵測，而非事後稽核
3. **本土化優先**：原生支援國內雲端服務商和中文報告，貼合國內開發者需求
4. **可擴展性**：支援自訂規則和 Webhook，適配不同團隊的安全工作流

### 技術選型原因

| 技術選擇 | 原因 |
|----------|------|
| **Go 語言** | 編譯為單二進制、跨平台、啟動速度快、部署簡單 |
| **Cobra CLI** | Go 生態中最成熟的 CLI 框架，提供參數解析、說明文件生成等能力 |
| **正規表示式匹配** | 密鑰格式偵測的核心手段，配合關鍵字預過濾提升效能 |
| **關鍵字預過濾** | 先透過字串包含檢查快速排除不相關行，再執行正規表示式匹配，大幅提升掃描速度 |

### 後續功能迭代計畫

#### v1.1.0 — 增強偵測能力

- 新增更多國內雲端服務商規則（如火山引擎、UCloud、青雲等）
- 支援忽略檔案 `.secretshieldignore`（類似 `.gitignore` 語法）
- 增加對 `.env` 檔案的專門解析和偵測

#### v1.2.0 — 團隊協作

- 支援企業微信 Webhook 告警
- 增加掃描結果基線（baseline）功能，僅報告新增洩漏
- 支援設定檔 `secretshield.yaml` 統一管理掃描參數

#### v2.0.0 — 架構升級

- 引入規則引擎插件化架構，支援外部規則包
- 增加 HTML 視覺化報告輸出
- 提供 Server 模式，支援 API 呼叫和定時掃描任務
- 支援增量掃描，僅掃描變更檔案

### 社群貢獻方向

我們歡迎以下方向的貢獻：

- 🆕 新增偵測規則（尤其是國內雲端服務商和開發平台）
- 🌐 多語言報告範本（日語、韓語等）
- 📦 新增 Webhook 通道（企業微信、Slack 等）
- 📖 文件改進和翻譯
- 🐛 Bug 修復和效能最佳化

---

## 📦 打包與部署指南

### 從原始碼編譯

```bash
git clone https://github.com/gitstq/SecretShield-CLI.git
cd SecretShield-CLI
go build -o bin/secretshield .
```

### 交叉編譯

使用 Makefile 提供的 `release` 目標進行交叉編譯：

```bash
make release
```

或手動指定目標平台：

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -o secretshield-linux-amd64 .

# Linux arm64
GOOS=linux GOARCH=arm64 go build -o secretshield-linux-arm64 .

# macOS amd64
GOOS=darwin GOARCH=amd64 go build -o secretshield-darwin-amd64 .

# macOS arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o secretshield-darwin-arm64 .

# Windows amd64
GOOS=windows GOARCH=amd64 go build -o secretshield-windows-amd64.exe .
```

### Release 下載

預編譯二進制檔可在 [GitHub Releases](https://github.com/gitstq/SecretShield-CLI/releases) 頁面取得。

### 相容環境說明

| 平台 | 架構 | 支援情況 |
|------|------|----------|
| Linux | amd64 | ✅ |
| Linux | arm64 | ✅ |
| macOS | amd64 | ✅ |
| macOS | arm64 (Apple Silicon) | ✅ |
| Windows | amd64 | ✅ |

> **最低要求**：Go 1.21+ 編譯環境（僅編譯時需要，執行期無 Go 依賴）

---

## 🤝 貢獻指南

我們非常歡迎社群貢獻！以下是參與貢獻的基本流程。

### PR 提交規範

1. **Fork** 本倉庫並建立你的特性分支：`git checkout -b feature/your-feature-name`
2. **撰寫程式碼**，確保通過所有測試：`go test ./... -v`
3. **提交訊息** 遵循 [Conventional Commits](https://www.conventionalcommits.org/) 規範：
   - `feat: 新增火山引擎偵測規則`
   - `fix: 修復正規表示式誤報問題`
   - `docs: 更新 README 使用說明`
   - `refactor: 重構掃描引擎效能最佳化`
4. **推送** 到你的 Fork：`git push origin feature/your-feature-name`
5. 建立 **Pull Request**，詳細描述變更內容

### Issue 回饋規則

提交 Issue 時，請包含以下資訊：

- **問題描述**：清晰描述遇到的問題或期望的功能
- **重現步驟**：如果是 Bug，請提供完整的重現步驟
- **環境資訊**：作業系統、Go 版本、SecretShield-CLI 版本
- **日誌輸出**：如有相關日誌或截圖，請一併提供

### 開發環境搭建

```bash
# 複製倉庫
git clone https://github.com/gitstq/SecretShield-CLI.git
cd SecretShield-CLI

# 安裝依賴
go mod download

# 執行測試
go test ./... -v

# 程式碼檢查
go vet ./...

# 本地建置
go build -o bin/secretshield .

# 執行
./bin/secretshield version
```

---

## 📄 開源協議說明

本專案基於 **MIT License** 開源。

```
MIT License

Copyright (c) 2024 SecretShield Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

這意味著你可以自由地使用、複製、修改、合併、發布、分發、再授權和/或銷售本軟體，唯一的要求是保留著作權聲明和許可聲明。

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/gitstq/SecretShield-CLI">SecretShield Contributors</a>
</p>

---

<a id="english"></a>

## 🎉 Introduction

**SecretShield-CLI** is a lightweight, terminal-based secret leak scanning engine designed to help developers quickly detect potential credential and secret leaks before code commits and within CI/CD pipelines. It supports scanning local directories, individual files, and Git repositories (including full commit history).

### Core Value

- **Shift Security Left**: Detect secret leaks before code enters the repository, rather than remediating after the fact
- **Risk Reduction**: Comprehensive coverage of both international and Chinese cloud provider credential formats
- **Developer Efficiency**: One-command scan with instant results, seamlessly integrating into developer workflows

### What Makes Us Different

Compared to mainstream tools like [TruffleHog](https://github.com/trufflesecurity/trufflehog) and [Gitleaks](https://github.com/gitleaks/gitleaks), SecretShield-CLI offers unique advantages:

| Feature | SecretShield-CLI | TruffleHog | Gitleaks |
|---------|:-:|:-:|:-:|
| Chinese cloud provider support | ✅ | ❌ | ❌ |
| Chinese report output | ✅ | ❌ | ❌ |
| Feishu/DingTalk webhook alerts | ✅ | ❌ | ❌ |
| Zero-dependency single binary | ✅ | ❌ (Python) | ✅ |
| Custom rule files | ✅ | ✅ | ✅ |
| SARIF format output | ✅ | ✅ | ✅ |

### Inspiration

Secret leaks are among the most common security vulnerabilities in the cloud-native era. According to GitHub security reports, millions of secrets are leaked in public repositories every year. Existing open-source tools primarily focus on international cloud providers (AWS, GCP, Azure), leaving Chinese cloud providers (Alibaba Cloud, Tencent Cloud, Huawei Cloud, etc.) severely under-covered. SecretShield-CLI was born from this gap -- **empowering Chinese developers with a secret scanning tool tailored to their own tech stack**.

---

## ✨ Key Features

- 🔍 **42+ Built-in Detection Rules**: Covering 24 generic rules (AWS/GCP/Azure/GitHub/GitLab/Stripe, etc.) and 18 Chinese cloud provider rules (Alibaba Cloud/Tencent Cloud/Huawei Cloud/Baidu Cloud/Qiniu Cloud/JD Cloud/Kingsoft Cloud/UpYun/Sina Cloud/MiniMax, etc.)
- 🇨🇳 **Native Chinese Cloud Provider Support**: Deep coverage of Alibaba Cloud AccessKey, Tencent Cloud SecretId/Key, Huawei Cloud AK/SK, Baidu Cloud, Qiniu Cloud, JD Cloud, Kingsoft Cloud, and other major Chinese platforms
- 📜 **Git Repository History Scanning**: Use the `--git` flag to scan the full commit history and diff content of Git repositories, discovering secrets that have been deleted but still exist in historical records
- 📊 **4 Output Formats**: Supports **JSON** (structured data), **SARIF** (CI/CD integration), **Table** (terminal table), and **Chinese Report** (fully localized field names for Chinese teams)
- 📢 **Feishu/DingTalk Webhook Alerts**: Automatically push scan results to Feishu or DingTalk group bots upon completion, with interactive card and Markdown format support
- ⚙️ **Custom Rule Configuration**: Extend detection capabilities through custom rule files, flexibly adapting to enterprise internal systems
- 💻 **Cross-Platform Single Binary Distribution**: Compiles to a single executable with no external runtime dependencies (except Cobra CLI framework), supporting Linux / macOS / Windows
- 🔄 **CI/CD Integration Friendly**: SARIF output integrates directly with GitHub Code Scanning, GitLab SAST, and other platforms; non-zero exit code on findings for pipeline gating
- 🚀 **Zero Runtime Dependencies**: No Python, Node.js, or other runtime environments required -- download and run

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** (only required when building from source)
- **Git** (required for `--git` history scanning)

### Installation

#### Option 1: go install (Recommended for Go users)

```bash
go install github.com/secretshield/cli@latest
```

After installation, the `secretshield` command will be available in your `GOPATH/bin` (ensure this directory is in your `PATH`).

#### Option 2: Download Pre-built Binary

Download the pre-built binary for your platform from the [GitHub Releases](https://github.com/gitstq/SecretShield-CLI/releases) page:

```bash
# Linux (amd64)
wget https://github.com/gitstq/SecretShield-CLI/releases/download/v1.0.0/secretshield-linux-amd64 -O secretshield
chmod +x secretshield
sudo mv secretshield /usr/local/bin/

# macOS (amd64)
wget https://github.com/gitstq/SecretShield-CLI/releases/download/v1.0.0/secretshield-darwin-amd64 -O secretshield
chmod +x secretshield
sudo mv secretshield /usr/local/bin/

# Windows (amd64)
# Download secretshield-windows-amd64.exe and add to PATH
```

#### Option 3: Build from Source

```bash
git clone https://github.com/gitstq/SecretShield-CLI.git
cd SecretShield-CLI
go build -o bin/secretshield .
```

### Verify Installation

```bash
secretshield version
# Output: SecretShield-CLI v1.0.0
```

---

## 📖 Detailed Usage Guide

### scan Command

The `scan` command is the core of SecretShield-CLI, used to scan directories, Git repositories, or individual files for leaked secrets.

#### Complete Parameter Reference

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--git` | `-g` | `false` | Enable Git repository history scanning mode |
| `--file` | `-f` | `false` | Scan a single file instead of a directory |
| `--output` | `-o` | `table` | Output format: `json`, `sarif`, `table`, `chinese` |
| `--report` | `-r` | `""` | Write report to file instead of stdout |
| `--rules` | | `""` | Path to custom rules file |
| `--exclude` | `-e` | `""` | Comma-separated list of directories to exclude |
| `--severity` | `-s` | `""` | Comma-separated severity filter: `critical`, `high`, `medium`, `low` |
| `--webhook` | `-w` | `""` | Webhook URL for notifications (auto-detects Feishu/DingTalk) |

> **Note**: The `--file` and `--git` flags cannot be used simultaneously.

### rules Command

List all built-in detection rules with optional category filtering.

```bash
# List all rules
secretshield rules

# List Chinese cloud provider rules only
secretshield rules --category china

# List generic/international rules only
secretshield rules --category generic
```

### version Command

Display the current version information.

```bash
secretshield version
```

### Common Usage Scenarios

#### Scenario 1: Scan the Current Directory

```bash
secretshield scan
```

Scans all files in the current working directory with the default `table` output format.

#### Scenario 2: Scan a Specific Project Directory

```bash
secretshield scan /path/to/your/project
```

#### Scenario 3: Scan Git Repository History

```bash
secretshield scan --git /path/to/your/repo
```

Enables Git history scanning, checking all historical commits and diffs for leaked secrets.

#### Scenario 4: Output SARIF Format for CI/CD

```bash
secretshield scan --output sarif --report results.sarif /path/to/project
```

The generated SARIF file can be uploaded directly to GitHub Code Scanning or GitLab SAST.

#### Scenario 5: Filter by High and Critical Severity Only

```bash
secretshield scan --severity high,critical /path/to/project
```

#### Scenario 6: Exclude Specific Directories

```bash
secretshield scan --exclude vendor,node_modules,dist,build /path/to/project
```

#### Scenario 7: Scan a Single File

```bash
secretshield scan --file config/production.yaml
```

#### Scenario 8: Chinese Report with Feishu Alert

```bash
secretshield scan --output chinese --report report.json --webhook "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxx" /path/to/project
```

#### Scenario 9: Use Custom Rules

```bash
secretshield scan --rules ./my-custom-rules.txt /path/to/project
```

#### Scenario 10: CI/CD Integration (GitHub Actions)

```yaml
- name: Secret Leak Scan
  run: |
    secretshield scan --output sarif --report results.sarif --severity high,critical .
  # secretshield returns non-zero exit code on findings, automatically gating the pipeline
```

### Custom Rule Configuration

Custom rule files use a plain text format with one rule per line, fields separated by `|`:

```text
# Format: rule_id|rule_name|regex_pattern|severity|category
# Severity levels: critical / high / medium / low
# Categories: generic / china

# Example: Detect internal API tokens
CUSTOM-001|Internal API Token|internal_token\s*[=:]\s*["']?[A-Za-z0-9]{32,}["']?|high|generic

# Example: Detect internal database passwords
CUSTOM-002|Internal DB Password|db_password\s*[=:]\s*["']?[^\s"']{8,}["']?|critical|china
```

### Webhook Configuration Examples

#### Feishu (Lark) Group Bot

```bash
secretshield scan --webhook "https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-token" /path/to/project
```

Feishu webhooks automatically send **interactive card messages** with a scan result summary and detailed findings list.

#### DingTalk Group Bot

```bash
secretshield scan --webhook "https://oapi.dingtalk.com/robot/send?access_token=your-access-token" /path/to/project
```

DingTalk webhooks automatically send **Markdown format messages** with a scan result summary.

#### Generic Webhook

If the URL does not match Feishu or DingTalk domains, a generic JSON payload is sent:

```json
{
  "tool": "SecretShield-CLI",
  "version": "1.0.0",
  "total": 3,
  "findings": [...]
}
```

---

## 💡 Design Philosophy & Roadmap

### Design Principles

1. **Simplicity First**: Single binary, zero configuration required out of the box, minimal learning curve
2. **Shift Security Left**: Developer-focused secret detection during the coding phase, not post-hoc auditing
3. **Localization Priority**: Native support for Chinese cloud providers and Chinese-language reports, tailored to Chinese developer needs
4. **Extensibility**: Custom rules and webhooks for adapting to different team security workflows

### Technical Choices

| Choice | Rationale |
|--------|-----------|
| **Go Language** | Compiles to a single binary, cross-platform, fast startup, simple deployment |
| **Cobra CLI** | The most mature CLI framework in the Go ecosystem, providing argument parsing, help generation, and more |
| **Regex Matching** | The core mechanism for secret format detection, paired with keyword pre-filtering for performance |
| **Keyword Pre-filtering** | Fast string-contains checks to eliminate irrelevant lines before regex matching, significantly improving scan speed |

### Roadmap

#### v1.1.0 -- Enhanced Detection

- Additional Chinese cloud provider rules (Volcengine, UCloud, QingCloud, etc.)
- Support for `.secretshieldignore` files (similar to `.gitignore` syntax)
- Dedicated `.env` file parsing and detection

#### v1.2.0 -- Team Collaboration

- WeCom (Enterprise WeChat) webhook alerts
- Scan result baseline feature, reporting only new leaks
- Configuration file `secretshield.yaml` for centralized scan parameter management

#### v2.0.0 -- Architecture Upgrade

- Plugin-based rule engine architecture with external rule pack support
- HTML visual report output
- Server mode with API support and scheduled scan tasks
- Incremental scanning, scanning only changed files

### Community Contribution Areas

We welcome contributions in the following areas:

- 🆕 New detection rules (especially Chinese cloud providers and development platforms)
- 🌐 Multi-language report templates (Japanese, Korean, etc.)
- 📦 New webhook channels (WeCom, Slack, etc.)
- 📖 Documentation improvements and translations
- 🐛 Bug fixes and performance optimizations

---

## 📦 Build & Deployment Guide

### Building from Source

```bash
git clone https://github.com/gitstq/SecretShield-CLI.git
cd SecretShield-CLI
go build -o bin/secretshield .
```

### Cross-Compilation

Use the `release` target in the Makefile for cross-compilation:

```bash
make release
```

Or specify the target platform manually:

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -o secretshield-linux-amd64 .

# Linux arm64
GOOS=linux GOARCH=arm64 go build -o secretshield-linux-arm64 .

# macOS amd64
GOOS=darwin GOARCH=amd64 go build -o secretshield-darwin-amd64 .

# macOS arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o secretshield-darwin-arm64 .

# Windows amd64
GOOS=windows GOARCH=amd64 go build -o secretshield-windows-amd64.exe .
```

### Release Downloads

Pre-built binaries are available on the [GitHub Releases](https://github.com/gitstq/SecretShield-CLI/releases) page.

### Platform Compatibility

| Platform | Architecture | Support |
|----------|-------------|---------|
| Linux | amd64 | ✅ |
| Linux | arm64 | ✅ |
| macOS | amd64 | ✅ |
| macOS | arm64 (Apple Silicon) | ✅ |
| Windows | amd64 | ✅ |

> **Minimum Requirement**: Go 1.21+ build environment (only needed at compile time; no Go dependency at runtime)

---

## 🤝 Contributing Guide

We warmly welcome community contributions! Below is the basic workflow for contributing.

### PR Submission Guidelines

1. **Fork** this repository and create your feature branch: `git checkout -b feature/your-feature-name`
2. **Write code**, ensuring all tests pass: `go test ./... -v`
3. **Commit messages** should follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:
   - `feat: add Volcengine detection rules`
   - `fix: resolve regex false positive issue`
   - `docs: update README usage instructions`
   - `refactor: optimize scanning engine performance`
4. **Push** to your fork: `git push origin feature/your-feature-name`
5. Create a **Pull Request** with a detailed description of your changes

### Issue Reporting Guidelines

When submitting an issue, please include the following information:

- **Problem Description**: Clearly describe the issue encountered or the feature you expect
- **Reproduction Steps**: For bugs, provide complete steps to reproduce
- **Environment Information**: Operating system, Go version, SecretShield-CLI version
- **Log Output**: Include relevant logs or screenshots if available

### Development Environment Setup

```bash
# Clone the repository
git clone https://github.com/gitstq/SecretShield-CLI.git
cd SecretShield-CLI

# Install dependencies
go mod download

# Run tests
go test ./... -v

# Code linting
go vet ./...

# Local build
go build -o bin/secretshield .

# Run
./bin/secretshield version
```

---

## 📄 License

This project is licensed under the **MIT License**.

```
MIT License

Copyright (c) 2024 SecretShield Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

This means you are free to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the software, with the only requirement being that the copyright notice and license declaration are preserved.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/gitstq/SecretShield-CLI">SecretShield Contributors</a>
</p>
