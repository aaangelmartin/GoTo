//go:build windows

package opener

import (
	"os/exec"
	"strings"
)

func openFile(path, app string) error {
	if app == "" || app == "default" {
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", path).Start()
	}
	return exec.Command("cmd", "/c", "start", "", app, path).Start()
}

func openSSH(target string, cfg Config) error {
	host := strings.TrimPrefix(target, "ssh://")
	term := cfg.Terminal
	if term == "" {
		term = "wt" // Windows Terminal if available
	}
	return exec.Command(term, "ssh", host).Start()
}
