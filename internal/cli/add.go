package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/alias"
	"github.com/aaangelmartin/goto/internal/urlx"
)

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name> <url>",
		Short: "Add a new alias",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, rawURL := args[0], args[1]
			if strings.ContainsAny(name, " \t/\\") {
				return fmt.Errorf("invalid alias name %q: no whitespace or slashes", name)
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
			fmt.Fprintf(cmd.OutOrStdout(), "added: %s -> %s\n", a.Name, a.URL)
			return nil
		},
	}
	cmd.Flags().StringVar(&flags.tagFlag, "tag", "", "comma-separated tags")
	cmd.Flags().StringVar(&flags.descFlag, "desc", "", "description")
	return cmd
}
