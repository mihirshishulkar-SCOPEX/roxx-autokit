#!/usr/bin/env bash
# ================================================================
#  ROXX-AUTOKIT — Master Installer + Self-Update Bootstrap
#  Installs the tool, sets up systemd daemon, configures GitHub sync
#  Usage: bash install.sh
# ================================================================

set -euo pipefail

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; BOLD='\033[1m'; NC='\033[0m'
MAGENTA='\033[0;35m'

TOOL_NAME="roxx-autokit"
INSTALL_DIR="$HOME/.local/bin"
DATA_DIR="$HOME/.local/share/roxx-autokit"
CONFIG_DIR="$HOME/.config/roxx-autokit"
REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SELF_UPDATE_INTERVAL=30  # minutes

banner() {
cat << "BANNER"
 ██████╗  ██████╗ ██╗  ██╗██╗  ██╗      █████╗ ██╗   ██╗████████╗ ██████╗ ██╗  ██╗██╗████████╗
 ██╔══██╗██╔═══██╗╚██╗██╔╝╚██╗██╔╝     ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██║ ██╔╝██║╚══██╔══╝
 ██████╔╝██║   ██║ ╚███╔╝  ╚███╔╝█████╗███████║██║   ██║   ██║   ██║   ██║█████╔╝ ██║   ██║   
 ██╔══██╗██║   ██║ ██╔██╗  ██╔██╗╚════╝██╔══██║██║   ██║   ██║   ██║   ██║██╔═██╗ ██║   ██║   
 ██║  ██║╚██████╔╝██╔╝ ██╗██╔╝ ██╗     ██║  ██║╚██████╔╝   ██║   ╚██████╔╝██║  ██╗██║   ██║   
 ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝    ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═╝   ╚═╝   

              AUTONOMOUS. SELF-UPDATING. RELENTLESS.
BANNER
}

log()     { echo -e "${GREEN}[+]${NC} $1"; }
warn()    { echo -e "${YELLOW}[!]${NC} $1"; }
error()   { echo -e "${RED}[✗]${NC} $1"; }
section() { echo -e "\n${CYAN}${BOLD}══════════════════════════════════════${NC}"; \
            echo -e "${CYAN}${BOLD}  $1${NC}"; \
            echo -e "${CYAN}${BOLD}══════════════════════════════════════${NC}"; }

banner
echo ""

# ── PRE-CHECKS ─────────────────────────────────────────────────────────────────
section "STEP 1/6 — SYSTEM CHECKS"

check_dep() {
    if command -v "$1" &>/dev/null; then
        log "$1 found: $(command -v $1)"
    else
        warn "$1 not found — $2"
    fi
}

check_dep "go"     "Install from https://go.dev"
check_dep "git"    "Required for GitHub sync"
check_dep "curl"   "Required for downloads"
check_dep "jq"     "Optional, for JSON processing"
check_dep "pip3"   "Required for Python tools"
check_dep "gh"     "Optional, for GitHub auth"

# Check Go version
if command -v go &>/dev/null; then
    GO_VER=$(go version | awk '{print $3}' | sed 's/go//')
    log "Go version: $GO_VER"
fi

# ── DIRECTORIES ──────────────────────────────────────────────────────────────
section "STEP 2/6 — CREATING DIRECTORIES"

mkdir -p "$INSTALL_DIR" "$DATA_DIR" "$CONFIG_DIR"
log "Created: $INSTALL_DIR"
log "Created: $DATA_DIR"
log "Created: $CONFIG_DIR"

# Ensure ~/.local/bin is in PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    SHELL_RC="$HOME/.bashrc"
    [[ -f "$HOME/.zshrc" ]] && SHELL_RC="$HOME/.zshrc"
    echo "export PATH=\"\$PATH:$HOME/.local/bin\"" >> "$SHELL_RC"
    warn "Added $INSTALL_DIR to PATH in $SHELL_RC — restart shell or run: export PATH=\$PATH:$INSTALL_DIR"
fi

# ── BUILD ROXX-AUTOKIT ────────────────────────────────────────────────────────
section "STEP 3/6 — BUILDING ROXX-AUTOKIT"

if [ -f "$REPO_DIR/main.go" ]; then
    log "Building from source in $REPO_DIR..."
    cd "$REPO_DIR"
    
    # Download dependencies
    go mod tidy 2>/dev/null || {
        warn "go mod tidy failed — trying go mod download..."
        go mod download 2>/dev/null || true
    }
    
    # Build with optimizations
    CGO_ENABLED=0 go build \
        -ldflags="-s -w -X main.VERSION=$(date +%Y%m%d.%H%M)" \
        -o "$INSTALL_DIR/$TOOL_NAME" \
        . && log "Built: $INSTALL_DIR/$TOOL_NAME" || {
        error "Build failed — check Go installation"
        exit 1
    }
else
    warn "Source not found — downloading pre-built binary..."
    # Fallback: download from GitHub releases
    LATEST=$(curl -s https://api.github.com/repos/mihirshishulkar-SCOPEX/roxx-autokit/releases/latest \
        | grep '"tag_name"' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
    
    if [ -n "$LATEST" ]; then
        ARCH=$(uname -m)
        OS=$(uname -s | tr '[:upper:]' '[:lower:]')
        curl -L "https://github.com/mihirshishulkar-SCOPEX/roxx-autokit/releases/download/$LATEST/roxx-autokit-${OS}-${ARCH}" \
            -o "$INSTALL_DIR/$TOOL_NAME"
        chmod +x "$INSTALL_DIR/$TOOL_NAME"
        log "Downloaded $TOOL_NAME $LATEST"
    else
        error "Could not build or download $TOOL_NAME"
        exit 1
    fi
fi

# Verify installation
if "$INSTALL_DIR/$TOOL_NAME" --version &>/dev/null 2>&1 || [ -x "$INSTALL_DIR/$TOOL_NAME" ]; then
    log "Binary verified: $INSTALL_DIR/$TOOL_NAME"
else
    error "Binary verification failed"
    exit 1
fi

# ── WRITE CONFIG ─────────────────────────────────────────────────────────────
section "STEP 4/6 — WRITING CONFIGURATION"

cat > "$CONFIG_DIR/config.toml" << EOF
# ROXX-AUTOKIT Configuration
# Generated: $(date)

[daemon]
update_interval = "30m"    # How often to update tools
github_sync = true         # Push updates to GitHub

[github]
repo = "${GITHUB_REPO:-mihirshishulkar-SCOPEX/roxxs-slave}"
token_env = "GITHUB_TOKEN" # Environment variable name for GitHub token
auto_push = true
commit_prefix = "🤖 auto-update"

[updater]
workers = 8                # Parallel update workers
timeout = "120s"
bin_dir = "$HOME/.local/bin"

[discovery]
enabled = true             # Auto-discover new tools from GitHub
min_stars = 100            # Minimum stars to consider a tool
categories = ["recon", "scanner", "fuzzer", "secrets"]

[logging]
level = "info"
file = "$DATA_DIR/daemon.log"
EOF

log "Config written: $CONFIG_DIR/config.toml"

# ── SELF-UPDATE CRON ─────────────────────────────────────────────────────────
section "STEP 5/6 — SETTING UP SELF-UPDATE"

# Create self-update script
cat > "$DATA_DIR/self_update.sh" << 'SELFUPDATE'
#!/usr/bin/env bash
# ROXX-AUTOKIT Self-Update Script
# Pulls latest version from GitHub and rebuilds

set -euo pipefail

REPO_DIR="$(dirname "$(readlink -f "$0")")"
TOOL_DIR="$HOME/roxx-autokit"
LOG="$HOME/.local/share/roxx-autokit/self_update.log"

log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG"; }

log "=== SELF-UPDATE TRIGGERED ==="

# Find the source repo
if [ -d "$TOOL_DIR/.git" ]; then
    log "Pulling latest source..."
    cd "$TOOL_DIR"
    git fetch origin --quiet
    LOCAL=$(git rev-parse HEAD)
    REMOTE=$(git rev-parse origin/main 2>/dev/null || git rev-parse origin/master 2>/dev/null)
    
    if [ "$LOCAL" != "$REMOTE" ]; then
        log "New version detected — updating..."
        git pull origin main 2>/dev/null || git pull origin master 2>/dev/null
        
        # Rebuild
        log "Rebuilding binary..."
        CGO_ENABLED=0 go build \
            -ldflags="-s -w -X main.VERSION=$(date +%Y%m%d.%H%M)" \
            -o "$HOME/.local/bin/roxx-autokit" . && \
            log "Binary updated successfully!" || \
            log "ERROR: Rebuild failed — keeping old binary"
    else
        log "Already at latest version ($LOCAL)"
    fi
else
    log "WARNING: Source repo not found at $TOOL_DIR"
fi

log "Self-update complete."
SELFUPDATE

chmod +x "$DATA_DIR/self_update.sh"
log "Self-update script: $DATA_DIR/self_update.sh"

# Install cron job for self-update every 30 minutes
CRON_JOB="*/$SELF_UPDATE_INTERVAL * * * * $DATA_DIR/self_update.sh >> $DATA_DIR/self_update.log 2>&1"
(crontab -l 2>/dev/null | grep -v "roxx-autokit\|self_update.sh"; echo "$CRON_JOB") | crontab - 2>/dev/null && \
    log "Cron job installed (every $SELF_UPDATE_INTERVAL minutes)" || \
    warn "Could not install cron — run manually: $DATA_DIR/self_update.sh"

# ── SYSTEMD SERVICE ──────────────────────────────────────────────────────────
section "STEP 6/6 — SYSTEMD SERVICE"

SYSTEMD_DIR="$HOME/.config/systemd/user"
mkdir -p "$SYSTEMD_DIR"

cat > "$SYSTEMD_DIR/roxx-autokit.service" << EOF
[Unit]
Description=ROXX-AUTOKIT Autonomous Security Tool Updater
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=$HOME/.local/bin/roxx-autokit daemon
Restart=always
RestartSec=60
StandardOutput=append:$DATA_DIR/daemon.log
StandardError=append:$DATA_DIR/daemon.log
Environment="PATH=$HOME/.local/bin:/usr/local/go/bin:/usr/bin:/bin"
Environment="HOME=$HOME"
Environment="GITHUB_TOKEN=%GITHUB_TOKEN%"

[Install]
WantedBy=default.target
EOF

# Replace %GITHUB_TOKEN% with actual value if set
if [ -n "${GITHUB_TOKEN:-}" ]; then
    sed -i "s|%GITHUB_TOKEN%|$GITHUB_TOKEN|g" "$SYSTEMD_DIR/roxx-autokit.service"
    log "GitHub token embedded in service"
else
    sed -i "s|Environment=\"GITHUB_TOKEN=%GITHUB_TOKEN%\"||g" "$SYSTEMD_DIR/roxx-autokit.service"
    warn "GITHUB_TOKEN not set — GitHub sync will be limited"
fi

# Enable and start service
if systemctl --user daemon-reload 2>/dev/null; then
    systemctl --user enable roxx-autokit.service 2>/dev/null && \
        log "Service enabled: roxx-autokit.service" || true
    systemctl --user restart roxx-autokit.service 2>/dev/null && \
        log "Service started!" || warn "Could not start service — run manually: $TOOL_NAME daemon"
else
    warn "systemd not available — start manually: $TOOL_NAME daemon"
fi

# ── FINAL SUMMARY ────────────────────────────────────────────────────────────
echo ""
echo -e "${CYAN}${BOLD}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}${BOLD}║     ✅  ROXX-AUTOKIT INSTALLED SUCCESSFULLY!         ║${NC}"
echo -e "${CYAN}${BOLD}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}Commands:${NC}"
echo -e "  ${YELLOW}roxx-autokit daemon${NC}     — Start autonomous 30-min update daemon"
echo -e "  ${YELLOW}roxx-autokit update${NC}     — Immediately update all tools"
echo -e "  ${YELLOW}roxx-autokit list${NC}       — List all tools + install status"
echo -e "  ${YELLOW}roxx-autokit status${NC}     — Show daemon + tool status"
echo -e "  ${YELLOW}roxx-autokit discover${NC}   — Hunt for new security tools on GitHub"
echo -e "  ${YELLOW}roxx-autokit sync${NC}       — Push updates to your GitHub repo"
echo -e "  ${YELLOW}roxx-autokit install <tool>${NC} — Install a specific tool"
echo ""
echo -e "${MAGENTA}Set GITHUB_TOKEN env var to enable full GitHub sync:${NC}"
echo -e "  ${YELLOW}export GITHUB_TOKEN=your_token_here${NC}"
echo ""
echo -e "${RED}${BOLD}⚡ LOCKED. LOADED. UNCHAINED.${NC}"
echo ""
