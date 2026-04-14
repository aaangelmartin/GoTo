package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/buildinfo"
	"github.com/aaangelmartin/goto/internal/i18n"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: i18n.T("version_short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), i18n.T("version_line"),
				buildinfo.Version, buildinfo.Commit, buildinfo.Date)
			return nil
		},
	}
}
