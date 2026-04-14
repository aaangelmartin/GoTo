// Package opener dispatches aliases to the appropriate action based on type.
//
// It is the single source of truth for "what happens when you press enter on
// X?": web browser for URLs, system default for files, mail client for
// mailto, terminal + ssh for SSH hosts, and a shell "cd" directive for
// directories (emitted as text when cmd output is consumed by the `goto`
// shell wrapper installed via `goto shell-init`).
package opener

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aaangelmartin/goto/internal/alias"
	"github.com/aaangelmartin/goto/internal/urlx"
)

// Config collects the runtime knobs the opener needs. It's populated from the
// user's config.toml but kept separate to avoid an import cycle.
type Config struct {
	// Browser is the default browser for URLs.
	Browser string
	// Openers maps a type (or ".ext") to a platform-specific app name or path.
	// Keys: "url", "mailto", "ssh", "file", "directory", ".pdf", ".md", ...
	Openers map[string]string
	// Terminal is the terminal emulator used to open SSH sessions.
	Terminal string
	// DirectoryMode controls how directories open: "shell" (emit cd) or "finder".
	DirectoryMode string
}

// Result represents the outcome of Open.
type Result struct {
	// ShellScript, when non-empty, should be printed to stdout so that the
	// `goto` shell wrapper can eval it (used for cd and similar shell-side
	// side effects that a child process can't do on behalf of its parent).
	ShellScript string
}

// Open dispatches an alias to its concrete handler. t is the already-resolved
// type (never TypeAuto).
func Open(a alias.Alias, t alias.Type, cfg Config) (Result, error) {
	switch t {
	case alias.TypeURL:
		return Result{}, openURL(a.Target, cfg)
	case alias.TypeMailto:
		target := a.Target
		if !strings.HasPrefix(target, "mailto:") {
			target = "mailto:" + target
		}
		return Result{}, openWithBrowser(target, cfg.pick("mailto"))
	case alias.TypeSSH:
		return Result{}, openSSH(a.Target, cfg)
	case alias.TypeFile:
		path := stripFilePrefix(alias.ExpandHome(a.Target))
		app := cfg.pick("file")
		if ext := strings.ToLower(filepath.Ext(path)); ext != "" {
			if byExt := cfg.Openers["."+strings.TrimPrefix(ext, ".")]; byExt != "" {
				app = byExt
			} else if byExt := cfg.Openers[ext]; byExt != "" {
				app = byExt
			}
		}
		return Result{}, openFile(path, app)
	case alias.TypeDirectory:
		path := stripFilePrefix(alias.ExpandHome(a.Target))
		mode := cfg.DirectoryMode
		if mode == "" {
			mode = "shell"
		}
		if mode == "shell" {
			// Emit a cd directive the parent shell wrapper will eval.
			return Result{ShellScript: fmt.Sprintf("cd %s", shellQuote(path))}, nil
		}
		return Result{}, openFile(path, cfg.pick("directory"))
	case alias.TypeCommand:
		return Result{}, runCommand(a.Target)
	}
	return Result{}, fmt.Errorf("unsupported alias type: %s", t)
}

func (c Config) pick(kind string) string {
	if v := c.Openers[kind]; v != "" {
		return v
	}
	if kind == "url" && c.Browser != "" {
		return c.Browser
	}
	return "default"
}

func openURL(target string, cfg Config) error {
	u := urlx.Normalize(target, false)
	return urlx.Open(u, cfg.pick("url"))
}

// openWithBrowser delegates to the URL opener so mailto: links reach the
// OS-registered mail client through the same pathway as web URLs.
func openWithBrowser(target, browser string) error {
	return urlx.Open(target, browser)
}

func stripFilePrefix(t string) string {
	return strings.TrimPrefix(t, "file://")
}

// shellQuote wraps a path in single quotes, escaping any embedded quotes for
// safe eval by the user's shell.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func runCommand(cmdline string) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
		if runtime.GOOS == "windows" {
			shell = "cmd"
		}
	}
	args := []string{"-c", cmdline}
	if runtime.GOOS == "windows" {
		args = []string{"/c", cmdline}
	}
	c := exec.Command(shell, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
