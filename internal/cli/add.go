package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/alias"
	"github.com/aaangelmartin/goto/internal/i18n"
	"github.com/aaangelmartin/goto/internal/urlx"
)

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name> <url>",
		Short: i18n.T("add_short"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, rawURL := args[0], args[1]
			if strings.ContainsAny(name, " \t/\\") {
				return fmt.Errorf(i18n.T("err_name_invalid"), name)
			}
			_, st, err := loadState()
			if err != nil {
				return err
			}
			a := alias.Alias{
				Name:        name,
				URL:         urlx.Normalize(rawURL, flags.noHTTPS),
				Description: flags.descFlag,
				CreatedAt:   time.Now(),
			}
			if flags.tagFlag != "" {
				for _, t := range strings.Split(flags.tagFlag, ",") {
					t = strings.TrimSpace(t)
					if t != "" {
						a.Tags = append(a.Tags, t)
					}
				}
			}
			if err := st.Put(a); err != nil {
				return err
			}
			if err := st.Save(); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), i18n.T("added"), a.Name, a.URL)
			return nil
		},
	}
	cmd.Flags().StringVar(&flags.tagFlag, "tag", "", i18n.T("add_tag"))
	cmd.Flags().StringVar(&flags.descFlag, "desc", "", i18n.T("add_desc"))
	return cmd
}
