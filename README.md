<div align="center">

```
██████╗  ██████╗ ██╗  ██╗██╗  ██╗      █████╗ ██╗   ██╗████████╗ ██████╗ ██╗  ██╗██╗████████╗
██╔══██╗██╔═══██╗╚██╗██╔╝╚██╗██╔╝     ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██║ ██╔╝██║╚══██╔══╝
██████╔╝██║   ██║ ╚███╔╝  ╚███╔╝█████╗███████║██║   ██║   ██║   ██║   ██║█████╔╝ ██║   ██║   
██╔══██╗██║   ██║ ██╔██╗  ██╔██╗╚════╝██╔══██║██║   ██║   ██║   ██║   ██║██╔═██╗ ██║   ██║   
██║  ██║╚██████╔╝██╔╝ ██╗██╔╝ ██╗     ██║  ██║╚██████╔╝   ██║   ╚██████╔╝██║  ██╗██║   ██║   
╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝    ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═╝   ╚═╝   
```

**AUTONOMOUS. SELF-UPDATING. RELENTLESS.**

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-red?style=for-the-badge)
![Auto-Update](https://img.shields.io/badge/Auto--Update-30min-FF6600?style=for-the-badge)
![Tools](https://img.shields.io/badge/Tools-30+-green?style=for-the-badge)
![GitHub Sync](https://img.shields.io/badge/GitHub-Auto--Sync-181717?style=for-the-badge&logo=github)

</div>

---

## ⚡ What is ROXX-AUTOKIT?

**ROXX-AUTOKIT** is an autonomous, self-updating security arsenal manager. It:

- 🔄 **Auto-updates every 30 minutes** — all 30+ security tools always at latest version
- 🤖 **Self-updates itself** — pulls its own source from GitHub and rebuilds
- 🔭 **Discovers new tools** — scans GitHub trending for new high-quality security tools  
- 📤 **Syncs to GitHub** — pushes updated registries and tool lists to your repo automatically
- ⚡ **8 parallel workers** — blazing fast concurrent updates
- 🎯 **Zero manual work** — install once, hunt forever

---

## 🚀 One Command Install

```bash
git clone https://github.com/mihirshishulkar-SCOPEX/roxx-autokit.git
cd roxx-autokit
chmod +x install.sh && ./install.sh
```

---

## 📋 Commands

```bash
roxx-autokit daemon      # 🔥 Start autonomous 30-min update daemon
roxx-autokit update      # ⚡ Immediately update ALL tools
roxx-autokit list        # 📋 List all tools + install status  
roxx-autokit status      # 💓 Show daemon + system status
roxx-autokit discover    # 🔭 Hunt for new tools on GitHub
roxx-autokit sync        # 📤 Push updates to your GitHub repo
roxx-autokit install <t> # 🔧 Install a specific tool
```

---

## 🔧 Tracked Tools (30+)

| Category | Tools |
|----------|-------|
| **Recon** | subfinder, amass, assetfinder, findomain, chaos-client, katana, gau, waybackurls |
| **DNS** | dnsx, shuffledns, puredns |
| **HTTP** | httpx, httprobe |
| **Scanner** | nuclei, dalfox, naabu, nmap |
| **Fuzzer** | ffuf, feroxbuster, gobuster |
| **Secrets** | trufflehog, gitleaks |
| **Cloud** | cloudlist, s3scanner |
| **OAST** | interactsh-client |
| **Util** | anew, qsreplace, notify, jwt_tool |

---

## ⚙️ Configuration

```bash
# Set your GitHub token for full sync
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx

# Set your target repo
export GITHUB_REPO=youruser/yourrepo
```

Config file: `~/.config/roxx-autokit/config.toml`

---

## 🏗 Architecture

```
roxx-autokit/
├── main.go                          # Entry + banner
├── cmd/root.go                      # CLI commands (cobra)
├── internal/
│   ├── updater/updater.go           # Core update engine (8 workers)
│   ├── daemon/daemon.go             # 30-min cron daemon  
│   ├── github/github.go             # GitHub sync + discovery
│   ├── toolkit/toolkit.go           # Tool management + status
│   └── registry/registry.go        # Master tool registry (30+ tools)
├── .github/workflows/auto-build.yml # CI/CD auto-build pipeline
└── install.sh                       # One-command installer
```

---

## 🔄 Self-Update Flow

```
Every 30 minutes:
  1. git pull (self-update source)
  2. go build (rebuild binary)
  3. Update all 30+ security tools (8 parallel workers)
  4. Search GitHub for new security tools
  5. Generate TOOLS.md + registry JSON
  6. git commit + push to GitHub
```

---

<div align="center">

**Built by ROXX'S SLAVE — Autonomous Bug Bounty Intelligence**

⚡ *LOCKED. LOADED. UNCHAINED.*

</div>
