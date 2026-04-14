package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/config"
	"github.com/aaangelmartin/goto/internal/i18n"
)

// newOpenerCmd exposes the [openers] table in config.toml as a CLI surface so
// users never have to hand-edit TOML. Keys are either alias types
// ("url", "mailto", "ssh", "file", "directory") or dotted extensions
// (".pdf", ".md"), and values are app names the per-OS opener resolves.
func newOpenerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "opener",
		Short: i18n.T("opener_short"),
	}
	cmd.AddCommand(newOpenerListCmd(), newOpenerSetCmd(), newOpenerUnsetCmd())
	return cmd
}

func newOpenerListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   i18n.T("opener_list_short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _, err := loadState()
			if err != nil {
				return err
			}
			keys := make([]string, 0, len(cfg.Openers))
			for k := range cfg.Openers {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			out := cmd.OutOrStdout()
			for _, k := range keys {
				v := cfg.Openers[k]
				if v == "" {
					v = "(unset)"
				}
				fmt.Fprintf(out, "%-12s  %s\n", k, v)
			}
			return nil
		},
	}
}

func newOpenerSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <type|.ext> <app>",
		Short: i18n.T("opener_set_short"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := normalizeOpenerKey(args[0])
			val := strings.TrimSpace(args[1])
			return persistOpener(key, val, cmd)
		},
	}
}

func newOpenerUnsetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unset <type|.ext>",
		Short: i18n.T("opener_unset_short"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := normalizeOpenerKey(args[0])
			return persistOpener(key, "", cmd)
		},
	}
}

// normalizeOpenerKey lowercases type keys and ensures extensions start with ".".
func normalizeOpenerKey(k string) string {
	k = strings.TrimSpace(k)
	if k == "" {
		return k
	}
	if strings.HasPrefix(k, ".") {
		return strings.ToLower(k)
	}
	// If it looks like an extension without the dot, prepend one.
	if !strings.ContainsAny(k, " /\\") && !isKnownOpenerType(k) && looksLikeExtension(k) {
		return "." + strings.ToLower(k)
	}
	return strings.ToLower(k)
}

func isKnownOpenerType(k string) bool {
	switch strings.ToLower(k) {
	case "url", "mailto", "ssh", "file", "directory":
		return true
	}
	return false
}

func looksLikeExtension(k string) bool {
	// A bare word up to 6 chars is probably an extension (pdf, md, tsx, …).
	if len(k) == 0 || len(k) > 6 {
		return false
	}
	for _, r := range k {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

func persistOpener(key, val string, cmd *cobra.Command) error {
	cfg, _, err := loadState()
	if err != nil {
		return err
	}
	if cfg.Openers == nil {
		cfg.Openers = map[string]string{}
	}
	if val == "" {
		delete(cfg.Openers, key)
	} else {
		cfg.Openers[key] = val
	}
	path, err := config.ConfigPath()
	if err != nil {
		return err
	}
	if err := config.Save(path, cfg); err != nil {
		return err
	}
	if val == "" {
		fmt.Fprintf(cmd.OutOrStdout(), "unset: %s\n", key)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "set: %s -> %s\n", key, val)
	}
	return nil
}
