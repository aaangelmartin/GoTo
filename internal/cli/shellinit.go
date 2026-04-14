package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/i18n"
)

// newShellInitCmd emits a shell wrapper that lets `goto <dir-alias>` actually
// change the caller's shell working directory (a child process can't `cd`
// for its parent, so we eval a script the wrapper captured).
func newShellInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:       "shell-init [bash|zsh|fish]",
		Short:     i18n.T("shellinit_short"),
		ValidArgs: []string{"bash", "zsh", "fish"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			switch args[0] {
			case "bash", "zsh":
				fmt.Fprint(out, bashZshWrapper)
			case "fish":
				fmt.Fprint(out, fishWrapper)
			}
			return nil
		},
	}
}

const bashZshWrapper = `# goto shell wrapper (bash/zsh). Source this from ~/.bashrc or ~/.zshrc:
#   eval "$(command goto shell-init zsh)"
goto() {
    local out
    if [ $# -eq 0 ]; then
        command goto
        return $?
    fi
    # Invoke the binary with an internal flag so it emits a shell-evalable
    # script when the target is a directory; otherwise it opens normally.
    out=$(command goto --_shell-exec "$@")
    local code=$?
    if [ -n "$out" ]; then
        eval "$out"
    fi
    return $code
}
`

const fishWrapper = `# goto shell wrapper (fish). Add this to ~/.config/fish/config.fish:
#   goto shell-init fish | source
function goto
    if test (count $argv) -eq 0
        command goto
        return $status
    end
    set -l out (command goto --_shell-exec $argv)
    set -l code $status
    if test -n "$out"
        eval $out
    end
    return $code
end
`
