//go:build darwin

package opener

import (
	"os/exec"
	"strings"
)

var macAppAliases = map[string]string{
	"chrome":  "Google Chrome",
	"firefox": "Firefox",
	"safari":  "Safari",
	"arc":     "Arc",
	"brave":   "Brave Browser",
	"edge":    "Microsoft Edge",
	"cursor":  "Cursor",
	"code":    "Visual Studio Code",
	"vscode":  "Visual Studio Code",
	"finder":  "Finder",
	"preview": "Preview",
	"mail":    "Mail",
}

func openFile(path, app string) error {
	if app == "" || app == "default" {
		return exec.Command("open", path).Start()
	}
	resolved := macAppAliases[strings.ToLower(app)]
	if resolved == "" {
		resolved = app
	}
	return exec.Command("open", "-a", resolved, path).Start()
}

func openSSH(target string, cfg Config) error {
	host := strings.TrimPrefix(target, "ssh://")
	term := cfg.Terminal
	if term == "" {
		term = "Terminal" // macOS system terminal
	}
	// Build AppleScript that opens a new tab running ssh.
	script := "tell application \"" + term + "\" to do script \"ssh " + host + "\""
	return exec.Command("osascript", "-e", script, "-e", "tell application \""+term+"\" to activate").Start()
}
