package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/config"
	"github.com/aaangelmartin/goto/internal/i18n"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: i18n.T("config_short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.ConfigPath()
			if err != nil {
				return err
			}
			if _, err := os.Stat(path); os.IsNotExist(err) {
				if err := config.Save(path, config.Default()); err != nil {
					return err
				}
			}
			editor := os.Getenv("EDITOR")
			if editor == "" {
				fmt.Fprintln(cmd.OutOrStdout(), path)
				return nil
			}
			c := exec.Command(editor, path)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			return c.Run()
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: i18n.T("config_path_short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.ConfigPath()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), path)
			return nil
		},
	})
	return cmd
}
