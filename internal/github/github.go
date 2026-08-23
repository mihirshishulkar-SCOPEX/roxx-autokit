package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mihirshishulkar-SCOPEX/roxx-autokit/internal/registry"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// GitHubClient handles all GitHub API and repo sync operations
type GitHubClient struct {
	Token     string
	RepoOwner string
	RepoName  string
	Client    *http.Client
	LocalPath string // path to local git repo clone
}

// New creates a GitHubClient from environment variables
func New() *GitHubClient {
	owner, repo := parseRepoEnv()
	localPath := filepath.Join(os.Getenv("HOME"), "roxxs-slave")

	return &GitHubClient{
		Token:     os.Getenv("GITHUB_TOKEN"),
		RepoOwner: owner,
		RepoName:  repo,
		LocalPath: localPath,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// parseRepoEnv reads GITHUB_REPO env var (format: "owner/repo")
func parseRepoEnv() (owner, repo string) {
	repoEnv := os.Getenv("GITHUB_REPO")
	if repoEnv == "" {
		repoEnv = "mihirshishulkar-SCOPEX/roxxs-slave" // default
	}
	parts := strings.SplitN(repoEnv, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "mihirshishulkar-SCOPEX", "roxxs-slave"
}

// SyncAll syncs the local repository to GitHub with all updates
func (g *GitHubClient) SyncAll() error {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	cyan.Println("\n🔄 [GITHUB SYNC] Starting repository sync...")

	// 1. Generate updated tool registry JSON
	if err := g.generateRegistryFile(); err != nil {
		yellow.Printf("  [!] Failed to generate registry: %v\n", err)
	}

	// 2. Generate TOOLS.md with current tool versions
	if err := g.generateToolsMD(); err != nil {
		yellow.Printf("  [!] Failed to generate TOOLS.md: %v\n", err)
	}

	// 3. Git commit and push
	if err := g.gitCommitAndPush(); err != nil {
		yellow.Printf("  [!] Git sync failed: %v\n", err)
		return err
	}

	green.Printf("  ✅ Repository synced: %s/%s\n", g.RepoOwner, g.RepoName)
	return nil
}

// generateRegistryFile writes the current tool registry as JSON
func (g *GitHubClient) generateRegistryFile() error {
	outDir := filepath.Join(g.LocalPath, "configs")
	os.MkdirAll(outDir, 0755)

	type RegistryOutput struct {
		GeneratedAt time.Time              `json:"generated_at"`
		Version     string                 `json:"version"`
		Tools       []registry.ToolEntry   `json:"tools"`
		Stats       map[string]int         `json:"stats"`
	}

	stats := make(map[string]int)
	for _, t := range registry.DefaultRegistry {
		stats["total"]++
		if t.Enabled {
			stats["enabled"]++
		}
		stats["cat_"+t.Category]++
	}

	out := RegistryOutput{
		GeneratedAt: time.Now().UTC(),
		Version:     "1.0.0",
		Tools:       registry.DefaultRegistry,
		Stats:       stats,
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outDir, "tools_registry.json"), data, 0644)
}

// generateToolsMD writes TOOLS.md with a human-readable table
func (g *GitHubClient) generateToolsMD() error {
	var sb strings.Builder

	sb.WriteString("# 🔧 ROXX-AUTOKIT Tool Registry\n\n")
	sb.WriteString(fmt.Sprintf("> Auto-generated on %s\n\n", time.Now().UTC().Format("2006-01-02 15:04:05 UTC")))
	sb.WriteString("| Tool | Category | Priority | Install | Tags |\n")
	sb.WriteString("|------|----------|----------|---------|------|\n")

	for _, t := range registry.DefaultRegistry {
		if t.Enabled {
			sb.WriteString(fmt.Sprintf("| [%s](https://github.com/%s) | %s | %d | `%s` | %s |\n",
				t.Name, t.Repo, t.Category, t.Priority, t.InstallType,
				strings.Join(t.Tags, ", "),
			))
		}
	}

	sb.WriteString(fmt.Sprintf("\n---\n*Total tools: %d*\n", len(registry.DefaultRegistry)))

	return os.WriteFile(filepath.Join(g.LocalPath, "TOOLS.md"), []byte(sb.String()), 0644)
}

// gitCommitAndPush commits all changes and pushes to origin
func (g *GitHubClient) gitCommitAndPush() error {
	if _, err := os.Stat(g.LocalPath); os.IsNotExist(err) {
		return fmt.Errorf("local repo not found at %s", g.LocalPath)
	}

	run := func(args ...string) error {
		cmd := exec.Command("git", args...)
		cmd.Dir = g.LocalPath
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git %s: %s", strings.Join(args, " "), strings.TrimSpace(string(out)))
		}
		return nil
	}

	// Configure git identity
	_ = run("config", "user.email", "roxx-autokit@auto.bot")
	_ = run("config", "user.name", "ROXX-AUTOKIT")

	// Stage all changes
	if err := run("add", "-A"); err != nil {
		return err
	}

	// Check if there's anything to commit
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.LocalPath
	out, _ := cmd.Output()
	if len(strings.TrimSpace(string(out))) == 0 {
		color.New(color.FgYellow).Println("  [i] No changes to commit — repo already up to date")
		return nil
	}

	commitMsg := fmt.Sprintf("🤖 auto-update: tools registry + TOOLS.md [%s]",
		time.Now().UTC().Format("2006-01-02 15:04:05"))

	if err := run("commit", "-m", commitMsg); err != nil {
		return err
	}

	if err := run("push", "origin", "main"); err != nil {
		// Try master
		if err2 := run("push", "origin", "master"); err2 != nil {
			return fmt.Errorf("push failed: %v | %v", err, err2)
		}
	}

	color.New(color.FgGreen).Printf("  [✓] Pushed: %s\n", commitMsg)
	return nil
}

// DiscoverNewTools searches GitHub for trending security tools and adds them
func (g *GitHubClient) DiscoverNewTools() error {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	cyan.Println("\n🔭 [DISCOVER] Hunting for new security tools on GitHub...")

	// Search queries for high-quality security tools
	queries := []string{
		"bug+bounty+tool+language:go+stars:>200",
		"recon+subdomain+language:go+stars:>100",
		"vulnerability+scanner+language:go+stars:>150",
		"penetration+testing+language:python+stars:>300",
	}

	discovered := 0
	for _, q := range queries {
		repos, err := g.searchGitHub(q)
		if err != nil {
			yellow.Printf("  [!] Search failed for '%s': %v\n", q, err)
			continue
		}

		for _, repo := range repos {
			if !g.isAlreadyTracked(repo.FullName) {
				g.addToDiscoveryLog(repo)
				discovered++
				green.Printf("  [+] Discovered: %s (%d ⭐)\n", repo.FullName, repo.Stars)
			}
		}
	}

	cyan.Printf("\n  📊 Discovered %d new tools\n", discovered)
	return nil
}

type GitHubRepo struct {
	FullName    string `json:"full_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Language    string `json:"language"`
	Topics      []string `json:"topics"`
	HTMLURL     string `json:"html_url"`
}

func (g *GitHubClient) searchGitHub(query string) ([]GitHubRepo, error) {
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&sort=stars&order=desc&per_page=10", query)

	req, _ := http.NewRequest("GET", url, nil)
	if g.Token != "" {
		req.Header.Set("Authorization", "token "+g.Token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "roxx-autokit/1.0")

	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub search returned %d", resp.StatusCode)
	}

	var result struct {
		Items []GitHubRepo `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (g *GitHubClient) isAlreadyTracked(fullName string) bool {
	for _, t := range registry.DefaultRegistry {
		if strings.EqualFold(t.Repo, fullName) {
			return true
		}
	}
	return false
}

func (g *GitHubClient) addToDiscoveryLog(repo GitHubRepo) {
	logPath := filepath.Join(os.Getenv("HOME"), ".local", "share", "roxx-autokit", "discovered_tools.json")
	os.MkdirAll(filepath.Dir(logPath), 0755)

	var existing []map[string]interface{}
	if data, err := os.ReadFile(logPath); err == nil {
		json.Unmarshal(data, &existing)
	}

	entry := map[string]interface{}{
		"repo":        repo.FullName,
		"name":        repo.Name,
		"description": repo.Description,
		"stars":       repo.Stars,
		"language":    repo.Language,
		"discovered":  time.Now().UTC().Format(time.RFC3339),
	}

	existing = append(existing, entry)
	data, _ := json.MarshalIndent(existing, "", "  ")
	os.WriteFile(logPath, data, 0644)
}

// CreateOrUpdateFile pushes a file directly to GitHub via API (no local clone needed)
func (g *GitHubClient) CreateOrUpdateFile(path, content, message string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s",
		g.RepoOwner, g.RepoName, path)

	// Get current SHA if file exists
	var sha string
	req, _ := http.NewRequest("GET", url, nil)
	if g.Token != "" {
		req.Header.Set("Authorization", "token "+g.Token)
	}
	req.Header.Set("User-Agent", "roxx-autokit/1.0")

	resp, err := g.Client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		var fileInfo struct {
			SHA string `json:"sha"`
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, &fileInfo)
		sha = fileInfo.SHA
	}

	// Encode content as base64
	encoded := encodeBase64([]byte(content))

	payload := map[string]interface{}{
		"message": message,
		"content": encoded,
	}
	if sha != "" {
		payload["sha"] = sha
	}

	data, _ := json.Marshal(payload)
	putReq, _ := http.NewRequest("PUT", url, bytes.NewReader(data))
	if g.Token != "" {
		putReq.Header.Set("Authorization", "token "+g.Token)
	}
	putReq.Header.Set("Content-Type", "application/json")
	putReq.Header.Set("User-Agent", "roxx-autokit/1.0")

	putResp, err := g.Client.Do(putReq)
	if err != nil {
		return err
	}
	defer putResp.Body.Close()

	if putResp.StatusCode != 200 && putResp.StatusCode != 201 {
		body, _ := io.ReadAll(putResp.Body)
		return fmt.Errorf("GitHub API returned %d: %s", putResp.StatusCode, string(body))
	}

	return nil
}

func encodeBase64(data []byte) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var buf bytes.Buffer
	for i := 0; i < len(data); i += 3 {
		b0 := data[i]
		var b1, b2 byte
		if i+1 < len(data) {
			b1 = data[i+1]
		}
		if i+2 < len(data) {
			b2 = data[i+2]
		}
		buf.WriteByte(chars[b0>>2])
		buf.WriteByte(chars[((b0&0x3)<<4)|(b1>>4)])
		if i+1 < len(data) {
			buf.WriteByte(chars[((b1&0xf)<<2)|(b2>>6)])
		} else {
			buf.WriteByte('=')
		}
		if i+2 < len(data) {
			buf.WriteByte(chars[b2&0x3f])
		} else {
			buf.WriteByte('=')
		}
	}
	return buf.String()
}
