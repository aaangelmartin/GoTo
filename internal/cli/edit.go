package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

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
		Short: "Edit an existing alias",
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
					return fmt.Errorf("alias %q already exists", rename)
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
			fmt.Fprintf(cmd.OutOrStdout(), "updated: %s -> %s\n", a.Name, a.URL)
			return nil
		},
	}
	cmd.Flags().StringVar(&newURL, "url", "", "new URL")
	cmd.Flags().StringVar(&newDesc, "desc", "", "new description")
	cmd.Flags().StringVar(&newTags, "tag", "", "replace tags (comma-separated)")
	cmd.Flags().StringVar(&rename, "name", "", "rename the alias")
	return cmd
}
