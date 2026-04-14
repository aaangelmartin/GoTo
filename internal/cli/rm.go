package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newRmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm <name>",
		Aliases: []string{"remove", "delete", "del"},
		Short:   "Remove an alias",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cfg, st, err := loadState()
			if err != nil {
				return err
			}
			a, err := st.Get(name)
			if err != nil {
				return err
			}
			if cfg.ConfirmDelete && !flags.yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Delete %q (%s)? [y/N] ", a.Name, a.URL)
				reader := bufio.NewReader(os.Stdin)
				resp, _ := reader.ReadString('\n')
				resp = strings.TrimSpace(strings.ToLower(resp))
				if resp != "y" && resp != "yes" {
					fmt.Fprintln(cmd.OutOrStdout(), "aborted")
					return nil
				}
			}
			if err := st.Delete(name); err != nil {
				return err
			}
			if err := st.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "removed: %s\n", a.Name)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&flags.yes, "yes", "y", false, "skip confirmation")
	return cmd
}
