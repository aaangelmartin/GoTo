package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/i18n"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 i18n.T("completion_short"),
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			return emitCompletion(cmd, args[0], os.Stdout)
		},
	}
	cmd.AddCommand(newCompletionInstallCmd())
	return cmd
}

// emitCompletion writes the shell-specific completion script to w.
func emitCompletion(cmd *cobra.Command, shell string, w interface {
	Write(p []byte) (n int, err error)
}) error {
	root := cmd.Root()
	switch shell {
	case "bash":
		return root.GenBashCompletionV2(w, true)
	case "zsh":
		return root.GenZshCompletion(w)
	case "fish":
		return root.GenFishCompletion(w, true)
	case "powershell":
		return root.GenPowerShellCompletionWithDesc(w)
	}
	return fmt.Errorf("unknown shell: %s", shell)
}

// newCompletionInstallCmd writes the script to the standard place for the
// detected (or requested) shell and tells the user what to do next.
func newCompletionInstallCmd() *cobra.Command {
	var shellFlag string
	cmd := &cobra.Command{
		Use:   "install [--shell bash|zsh|fish]",
		Short: i18n.T("completion_install_short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := shellFlag
			if shell == "" {
				shell = detectShell()
			}
			if shell == "" {
				return fmt.Errorf("could not detect shell; pass --shell bash|zsh|fish")
			}
			target, err := completionTarget(shell)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			var buf bytes.Buffer
			if err := emitCompletion(cmd, shell, &buf); err != nil {
				return err
			}
			if err := os.WriteFile(target, buf.Bytes(), 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "✓ %s completion installed: %s\n", shell, target)
			fmt.Fprintln(cmd.OutOrStdout(), activationHint(shell, target))
			return nil
		},
	}
	cmd.Flags().StringVar(&shellFlag, "shell", "", "target shell (bash|zsh|fish); defaults to $SHELL")
	return cmd
}

func detectShell() string {
	s := os.Getenv("SHELL")
	base := filepath.Base(s)
	switch base {
	case "bash", "zsh", "fish":
		return base
	}
	return ""
}

// completionTarget returns the conventional install path per shell. macOS
// with Homebrew keeps zsh site-functions under /opt/homebrew/share; we aim
// for a user-writable dir so no sudo is needed.
func completionTarget(shell string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch shell {
	case "zsh":
		// Users must have this dir in their fpath before `compinit`.
		return filepath.Join(home, ".zsh", "completions", "_goto"), nil
	case "bash":
		// Sourced by bash-completion; some systems prefer ~/.local/share/bash-completion/completions.
		return filepath.Join(home, ".local", "share", "bash-completion", "completions", "goto"), nil
	case "fish":
		return filepath.Join(home, ".config", "fish", "completions", "goto.fish"), nil
	}
	return "", fmt.Errorf("unsupported shell: %s", shell)
}

func activationHint(shell, target string) string {
	switch shell {
	case "zsh":
		dir := filepath.Dir(target)
		return strings.Join([]string{
			"",
			"To activate, add this to ~/.zshrc (once) and start a new shell:",
			"",
			fmt.Sprintf("    fpath=(%s $fpath)", dir),
			"    autoload -Uz compinit && compinit",
			"",
			"Then: goto <TAB> completes your alias names.",
		}, "\n")
	case "bash":
		return strings.Join([]string{
			"",
			"If bash-completion is installed, the script is picked up on the next login.",
			"Otherwise source it from ~/.bashrc:",
			"",
			fmt.Sprintf("    source %s", target),
			"",
		}, "\n")
	case "fish":
		return "\nFish picks it up automatically on the next shell; no extra config needed.\n"
	}
	return ""
}
