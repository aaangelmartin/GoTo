package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/i18n"
	"github.com/aaangelmartin/goto/internal/urlx"
)

func newEditCmd() *cobra.Command {
	var (
		newURL  string
		newDesc string
		newTags string
		rename  string
	)
	cmd := &cobra.Command{
		Use:   "edit <name>",
		Short: i18n.T("edit_short"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			_, st, err := loadState()
			if err != nil {
				return err
			}
			a, err := st.Get(name)
			if err != nil {
				return err
			}
			if newURL != "" {
				a.URL = urlx.Normalize(newURL, flags.noHTTPS)
			}
			if cmd.Flags().Changed("desc") {
				a.Description = newDesc
			}
			if cmd.Flags().Changed("tag") {
				a.Tags = nil
				for _, t := range strings.Split(newTags, ",") {
					t = strings.TrimSpace(t)
					if t != "" {
						a.Tags = append(a.Tags, t)
					}
				}
			}
			if rename != "" && rename != a.Name {
				if _, err := st.Get(rename); err == nil {
					return fmt.Errorf(i18n.T("edit_exists"), rename)
				}
				if err := st.Delete(a.Name); err != nil {
					return err
				}
				a.Name = rename
			}
			st.Set(a)
			if err := st.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), i18n.T("updated"), a.Name, a.URL)
			return nil
		},
	}
	cmd.Flags().StringVar(&newURL, "url", "", i18n.T("edit_url"))
	cmd.Flags().StringVar(&newDesc, "desc", "", i18n.T("edit_desc"))
	cmd.Flags().StringVar(&newTags, "tag", "", i18n.T("edit_tag"))
	cmd.Flags().StringVar(&rename, "name", "", i18n.T("edit_name"))
	return cmd
}
