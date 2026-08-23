package registry

// ToolEntry defines a security tool tracked by ROXX-AUTOKIT
type ToolEntry struct {
	Name        string   `toml:"name"`
	Repo        string   `toml:"repo"`         // GitHub owner/repo
	Binary      string   `toml:"binary"`       // binary name after install
	InstallType string   `toml:"install_type"` // "go", "release", "script", "pip"
	GoPackage   string   `toml:"go_package"`   // for install_type=go
	ScriptURL   string   `toml:"script_url"`   // for install_type=script
	PipPackage  string   `toml:"pip_package"`  // for install_type=pip
	Tags        []string `toml:"tags"`         // e.g. ["recon","subdomain"]
	Category    string   `toml:"category"`
	Priority    int      `toml:"priority"` // 1=highest
	Enabled     bool     `toml:"enabled"`
	Version     string   `toml:"version"` // currently installed version
}

// DefaultRegistry is the master list of security tools tracked by ROXX-AUTOKIT
var DefaultRegistry = []ToolEntry{
	// ── RECON / SUBDOMAIN ──────────────────────────────────────────
	{Name: "subfinder", Repo: "projectdiscovery/subfinder",
		Binary: "subfinder", InstallType: "go",
		GoPackage: "github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest",
		Tags: []string{"recon", "subdomain"}, Category: "recon", Priority: 1, Enabled: true},

	{Name: "amass", Repo: "owasp-amass/amass",
		Binary: "amass", InstallType: "go",
		GoPackage: "github.com/owasp-amass/amass/v4/...@master",
		Tags: []string{"recon", "subdomain"}, Category: "recon", Priority: 1, Enabled: true},

	{Name: "assetfinder", Repo: "tomnomnom/assetfinder",
		Binary: "assetfinder", InstallType: "go",
		GoPackage: "github.com/tomnomnom/assetfinder@latest",
		Tags: []string{"recon", "subdomain"}, Category: "recon", Priority: 2, Enabled: true},

	{Name: "findomain", Repo: "findomain/findomain",
		Binary: "findomain", InstallType: "release",
		Tags: []string{"recon", "subdomain"}, Category: "recon", Priority: 2, Enabled: true},

	{Name: "chaos-client", Repo: "projectdiscovery/chaos-client",
		Binary: "chaos", InstallType: "go",
		GoPackage: "github.com/projectdiscovery/chaos-client/cmd/chaos@latest",
		Tags: []string{"recon", "subdomain"}, Category: "recon", Priority: 2, Enabled: true},

	// ── HTTP PROBING ────────────────────────────────────────────────
	{Name: "httpx", Repo: "projectdiscovery/httpx",
		Binary: "httpx", InstallType: "go",
		GoPackage: "github.com/projectdiscovery/httpx/cmd/httpx@latest",
		Tags: []string{"probe", "http"}, Category: "http", Priority: 1, Enabled: true},

	{Name: "httprobe", Repo: "tomnomnom/httprobe",
		Binary: "httprobe", InstallType: "go",
		GoPackage: "github.com/tomnomnom/httprobe@latest",
		Tags: []string{"probe", "http"}, Category: "http", Priority: 2, Enabled: true},

	// ── VULNERABILITY SCANNING ─────────────────────────────────────
	{Name: "nuclei", Repo: "projectdiscovery/nuclei",
		Binary: "nuclei", InstallType: "go",
		GoPackage: "github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest",
		Tags: []string{"vuln", "scanner"}, Category: "scanner", Priority: 1, Enabled: true},

	{Name: "dalfox", Repo: "hahwul/dalfox",
		Binary: "dalfox", InstallType: "go",
		GoPackage: "github.com/hahwul/dalfox/v2@latest",
		Tags: []string{"vuln", "xss"}, Category: "scanner", Priority: 1, Enabled: true},

	{Name: "sqlmap", Repo: "sqlmapproject/sqlmap",
		Binary: "sqlmap", InstallType: "script",
		ScriptURL: "https://raw.githubusercontent.com/sqlmapproject/sqlmap/master/sqlmap.py",
		Tags: []string{"vuln", "sqli"}, Category: "scanner", Priority: 1, Enabled: true},

	// ── FUZZING ─────────────────────────────────────────────────────
	{Name: "ffuf", Repo: "ffuf/ffuf",
		Binary: "ffuf", InstallType: "go",
		GoPackage: "github.com/ffuf/ffuf/v2@latest",
		Tags: []string{"fuzz", "web"}, Category: "fuzzer", Priority: 1, Enabled: true},

	{Name: "feroxbuster", Repo: "epi052/feroxbuster",
		Binary: "feroxbuster", InstallType: "release",
		Tags: []string{"fuzz", "web"}, Category: "fuzzer", Priority: 1, Enabled: true},

	{Name: "gobuster", Repo: "OJ/gobuster",
		Binary: "gobuster", InstallType: "go",
		GoPackage: "github.com/OJ/gobuster/v3@latest",
		Tags: []string{"fuzz", "web"}, Category: "fuzzer", Priority: 2, Enabled: true},

	// ── PARAMETERS / JS ANALYSIS ───────────────────────────────────
	{Name: "gau", Repo: "lc/gau",
		Binary: "gau", InstallType: "go",
		GoPackage: "github.com/lc/gau/v2/cmd/gau@latest",
		Tags: []string{"recon", "urls"}, Category: "recon", Priority: 1, Enabled: true},

	{Name: "waybackurls", Repo: "tomnomnom/waybackurls",
		Binary: "waybackurls", InstallType: "go",
		GoPackage: "github.com/tomnomnom/waybackurls@latest",
		Tags: []string{"recon", "urls"}, Category: "recon", Priority: 2, Enabled: true},

	{Name: "katana", Repo: "projectdiscovery/katana",
		Binary: "katana", InstallType: "go",
		GoPackage: "github.com/projectdiscovery/katana/cmd/katana@latest",
		Tags: []string{"recon", "crawler"}, Category: "recon", Priority: 1, Enabled: true},

	{Name: "paramspider", Repo: "devanshbatham/paramspider",
		Binary: "paramspider", InstallType: "pip",
		PipPackage: "paramspider",
		Tags: []string{"recon", "params"}, Category: "recon", Priority: 2, Enabled: true},

	// ── PORT SCANNING ───────────────────────────────────────────────
	{Name: "naabu", Repo: "projectdiscovery/naabu",
		Binary: "naabu", InstallType: "go",
		GoPackage: "github.com/projectdiscovery/naabu/v2/cmd/naabu@latest",
		Tags: []string{"scan", "ports"}, Category: "scanner", Priority: 1, Enabled: true},

	{Name: "masscan", Repo: "robertdavidgraham/masscan",
		Binary: "masscan", InstallType: "release",
		Tags: []string{"scan", "ports"}, Category: "scanner", Priority: 2, Enabled: true},

	// ── SECRETS / EXPOSURE ─────────────────────────────────────────
	{Name: "trufflehog", Repo: "trufflesecurity/trufflehog",
		Binary: "trufflehog", InstallType: "release",
		Tags: []string{"secrets", "exposure"}, Category: "secrets", Priority: 1, Enabled: true},

	{Name: "gitleaks", Repo: "gitleaks/gitleaks",
		Binary: "gitleaks", InstallType: "release",
		Tags: []string{"secrets", "git"}, Category: "secrets", Priority: 1, Enabled: true},

	{Name: "secretfinder", Repo: "m4ll0k/SecretFinder",
		Binary: "secretfinder", InstallType: "pip",
		PipPackage: "jsbeautifier",
		Tags: []string{"secrets", "js"}, Category: "secrets", Priority: 2, Enabled: true},

	// ── OAST / INTERACTION ─────────────────────────────────────────
	{Name: "interactsh-client", Repo: "projectdiscovery/interactsh",
		Binary: "interactsh-client", InstallType: "go",
		GoPackage: "github.com/projectdiscovery/interactsh/cmd/interactsh-client@latest",
		Tags: []string{"oast", "callback"}, Category: "oast", Priority: 1, Enabled: true},

	// ── CLOUD / MISCONFIG ──────────────────────────────────────────
	{Name: "cloudlist", Repo: "projectdiscovery/cloudlist",
		Binary: "cloudlist", InstallType: "go",
		GoPackage: "github.com/projectdiscovery/cloudlist/cmd/cloudlist@latest",
		Tags: []string{"cloud", "recon"}, Category: "cloud", Priority: 2, Enabled: true},

	{Name: "s3scanner", Repo: "sa7mon/S3Scanner",
		Binary: "s3scanner", InstallType: "pip",
		PipPackage: "s3scanner",
		Tags: []string{"cloud", "aws", "s3"}, Category: "cloud", Priority: 2, Enabled: true},

	// ── UTILITY ────────────────────────────────────────────────────
	{Name: "anew", Repo: "tomnomnom/anew",
		Binary: "anew", InstallType: "go",
		GoPackage: "github.com/tomnomnom/anew@latest",
		Tags: []string{"util"}, Category: "util", Priority: 3, Enabled: true},

	{Name: "qsreplace", Repo: "tomnomnom/qsreplace",
		Binary: "qsreplace", InstallType: "go",
		GoPackage: "github.com/tomnomnom/qsreplace@latest",
		Tags: []string{"util"}, Category: "util", Priority: 3, Enabled: true},

	{Name: "notify", Repo: "projectdiscovery/notify",
		Binary: "notify", InstallType: "go",
		GoPackage: "github.com/projectdiscovery/notify/cmd/notify@latest",
		Tags: []string{"util", "notify"}, Category: "util", Priority: 2, Enabled: true},

	{Name: "dnsx", Repo: "projectdiscovery/dnsx",
		Binary: "dnsx", InstallType: "go",
		GoPackage: "github.com/projectdiscovery/dnsx/cmd/dnsx@latest",
		Tags: []string{"recon", "dns"}, Category: "recon", Priority: 1, Enabled: true},

	{Name: "shuffledns", Repo: "projectdiscovery/shuffledns",
		Binary: "shuffledns", InstallType: "go",
		GoPackage: "github.com/projectdiscovery/shuffledns/cmd/shuffledns@latest",
		Tags: []string{"recon", "dns"}, Category: "recon", Priority: 2, Enabled: true},

	{Name: "puredns", Repo: "d3mondev/puredns",
		Binary: "puredns", InstallType: "go",
		GoPackage: "github.com/d3mondev/puredns/v2@latest",
		Tags: []string{"recon", "dns"}, Category: "recon", Priority: 1, Enabled: true},

	{Name: "nmap", Repo: "nmap/nmap",
		Binary: "nmap", InstallType: "script",
		Tags: []string{"scan", "ports"}, Category: "scanner", Priority: 1, Enabled: true},

	{Name: "jwt_tool", Repo: "ticarpi/jwt_tool",
		Binary: "jwt_tool", InstallType: "pip",
		PipPackage: "termcolor",
		Tags: []string{"vuln", "jwt"}, Category: "scanner", Priority: 2, Enabled: true},
}
