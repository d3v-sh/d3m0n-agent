```
 ██████╗ ██████╗ ███╗   ███╗ ██████╗ ███╗   ██╗
 ██╔══██╗╚════██╗████╗ ████║██╔═████╗████╗  ██║
 ██║  ██║ █████╔╝██╔████╔██║██║██╔██║██╔██╗ ██║
 ██║  ██║ ╚═══██╗██║╚██╔╝██║████╔╝██║██║╚██╗██║
 ██████╔╝██████╔╝██║ ╚═╝ ██║╚██████╔╝██║ ╚████║
 ╚═════╝ ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝
          [ AI-Powered Cybersecurity Agent ]
               [ by d3v-sh | v1.0.0 ]
```

<div align="center">

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Kali-important?style=for-the-badge)
![AI](https://img.shields.io/badge/AI-OpenRouter%20%7C%20Local%20LLM-blueviolet?style=for-the-badge)
![Build](https://img.shields.io/github/actions/workflow/status/d3v-sh/d3m0n/build.yml?style=for-the-badge)

**A fully autonomous AI-powered cybersecurity agent built in Go.**
Runs real tools. Returns real results. No simulated output.

</div>

---

## What is d3m0n?

d3m0n is an AI agent for offensive security work. You describe a task in natural language — the agent decides which tools to call, executes them on your system, reasons over the real output, and continues until the task is complete.

It is not a wrapper that fakes tool output. When you ask it to scan a target, it actually runs `nmap`. When you ask it to enumerate subdomains, it actually runs `amass`. Every tool call is logged to SQLite with full arguments and output.

Built in Go — compiles to a single binary with no runtime dependencies.

---

## Features

- **Autonomous agent loop** — model calls tools, reasons over real output, chains multiple tools automatically
- **22 integrated security tools** — from recon to exploitation to reporting
- **CVE lookup** — searches NVD database for vulnerabilities in discovered services
- **Target memory** — agent saves and recalls information about targets across sessions
- **Professional report generation** — generates Markdown and dark-themed HTML pentest reports from conversation history
- **Session logging** — every tool call stored in SQLite with full args and output
- **Config-driven** — model, timeout, safe mode, allowed paths all set in `config.yaml`
- **Local LLM support** — works with Ollama, LM Studio, or any OpenAI-compatible endpoint
- **Safe mode** — path traversal protection, restricts file operations to configured directories
- **Streaming output** — responses stream token by token, re-rendered with glamour markdown
- **Beautiful CLI** — lipgloss styled UI, spinners, rendered markdown, slash commands

---

## Requirements

- Go 1.22+
- Linux (Kali recommended for full tool availability)
- API key from [OpenRouter](https://openrouter.ai) **or** a local LLM server (Ollama, LM Studio)

**Security tools** (install for full functionality):

```bash
# Kali / Debian
sudo apt install nmap whois curl gobuster ffuf amass nikto sqlmap sslscan gitleaks

# Python tools
pip install theHarvester sherlock-project spiderfoot

# Other
# EyeWitness: https://github.com/RedSiege/EyeWitness
# testssl.sh:  https://github.com/drwetter/testssl.sh
# recon-ng:    https://github.com/lanmaster53/recon-ng
# dradis:      https://dradisframework.com
```

---

## Installation

```bash
# clone
git clone https://github.com/d3v-sh/d3m0n.git
cd d3m0n

# configure
cp config.yaml.example config.yaml
# edit config.yaml and add your API key

# build
make build

# run
./d3m0n

# or run directly
make run
```

---

## Configuration

Copy `config.yaml.example` to `config.yaml`:

```yaml
# LLM
model: "anthropic/claude-3.5-haiku"     # any OpenRouter model
base_url: "https://openrouter.ai/api/v1"
api_key: "your-api-key-here"            # or set OPENROUTER_API_KEY env var

# Agent
max_history: 20      # messages kept in context window
timeout: 300         # seconds per request before timeout

# Security
safe_mode: true      # restrict file ops to allowed_paths only
allowed_paths:
  - .
  - ./output
  - /tmp

# Prompt
system_prompt: "You are a professional cybersecurity red team expert."
```

`config.yaml` is gitignored — your keys stay local.

### Local LLM

```yaml
# Ollama
model: "llama3.2"
base_url: "http://localhost:11434/v1"
api_key: "ollama"

# LM Studio
model: "local-model"
base_url: "http://localhost:1234/v1"
api_key: "lm-studio"
```

> Tool calling requires a capable model. Llama 3.2, Mistral 7B+, and Qwen2.5 work well. Models under 7B parameters are unreliable for agentic tasks. Avoid free-tier OpenRouter models (`:free` suffix) — they hallucinate tool outputs.

---

## Usage

```bash
./d3m0n
```

Describe tasks in plain English:

```
❯ perform a full recon on example.com
❯ scan 192.168.1.1 for open ports and look up CVEs for discovered services
❯ find all subdomains of target.com
❯ check ssl/tls configuration of example.com:443
❯ search for secrets in /path/to/repo
❯ run a whois on vulnweb.com then save results to output/vulnweb.txt
```

### Slash Commands

| Command | Description |
|---------|-------------|
| `/help` | Show available commands |
| `/tools` | List all integrated tools |
| `/report` | Generate pentest report from conversation |
| `/history` | Show last 10 sessions from database |
| `/targets` | Show all remembered targets |
| `/session` | Show current session info |
| `/config` | Show active configuration |
| `clear` | Clear screen and reset conversation history |
| `exit` | Quit d3m0n |

### Report Generation

After a recon session, run `/report` to generate a professional penetration testing report:

```
❯ perform recon on vulnweb.com
  ⚙  calling → run_whois
  ✓  run_whois completed
  ⚙  calling → run_amass
  ✓  run_amass completed
  ⚙  calling → run_curl
  ✓  run_curl completed

│ Assistant
  [summary of findings]

❯ /report
  ⏳ Generating report...
  ✓  Reports saved:
       MD:   output/report-20260601-143022.md
       HTML: output/report-20260601-143022.html
```

The HTML report is a dark-themed, fully formatted document with findings table, evidence, and recommendations.

---

## Integrated Tools

| Tool | Purpose |
|------|---------|
| `run_nmap` | Port scanning and service/version detection |
| `run_whois` | Domain and IP registration lookup |
| `run_curl` | HTTP probing, header inspection, content fetching |
| `run_gobuster` | Directory, DNS, and vhost brute forcing |
| `run_ffuf` | Fast web fuzzing with FUZZ keyword support |
| `run_amass` | Subdomain enumeration and DNS intelligence |
| `run_recon_ng` | Modular web reconnaissance framework |
| `run_theharvester` | Email, subdomain, and name OSINT from public sources |
| `run_sherlock` | Username search across 300+ social platforms |
| `run_spiderfoot` | Automated OSINT collection across multiple modules |
| `run_eyewitness` | Web service screenshot capture |
| `run_sslscan` | SSL/TLS configuration and cipher testing |
| `run_testssl` | Thorough SSL/TLS vulnerability analysis |
| `run_gitleaks` | Secret and credential detection in git repos |
| `run_sqlmap` | Automated SQL injection detection and exploitation |
| `run_crtsh` | Certificate transparency log lookup for subdomains |
| `run_dradis` | Import scan results into Dradis reporting framework |
| `search_cve` | NVD CVE database lookup for products and versions |
| `read_file` | Read files from disk (safe mode enforced) |
| `write_file` | Write or append to files (safe mode enforced) |
| `remember` | Save target intelligence to persistent memory |
| `recall` | Recall saved information about a target |

---

## Project Structure

```
d3m0n/
├── main.go                      # entry point, CLI loop, agent loop
├── config/
│   └── config.go                # config struct, YAML loading, defaults
├── tools/
│   ├── tools.go                 # tool definitions (OpenAI function schemas)
│   └── tool_wrapper.go          # tool implementations and dispatcher
├── logger/
│   └── logger.go                # SQLite session and tool call logging
├── ui/
│   └── ui.go                    # terminal UI (lipgloss styles, banner, prompts)
├── .github/
│   └── workflows/
│       └── build.yml            # GitHub Actions CI build
├── config.yaml.example          # example configuration file
├── Makefile                     # build, run, clean
├── .gitignore
├── go.mod
└── go.sum
```

---

## Session Logging

Every session is stored in `data/agent.db` (SQLite — no server required):

```bash
# list recent sessions
sqlite3 data/agent.db \
  "SELECT id, start_time FROM sessions ORDER BY start_time DESC LIMIT 10;"

# view tool calls from a session
sqlite3 data/agent.db \
  "SELECT tool, args FROM tool_logs WHERE session_id='20260601-143022';"

# view all remembered targets
sqlite3 data/agent.db "SELECT host, info, updated FROM targets;"
```

Or use the `/history` and `/targets` slash commands inside d3m0n.

---

## Security

- **Safe mode** restricts all file read/write to `allowed_paths` in config — path traversal is blocked
- **API key** loaded from `config.yaml` or `OPENROUTER_API_KEY` env var — never hardcoded
- **`config.yaml` is gitignored** — credentials stay local
- All tool execution happens locally on your machine
- Session data and reports stored locally in `data/` and `output/`

> d3m0n is intended for authorized security testing only. Only use against systems you own or have explicit written permission to test. The author is not responsible for misuse.

---

## Recommended Models

| Model | Speed | Tool Calling | Cost |
|-------|-------|-------------|------|
| `anthropic/claude-3.5-haiku` | Fast | Excellent | Low |
| `openai/gpt-4o-mini` | Fast | Reliable | Low |
| `google/gemini-2.0-flash-001` | Very fast | Good | Very low |
| `anthropic/claude-sonnet-4-5` | Medium | Best | Medium |
| `llama3.2` (Ollama) | Local | Good (7B+) | Free |
| `qwen2.5:14b` (Ollama) | Local | Best local | Free |

---

## Build

```bash
make build    # build binary → ./d3m0n
make run      # go run main.go
make clean    # remove binary and data/sessions directories
```

---

## Roadmap

- [ ] Docker + Kali image for fully isolated execution
- [ ] Telegram bot microservice (Python, separate repo)
- [ ] HTTP REST API for remote/multi-user access
- [ ] Kubernetes autoscaling support
- [ ] Playbook/workflow system for automated recon pipelines
- [ ] Web UI

---

## Contributing

Pull requests are welcome. For major changes, open an issue first.

```bash
git clone https://github.com/d3v-sh/d3m0n.git
cd d3m0n
git checkout -b feature/your-feature
# make changes
git push origin feature/your-feature
# open a pull request
```

---

## License

MIT — see [LICENSE](LICENSE) for details.

---

<div align="center">

Built in Go &nbsp;•&nbsp; Powered by AI &nbsp;•&nbsp; Made for red teamers

**[github.com/d3v-sh](https://github.com/d3v-sh)**

</div>