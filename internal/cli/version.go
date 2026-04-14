package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/buildinfo"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version, commit and build date",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "goto %s (commit %s, built %s)\n",
				buildinfo.Version, buildinfo.Commit, buildinfo.Date)
			return nil
		},
	}
}
