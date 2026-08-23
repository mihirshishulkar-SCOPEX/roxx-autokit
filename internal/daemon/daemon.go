package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/robfig/cron/v3"
	"github.com/mihirshishulkar-SCOPEX/roxx-autokit/internal/updater"
	"github.com/mihirshishulkar-SCOPEX/roxx-autokit/internal/github"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Start launches the autonomous update daemon
// Runs every 30 minutes, updates all tools + syncs to GitHub
func Start() error {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow)

	cyan.Println("\n🔥 [DAEMON] ROXX-AUTOKIT autonomous daemon starting...")
	yellow.Println("   Cycle: every 30 minutes")
	yellow.Println("   Tasks: update tools → discover new tools → sync GitHub")
	cyan.Println("   Press Ctrl+C to stop\n")

	// Write PID file
	pidFile := "/tmp/roxx-autokit.pid"
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		log.Warnf("Could not write PID file: %v", err)
	}
	defer os.Remove(pidFile)

	logFile, _ := os.OpenFile(
		os.Getenv("HOME")+"/.local/share/roxx-autokit/daemon.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)
	if logFile != nil {
		defer logFile.Close()
		log.SetOutput(logFile)
	}
	log.SetFormatter(&logrus.JSONFormatter{})

	// Run immediately on start
	cyan.Println("🚀 Running initial update cycle...")
	runCycle(green, yellow)

	// Schedule cron
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 */30 * * * *", func() {
		cyan.Printf("\n⏰ [%s] 30-minute cycle triggered\n", time.Now().Format("2006-01-02 15:04:05"))
		runCycle(green, yellow)
	})
	c.Start()
	defer c.Stop()

	// Status ticker every 10 minutes
	statusTicker := time.NewTicker(10 * time.Minute)
	defer statusTicker.Stop()

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sig:
			cyan.Println("\n🛑 [DAEMON] Shutdown signal received. Stopping gracefully...")
			return nil
		case t := <-statusTicker.C:
			cyan.Printf("💓 [DAEMON] Alive — %s | Next update in ~%s\n",
				t.Format("15:04:05"),
				nextRunIn(c),
			)
		}
	}
}

// runCycle runs one full update+sync cycle
func runCycle(green, yellow *color.Color) {
	log := logrus.New()
	log.Info("Starting update cycle")

	start := time.Now()

	// 1. Update all tools
	yellow.Println("  [1/3] Updating all installed tools...")
	u := updater.New()
	if err := u.UpdateAll(); err != nil {
		log.WithError(err).Error("UpdateAll failed")
	}

	// 2. Discover new tools from GitHub trending
	yellow.Println("  [2/3] Discovering new security tools...")
	gh := github.New()
	if err := gh.DiscoverNewTools(); err != nil {
		log.WithError(err).Warn("DiscoverNewTools had issues")
	}

	// 3. Sync to GitHub repo
	yellow.Println("  [3/3] Syncing to GitHub...")
	if err := gh.SyncAll(); err != nil {
		log.WithError(err).Warn("GitHub sync had issues")
	}

	green.Printf("  ✅ Cycle complete in %s\n\n", time.Since(start).Round(time.Second))
	log.WithField("duration", time.Since(start)).Info("Cycle complete")
}

// nextRunIn returns human-readable time until next cron execution
func nextRunIn(c *cron.Cron) string {
	entries := c.Entries()
	if len(entries) == 0 {
		return "unknown"
	}
	next := time.Until(entries[0].Next)
	return next.Round(time.Minute).String()
}
