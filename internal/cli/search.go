package cli

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/alias"
	"github.com/aaangelmartin/goto/internal/i18n"
)

func newSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query...>",
		Short: i18n.T("search_short"),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _, err := loadState()
			if err != nil {
				return err
			}
			query := strings.Join(args, " ")
			tmpl := cfg.SearchEngine
			if !strings.Contains(tmpl, "{q}") {
				return fmt.Errorf("%s", i18n.T("err_search_config"))
			}
			u := strings.ReplaceAll(tmpl, "{q}", url.QueryEscape(query))
			return dispatchOpen(alias.Alias{Target: u}, alias.TypeURL, cfg, cmd)
		},
	}
}
