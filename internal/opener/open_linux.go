//go:build linux

package opener

import (
	"os/exec"
	"strings"
)

func openFile(path, app string) error {
	if app == "" || app == "default" {
		return exec.Command("xdg-open", path).Start()
	}
	if _, err := exec.LookPath(app); err != nil {
		return exec.Command("xdg-open", path).Start()
	}
	return exec.Command(app, path).Start()
}

func openSSH(target string, cfg Config) error {
	host := strings.TrimPrefix(target, "ssh://")
	term := cfg.Terminal
	if term == "" {
		// try some common terminals in order
		for _, t := range []string{"x-terminal-emulator", "gnome-terminal", "konsole", "xterm"} {
			if _, err := exec.LookPath(t); err == nil {
				term = t
				break
			}
		}
	}
	if term == "" {
		return exec.Command("ssh", host).Start()
	}
	return exec.Command(term, "-e", "ssh", host).Start()
}
