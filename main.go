package main

import (
	"fmt"
	"os"

	"github.com/mihirshishulkar-SCOPEX/roxx-autokit/cmd"
	"github.com/fatih/color"
)

const BANNER = `
██████╗  ██████╗ ██╗  ██╗██╗  ██╗      █████╗ ██╗   ██╗████████╗ ██████╗ ██╗  ██╗██╗████████╗
██╔══██╗██╔═══██╗╚██╗██╔╝╚██╗██╔╝     ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██║ ██╔╝██║╚══██╔══╝
██████╔╝██║   ██║ ╚███╔╝  ╚███╔╝█████╗███████║██║   ██║   ██║   ██║   ██║█████╔╝ ██║   ██║   
██╔══██╗██║   ██║ ██╔██╗  ██╔██╗╚════╝██╔══██║██║   ██║   ██║   ██║   ██║██╔═██╗ ██║   ██║   
██║  ██║╚██████╔╝██╔╝ ██╗██╔╝ ██╗     ██║  ██║╚██████╔╝   ██║   ╚██████╔╝██║  ██╗██║   ██║   
╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝    ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═╝   ╚═╝   
`

const VERSION = "1.0.0"

func main() {
	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	cyan := color.New(color.FgCyan)

	red.Print(BANNER)
	yellow.Printf("  ⚡ ROXX-AUTOKIT v%s — Self-Updating Security Arsenal\n", VERSION)
	cyan.Println("  🔥 Autonomous. Relentless. Self-Evolving.\n")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
