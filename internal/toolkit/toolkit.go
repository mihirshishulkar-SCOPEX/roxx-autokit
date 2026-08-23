package toolkit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mihirshishulkar-SCOPEX/roxx-autokit/internal/registry"
	"github.com/mihirshishulkar-SCOPEX/roxx-autokit/internal/updater"
)

// Toolkit manages the local installation and status of all tools
type Toolkit struct {
	BinDir string
	u      *updater.Updater
}

// New creates a Toolkit
func New() *Toolkit {
	return &Toolkit{
		BinDir: filepath.Join(os.Getenv("HOME"), ".local", "bin"),
		u:      updater.New(),
	}
}

// Install installs a single tool by name
func (tk *Toolkit) Install(name string) error {
	green := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed, color.Bold)
	cyan := color.New(color.FgCyan)

	for _, t := range registry.DefaultRegistry {
		if strings.EqualFold(t.Name, name) {
			cyan.Printf("\n⚡ Installing %s...\n", t.Name)
			if err := tk.u.UpdateTool(t); err != nil {
				red.Printf("[✗] Failed: %v\n", err)
				return err
			}
			green.Printf("[✓] %s installed successfully!\n", t.Name)
			return nil
		}
	}
	return fmt.Errorf("tool '%s' not found in registry", name)
}

// List prints all tools in the registry with their install status
func (tk *Toolkit) List() error {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	cyan.Println("\n╔═══════════════════════════════════════════════════════════════╗")
	cyan.Println("║              ROXX-AUTOKIT — Tool Registry                    ║")
	cyan.Println("╚═══════════════════════════════════════════════════════════════╝\n")

	// Group by category
	cats := make(map[string][]registry.ToolEntry)
	catOrder := []string{}
	for _, t := range registry.DefaultRegistry {
		if _, ok := cats[t.Category]; !ok {
			catOrder = append(catOrder, t.Category)
		}
		cats[t.Category] = append(cats[t.Category], t)
	}

	totalInstalled, totalTools := 0, 0
	for _, cat := range catOrder {
		tools := cats[cat]
		cyan.Printf("\n  📂 %s\n", strings.ToUpper(cat))
		for _, t := range tools {
			totalTools++
			installed := tk.isInstalled(t.Binary)
			status := ""
			if installed {
				totalInstalled++
				green.Printf("    ✅ %-20s", t.Name)
			} else {
				yellow.Printf("    ⚠️  %-20s", t.Name)
			}
			_ = status

			// Check version
			ver := tk.getVersion(t.Binary)
			if ver != "" {
				green.Printf(" %s", ver)
			}
			fmt.Printf("  [%s] github.com/%s\n", t.InstallType, t.Repo)
		}
	}

	fmt.Println()
	cyan.Printf("  📊 Installed: ")
	if totalInstalled == totalTools {
		green.Printf("%d/%d", totalInstalled, totalTools)
	} else {
		red.Printf("%d/%d", totalInstalled, totalTools)
	}
	cyan.Printf(" tools\n\n")
	return nil
}

// Status shows detailed status of all tools and daemon
func (tk *Toolkit) Status() error {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	cyan.Println("\n⚡ ROXX-AUTOKIT Status\n")

	// Daemon status
	pidFile := "/tmp/roxx-autokit.pid"
	if data, err := os.ReadFile(pidFile); err == nil {
		pid := strings.TrimSpace(string(data))
		// Check if process is running
		checkCmd := exec.Command("kill", "-0", pid)
		if checkCmd.Run() == nil {
			green.Printf("  🟢 Daemon: RUNNING (PID %s)\n", pid)
		} else {
			yellow.Printf("  🔴 Daemon: STOPPED\n")
			os.Remove(pidFile)
		}
	} else {
		yellow.Println("  🔴 Daemon: NOT RUNNING")
	}

	// Last update time from log
	logPath := os.Getenv("HOME") + "/.local/share/roxx-autokit/daemon.log"
	if info, err := os.Stat(logPath); err == nil {
		cyan.Printf("  📋 Last log: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
	}

	// Tool counts
	installed, total := 0, 0
	for _, t := range registry.DefaultRegistry {
		if t.Enabled {
			total++
			if tk.isInstalled(t.Binary) {
				installed++
			}
		}
	}

	cyan.Printf("  🔧 Tools: %d/%d installed\n", installed, total)
	cyan.Printf("  📁 Bin dir: %s\n\n", tk.BinDir)

	// Show recently discovered tools
	discoveredPath := os.Getenv("HOME") + "/.local/share/roxx-autokit/discovered_tools.json"
	if info, err := os.Stat(discoveredPath); err == nil {
		yellow.Printf("  🔭 Discovered tools log: %s (%s)\n",
			discoveredPath, formatSize(info.Size()))
	}

	green.Printf("  ⏰ Current time: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

func (tk *Toolkit) isInstalled(binary string) bool {
	// Check ~/.local/bin
	if _, err := os.Stat(filepath.Join(tk.BinDir, binary)); err == nil {
		return true
	}
	// Check PATH
	_, err := exec.LookPath(binary)
	return err == nil
}

func (tk *Toolkit) getVersion(binary string) string {
	// Try common version flags
	for _, flag := range []string{"--version", "-version", "version"} {
		cmd := exec.Command(binary, flag)
		out, err := cmd.CombinedOutput()
		if err == nil && len(out) > 0 {
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			if len(lines) > 0 {
				// Extract first meaningful version string
				v := strings.TrimSpace(lines[0])
				if len(v) > 60 {
					v = v[:60] + "..."
				}
				return "(" + v + ")"
			}
		}
	}
	return ""
}

func formatSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "KMGTPE"[exp])
}
