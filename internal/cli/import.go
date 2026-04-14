package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/alias"
	"github.com/aaangelmartin/goto/internal/i18n"
	"github.com/aaangelmartin/goto/internal/store"
)

func newImportCmd() *cobra.Command {
	var overwrite bool
	cmd := &cobra.Command{
		Use:   "import <file.json>",
		Short: i18n.T("import_short"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			var list []alias.Alias
			if err := json.Unmarshal(b, &list); err != nil {
				return fmt.Errorf(i18n.T("err_parse"), args[0], err)
			}
			_, st, err := loadState()
			if err != nil {
				return err
			}
			added, updated, skipped := 0, 0, 0
			for _, a := range list {
				if _, err := st.Get(a.Name); err == nil {
					if overwrite {
						st.Set(a)
						updated++
					} else {
						skipped++
					}
					continue
				} else if !errors.Is(err, store.ErrNotFound) {
					return err
				}
				if err := st.Put(a); err != nil {
					return err
				}
				added++
			}
			if err := st.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), i18n.T("imported"), added, updated, skipped)
			return nil
		},
	}
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, i18n.T("import_overwrite"))
	return cmd
}
