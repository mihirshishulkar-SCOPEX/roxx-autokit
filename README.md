<div align="center">

<img width="100%" src="https://capsule-render.vercel.app/api?type=waving&color=0:0d0d0d,40:1a0000,100:8B0000&height=220&section=header&text=ROXX-AUTOKIT&fontSize=72&fontColor=FF6600&animation=fadeIn&fontAlignY=40&stroke=FF4400&strokeWidth=2&desc=Autonomous%20Self-Updating%20Security%20Arsenal&descAlignY=65&descSize=18&descColor=ff9966"/>

<img src="https://readme-typing-svg.demolab.com?font=Fira+Code&weight=900&size=15&duration=1800&pause=500&color=FF6600&background=0D0D0D&center=true&vCenter=true&width=900&lines=33+security+tools.+8+parallel+workers.+0+manual+effort.;Self-updates+every+30+minutes.+Autonomously.;Discovers+new+tools.+Pushes+to+GitHub.+Stays+sharp.;Clone.+Install.+Never+update+manually+again." alt="Typing"/>

<br/>

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white&labelColor=0d0d0d)
![Auto-Update](https://img.shields.io/badge/Auto--Update-Every_30min-FF6600?style=for-the-badge&labelColor=0d0d0d)
![Tools](https://img.shields.io/badge/Tools_Tracked-33+-brightgreen?style=for-the-badge&labelColor=0d0d0d)
![Workers](https://img.shields.io/badge/Parallel_Workers-8-blue?style=for-the-badge&labelColor=0d0d0d)
![License](https://img.shields.io/badge/License-MIT-red?style=for-the-badge&labelColor=0d0d0d)
![GitHub Sync](https://img.shields.io/badge/GitHub-Auto--Sync-181717?style=for-the-badge&logo=github&labelColor=0d0d0d)

<br/>

<img src="https://img.shields.io/badge/Kali_Linux-557C94?style=for-the-badge&logo=kalilinux&logoColor=white"/>
<img src="https://img.shields.io/badge/macOS-000000?style=for-the-badge&logo=apple&logoColor=white"/>
<img src="https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white"/>
<img src="https://img.shields.io/badge/Built_by-Mihir_Shishulkar-FF6600?style=for-the-badge&labelColor=0d0d0d"/>

</div>

---

## ⚡ What is ROXX-AUTOKIT?

**ROXX-AUTOKIT** is a blazing-fast, autonomous security tool manager built in Go.

It **installs, updates, discovers, and syncs** your entire security toolkit — without you touching it.

> Set it up once. It runs forever. Every 30 minutes, your arsenal is sharper than it was.

```
╔══════════════════════════════════════════════════════════════╗
║  Every 30 minutes, ROXX-AUTOKIT:                            ║
║                                                              ║
║  1. Self-updates its own binary from GitHub                 ║
║  2. Updates all 33 security tools in parallel               ║
║  3. Discovers new trending tools from GitHub                ║
║  4. Generates an updated tool registry (JSON + Markdown)    ║
║  5. Git commits and pushes everything to your GitHub repo   ║
╚══════════════════════════════════════════════════════════════╝
```

---

## 🚀 One Command Install

```bash
git clone https://github.com/mihirshishulkar-SCOPEX/roxx-autokit.git
cd roxx-autokit
chmod +x install.sh && ./install.sh
```

> That's it. The daemon starts automatically. Your tools stay updated forever.

---

## 🧠 How It Works

```
┌─────────────────────────────────────────────────────────────┐
│                    ROXX-AUTOKIT DAEMON                      │
│                    (every 30 minutes)                       │
└─────────────────────────────────────────────────────────────┘
           │
           ▼
    ┌─────────────┐
    │  SELF-UPDATE │  ← git pull + go build (rebuilds itself)
    └──────┬──────┘
           │
           ▼
    ┌─────────────────────────────────────────────┐
    │            UPDATE ALL TOOLS                 │
    │   8 parallel workers running concurrently   │
    │                                             │
    │  Go tools   → go install pkg@latest         │
    │  Releases   → GitHub API → latest binary    │
    │  Pip tools  → pip3 install --upgrade        │
    │  Scripts    → curl latest version           │
    └──────────────────┬──────────────────────────┘
                       │
                       ▼
              ┌─────────────────┐
              │  DISCOVER NEW   │  ← GitHub search API
              │     TOOLS       │    trending security repos
              └────────┬────────┘
                       │
                       ▼
              ┌─────────────────┐
              │  GITHUB SYNC    │  ← generate registry JSON
              │                 │    generate TOOLS.md
              │                 │    git commit + push
              └─────────────────┘
```

---

## 📋 Commands

```bash
# Start the autonomous 30-minute update daemon
roxx-autokit daemon

# Immediately update ALL 33 tools right now
roxx-autokit update

# List all tracked tools with install status + versions
roxx-autokit list

# Show daemon status, tool counts, last run time
roxx-autokit status

# Hunt GitHub trending for new security tools
roxx-autokit discover

# Push current registry and TOOLS.md to GitHub
roxx-autokit sync

# Install a specific tool from the registry
roxx-autokit install nuclei
roxx-autokit install ffuf
roxx-autokit install dalfox
```

---

## 🔧 Tools Arsenal (33 Tracked)

<details>
<summary><b>🔭 Recon & Subdomain (9 tools)</b></summary>

| Tool | Repo | Install Type |
|------|------|-------------|
| subfinder | [projectdiscovery/subfinder](https://github.com/projectdiscovery/subfinder) | Go |
| amass | [owasp-amass/amass](https://github.com/owasp-amass/amass) | Go |
| assetfinder | [tomnomnom/assetfinder](https://github.com/tomnomnom/assetfinder) | Go |
| findomain | [findomain/findomain](https://github.com/findomain/findomain) | Release |
| chaos-client | [projectdiscovery/chaos-client](https://github.com/projectdiscovery/chaos-client) | Go |
| katana | [projectdiscovery/katana](https://github.com/projectdiscovery/katana) | Go |
| gau | [lc/gau](https://github.com/lc/gau) | Go |
| waybackurls | [tomnomnom/waybackurls](https://github.com/tomnomnom/waybackurls) | Go |
| paramspider | [devanshbatham/paramspider](https://github.com/devanshbatham/paramspider) | Pip |

</details>

<details>
<summary><b>🌐 DNS Resolution (3 tools)</b></summary>

| Tool | Repo | Install Type |
|------|------|-------------|
| dnsx | [projectdiscovery/dnsx](https://github.com/projectdiscovery/dnsx) | Go |
| shuffledns | [projectdiscovery/shuffledns](https://github.com/projectdiscovery/shuffledns) | Go |
| puredns | [d3mondev/puredns](https://github.com/d3mondev/puredns) | Go |

</details>

<details>
<summary><b>🔌 HTTP Probing (2 tools)</b></summary>

| Tool | Repo | Install Type |
|------|------|-------------|
| httpx | [projectdiscovery/httpx](https://github.com/projectdiscovery/httpx) | Go |
| httprobe | [tomnomnom/httprobe](https://github.com/tomnomnom/httprobe) | Go |

</details>

<details>
<summary><b>🎯 Vulnerability Scanning (6 tools)</b></summary>

| Tool | Repo | Install Type |
|------|------|-------------|
| nuclei | [projectdiscovery/nuclei](https://github.com/projectdiscovery/nuclei) | Go |
| dalfox | [hahwul/dalfox](https://github.com/hahwul/dalfox) | Go |
| naabu | [projectdiscovery/naabu](https://github.com/projectdiscovery/naabu) | Go |
| nmap | [nmap/nmap](https://github.com/nmap/nmap) | System |
| sqlmap | [sqlmapproject/sqlmap](https://github.com/sqlmapproject/sqlmap) | Script |
| jwt_tool | [ticarpi/jwt_tool](https://github.com/ticarpi/jwt_tool) | Pip |

</details>

<details>
<summary><b>💣 Fuzzing (3 tools)</b></summary>

| Tool | Repo | Install Type |
|------|------|-------------|
| ffuf | [ffuf/ffuf](https://github.com/ffuf/ffuf) | Go |
| feroxbuster | [epi052/feroxbuster](https://github.com/epi052/feroxbuster) | Release |
| gobuster | [OJ/gobuster](https://github.com/OJ/gobuster) | Go |

</details>

<details>
<summary><b>🔑 Secrets & Exposure (3 tools)</b></summary>

| Tool | Repo | Install Type |
|------|------|-------------|
| trufflehog | [trufflesecurity/trufflehog](https://github.com/trufflesecurity/trufflehog) | Release |
| gitleaks | [gitleaks/gitleaks](https://github.com/gitleaks/gitleaks) | Release |
| secretfinder | [m4ll0k/SecretFinder](https://github.com/m4ll0k/SecretFinder) | Pip |

</details>

<details>
<summary><b>☁️ Cloud & Misconfig (2 tools)</b></summary>

| Tool | Repo | Install Type |
|------|------|-------------|
| cloudlist | [projectdiscovery/cloudlist](https://github.com/projectdiscovery/cloudlist) | Go |
| s3scanner | [sa7mon/S3Scanner](https://github.com/sa7mon/S3Scanner) | Pip |

</details>

<details>
<summary><b>📡 OAST / Callbacks (1 tool)</b></summary>

| Tool | Repo | Install Type |
|------|------|-------------|
| interactsh-client | [projectdiscovery/interactsh](https://github.com/projectdiscovery/interactsh) | Go |

</details>

<details>
<summary><b>🛠️ Utilities (4 tools)</b></summary>

| Tool | Repo | Install Type |
|------|------|-------------|
| anew | [tomnomnom/anew](https://github.com/tomnomnom/anew) | Go |
| qsreplace | [tomnomnom/qsreplace](https://github.com/tomnomnom/qsreplace) | Go |
| notify | [projectdiscovery/notify](https://github.com/projectdiscovery/notify) | Go |

</details>

---

## ⚙️ Configuration

Config file is auto-generated at `~/.config/roxx-autokit/config.toml`:

```toml
[daemon]
update_interval = "30m"    # How often to auto-update

[github]
repo   = "youruser/yourrepo"   # Your target repo for sync
auto_push = true
commit_prefix = "🤖 auto-update"

[updater]
workers = 8                # Parallel update workers
timeout = "120s"

[discovery]
enabled   = true           # Auto-discover new GitHub tools
min_stars = 100
```

### GitHub Sync Setup

```bash
# Set your GitHub Personal Access Token
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx

# Set your target repo (where registry + TOOLS.md get pushed)
export GITHUB_REPO=youruser/yourrepo

# Start daemon — it will now sync to GitHub every 30 min
roxx-autokit daemon
```

---

## 🏗 Project Structure

```
roxx-autokit/
│
├── main.go                              # Entry point + ASCII banner
├── go.mod / go.sum                      # Go module
│
├── cmd/
│   └── root.go                          # Cobra CLI — 7 commands
│
├── internal/
│   ├── registry/registry.go            # Master tool registry (33 tools)
│   ├── updater/updater.go              # Core engine — 8 parallel workers
│   ├── daemon/daemon.go                # 30-min cron daemon + graceful shutdown
│   ├── github/github.go               # GitHub API: sync + discovery
│   └── toolkit/toolkit.go             # list / status / install commands
│
├── .github/
│   └── workflows/
│       └── auto-build.yml             # CI/CD: cross-platform builds + auto-release
│
└── install.sh                          # Full one-command installer + systemd
```

---

## 🖥 Systemd Service

```bash
# Installed automatically by install.sh
# Service file: ~/.config/systemd/user/roxx-autokit.service

systemctl --user status roxx-autokit    # Check status
systemctl --user start roxx-autokit     # Start daemon
systemctl --user stop roxx-autokit      # Stop daemon
systemctl --user enable roxx-autokit    # Auto-start on login

# Follow live logs
journalctl --user -u roxx-autokit -f
```

---

## 🔄 CI/CD Pipeline

The included GitHub Actions workflow (`.github/workflows/auto-build.yml`) auto-builds:

| Platform | Arch |
|----------|------|
| Linux | amd64 |
| Linux | arm64 |
| macOS | amd64 |
| macOS | arm64 (Apple Silicon) |
| Windows | amd64 |

- Builds fire on every push to `main`
- Releases auto-created on tags (`v*`)
- Scheduled builds every 6 hours

---

## 📊 Performance

| Metric | Value |
|--------|-------|
| Binary size | ~8.4MB (stripped, static) |
| Parallel workers | 8 |
| Avg full update time | ~2-4 minutes |
| Memory usage (daemon) | ~15MB |
| Dependencies | Zero (single static binary) |

---

## 🤝 Adding Custom Tools

Edit [`internal/registry/registry.go`](internal/registry/registry.go) and add a new entry:

```go
{
    Name:        "mytool",
    Repo:        "author/mytool",
    Binary:      "mytool",
    InstallType: "go",               // "go", "release", "pip", "script"
    GoPackage:   "github.com/author/mytool@latest",
    Tags:        []string{"recon"},
    Category:    "recon",
    Priority:    2,
    Enabled:     true,
},
```

Then rebuild:

```bash
cd roxx-autokit && go build -o ~/.local/bin/roxx-autokit .
```

---

<div align="center">

<img width="100%" src="https://capsule-render.vercel.app/api?type=waving&color=0:8B0000,60:1a0000,100:0d0d0d&height=120&section=footer"/>

<img src="https://img.shields.io/badge/Built_by-Mihir_Shishulkar-FF6600?style=for-the-badge&labelColor=0d0d0d"/>
<img src="https://img.shields.io/badge/Mode-AUTONOMOUS-8B0000?style=for-the-badge&labelColor=0d0d0d"/>

**⚡ LOCKED. LOADED. UNCHAINED.**

*Your arsenal. Always sharp. Never stale.*

</div>
