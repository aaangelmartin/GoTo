package cli

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/i18n"
)

func newExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export [file.json]",
		Short: i18n.T("export_short"),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, st, err := loadState()
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(st.List(), "", "  ")
			if err != nil {
				return err
			}
			b = append(b, '\n')
			if len(args) == 0 {
				_, err := cmd.OutOrStdout().Write(b)
				return err
			}
			return os.WriteFile(args[0], b, 0o644)
		},
	}
}
