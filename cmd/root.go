package cmd

import (
	"github.com/spf13/cobra"
	"github.com/mihirshishulkar-SCOPEX/roxx-autokit/internal/daemon"
	"github.com/mihirshishulkar-SCOPEX/roxx-autokit/internal/updater"
	"github.com/mihirshishulkar-SCOPEX/roxx-autokit/internal/toolkit"
	"github.com/mihirshishulkar-SCOPEX/roxx-autokit/internal/github"
	"github.com/fatih/color"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "roxx-autokit",
	Short: "ROXX-AUTOKIT — Autonomous self-updating security arsenal",
	Long: `ROXX-AUTOKIT is a relentless self-updating security toolkit.
It auto-updates all tools every 30 minutes, syncs new tools from 
the ecosystem, and pushes updates to your GitHub repo autonomously.`,
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the auto-update daemon (30-min cycle)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return daemon.Start()
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Immediately update all installed tools to latest versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		u := updater.New()
		return u.UpdateAll()
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync tool registry and push changes to GitHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		gh := github.New()
		return gh.SyncAll()
	},
}

var installCmd = &cobra.Command{
	Use:   "install [tool]",
	Short: "Install a specific tool from the registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tk := toolkit.New()
		return tk.Install(args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracked tools and their versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		tk := toolkit.New()
		return tk.List()
	},
}

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover and add new security tools from GitHub trending",
	RunE: func(cmd *cobra.Command, args []string) error {
		gh := github.New()
		return gh.DiscoverNewTools()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current status of all tools and last update time",
	RunE: func(cmd *cobra.Command, args []string) error {
		tk := toolkit.New()
		return tk.Status()
	},
}

func Execute() error {
	rootCmd.AddCommand(
		daemonCmd,
		updateCmd,
		syncCmd,
		installCmd,
		listCmd,
		discoverCmd,
		statusCmd,
	)

	// Check if no args — show help with style
	if len(os.Args) == 1 {
		color.New(color.FgYellow).Println("  Use 'roxx-autokit daemon' to start the autonomous update engine.")
		color.New(color.FgCyan).Println("  Use 'roxx-autokit --help' to see all commands.\n")
	}

	return rootCmd.Execute()
}
