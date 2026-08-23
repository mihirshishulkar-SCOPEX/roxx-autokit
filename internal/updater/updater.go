package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/mihirshishulkar-SCOPEX/roxx-autokit/internal/registry"
	"github.com/sirupsen/logrus"
)

var (
	green  = color.New(color.FgGreen, color.Bold)
	red    = color.New(color.FgRed, color.Bold)
	yellow = color.New(color.FgYellow)
	cyan   = color.New(color.FgCyan, color.Bold)
	log    = logrus.New()
)

// Updater orchestrates updates for all registered tools
type Updater struct {
	BinDir     string
	MaxWorkers int
	Timeout    time.Duration
	Client     *http.Client
}

// New creates a new Updater with sane defaults
func New() *Updater {
	binDir := filepath.Join(os.Getenv("HOME"), ".local", "bin")
	os.MkdirAll(binDir, 0755)

	return &Updater{
		BinDir:     binDir,
		MaxWorkers: 8, // Parallel update workers
		Timeout:    120 * time.Second,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UpdateAll updates every enabled tool in the registry concurrently
func (u *Updater) UpdateAll() error {
	cyan.Println("\n⚡ [ROXX-AUTOKIT] Starting full arsenal update...")

	tools := registry.DefaultRegistry
	enabled := make([]registry.ToolEntry, 0, len(tools))
	for _, t := range tools {
		if t.Enabled {
			enabled = append(enabled, t)
		}
	}

	// Semaphore for parallel workers
	sem := make(chan struct{}, u.MaxWorkers)
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(map[string]string)

	start := time.Now()

	for _, tool := range enabled {
		wg.Add(1)
		sem <- struct{}{} // acquire
		go func(t registry.ToolEntry) {
			defer wg.Done()
			defer func() { <-sem }() // release

			err := u.UpdateTool(t)
			mu.Lock()
			if err != nil {
				results[t.Name] = fmt.Sprintf("❌ FAILED: %v", err)
				red.Printf("  [✗] %-20s FAILED: %v\n", t.Name, err)
			} else {
				results[t.Name] = "✅ OK"
				green.Printf("  [✓] %-20s Updated\n", t.Name)
			}
			mu.Unlock()
		}(tool)
	}

	wg.Wait()
	elapsed := time.Since(start)

	// Summary
	ok, fail := 0, 0
	for _, v := range results {
		if strings.HasPrefix(v, "✅") {
			ok++
		} else {
			fail++
		}
	}

	cyan.Printf("\n🏁 Update complete in %s — %d OK, %d Failed\n\n", elapsed.Round(time.Second), ok, fail)
	return nil
}

// UpdateTool updates a single tool based on its install type
func (u *Updater) UpdateTool(t registry.ToolEntry) error {
	switch t.InstallType {
	case "go":
		return u.updateGo(t)
	case "release":
		return u.updateRelease(t)
	case "pip":
		return u.updatePip(t)
	case "script":
		return u.updateScript(t)
	default:
		return fmt.Errorf("unknown install type: %s", t.InstallType)
	}
}

// updateGo installs/updates a Go-based tool
func (u *Updater) updateGo(t registry.ToolEntry) error {
	if t.GoPackage == "" {
		return fmt.Errorf("no go package defined for %s", t.Name)
	}

	gobin := filepath.Join(os.Getenv("HOME"), ".local", "bin")
	cmd := exec.Command("go", "install", t.GoPackage)
	cmd.Env = append(os.Environ(),
		"GOBIN="+gobin,
		"CGO_ENABLED=0",
		"GOFLAGS=-mod=mod",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go install failed: %s — %v", strings.TrimSpace(string(out)), err)
	}
	return nil
}

// updatePip installs/updates a Python pip tool
func (u *Updater) updatePip(t registry.ToolEntry) error {
	pkg := t.PipPackage
	if pkg == "" {
		pkg = t.Name
	}

	cmd := exec.Command("pip3", "install", "--upgrade", "--quiet", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pip install failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// updateScript downloads and installs a script-based tool
func (u *Updater) updateScript(t registry.ToolEntry) error {
	if t.ScriptURL == "" {
		return nil // some script tools like nmap are system-installed
	}

	resp, err := u.Client.Get(t.ScriptURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dest := filepath.Join(u.BinDir, t.Binary+".py")
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return os.Chmod(dest, 0755)
}

// updateRelease downloads the latest GitHub release binary
func (u *Updater) updateRelease(t registry.ToolEntry) error {
	latest, downloadURL, err := u.getLatestRelease(t)
	if err != nil {
		return err
	}

	if downloadURL == "" {
		return fmt.Errorf("no compatible release found for %s (latest: %s)", t.Name, latest)
	}

	yellow.Printf("    → Downloading %s %s...\n", t.Name, latest)
	return u.downloadAndInstall(t, downloadURL)
}

// getLatestRelease queries GitHub API for the latest release
func (u *Updater) getLatestRelease(t registry.ToolEntry) (version, downloadURL string, err error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", t.Repo)

	req, _ := http.NewRequest("GET", apiURL, nil)

	// Use GitHub token if available
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "roxx-autokit/1.0")

	resp, err := u.Client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return "", "", fmt.Errorf("rate limited by GitHub API")
	}
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("GitHub API returned %d for %s", resp.StatusCode, t.Repo)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", err
	}

	version = release.TagName

	// Find a matching asset for current OS/arch
	os_name := runtime.GOOS
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	}

	for _, asset := range release.Assets {
		nameLower := strings.ToLower(asset.Name)
		if strings.Contains(nameLower, os_name) &&
			(strings.Contains(nameLower, arch) || strings.Contains(nameLower, "amd64")) &&
			!strings.Contains(nameLower, ".sha") {
			return version, asset.BrowserDownloadURL, nil
		}
	}

	return version, "", nil
}

// downloadAndInstall downloads a release asset and installs it
func (u *Updater) downloadAndInstall(t registry.ToolEntry, url string) error {
	resp, err := u.Client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpFile := filepath.Join(os.TempDir(), t.Name+"_download")
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.Copy(f, resp.Body); err != nil {
		return err
	}
	f.Close()

	dest := filepath.Join(u.BinDir, t.Binary)

	// Handle archives
	urlLower := strings.ToLower(url)
	if strings.Contains(urlLower, ".tar.gz") || strings.Contains(urlLower, ".tgz") {
		return u.extractTar(tmpFile, t.Binary, dest)
	} else if strings.Contains(urlLower, ".zip") {
		return u.extractZip(tmpFile, t.Binary, dest)
	}

	// Plain binary
	if err := os.Rename(tmpFile, dest); err != nil {
		// Cross-device: copy
		return u.copyFile(tmpFile, dest)
	}

	return os.Chmod(dest, 0755)
}

func (u *Updater) extractTar(src, binary, dest string) error {
	tmpDir := src + "_dir"
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("tar", "xzf", src, "-C", tmpDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tar extract failed: %v", err)
	}

	// Find the binary recursively
	var found string
	filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && info.Name() == binary {
			found = path
		}
		return nil
	})

	if found == "" {
		return fmt.Errorf("binary '%s' not found in archive", binary)
	}

	return u.copyFile(found, dest)
}

func (u *Updater) extractZip(src, binary, dest string) error {
	tmpDir := src + "_dir"
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("unzip", "-o", src, "-d", tmpDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unzip failed: %v", err)
	}

	var found string
	filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && info.Name() == binary {
			found = path
		}
		return nil
	})

	if found == "" {
		return fmt.Errorf("binary '%s' not found in zip", binary)
	}

	return u.copyFile(found, dest)
}

func (u *Updater) copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return os.Chmod(dst, 0755)
}
