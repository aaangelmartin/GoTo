package cli

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query...>",
		Short: "Open a web search for the given query",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _, err := loadState()
			if err != nil {
				return err
			}
			query := strings.Join(args, " ")
			tmpl := cfg.SearchEngine
			if !strings.Contains(tmpl, "{q}") {
				return fmt.Errorf("search_engine in config must contain {q}")
			}
			u := strings.ReplaceAll(tmpl, "{q}", url.QueryEscape(query))
			return openURL(u, cfg, cmd)
		},
	}
}
