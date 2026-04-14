// Package cli wires all CLI subcommands using Cobra.
package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/aaangelmartin/goto/internal/alias"
	"github.com/aaangelmartin/goto/internal/config"
	"github.com/aaangelmartin/goto/internal/store"
	"github.com/aaangelmartin/goto/internal/tui"
	"github.com/aaangelmartin/goto/internal/urlx"
)

type globalFlags struct {
	browser  string
	noHTTPS  bool
	dryRun   bool
	yes      bool
	useJSON  bool
	tagFlag  string
	descFlag string
}

var flags globalFlags

// Execute runs the root command.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "goto [target]",
		Short: "Open URLs and manage link aliases from the terminal",
		Long: `goto opens URLs in your browser (auto-prepends https:// if missing)
and lets you manage personal link aliases with a beautiful TUI.

  goto google.com         Opens https://google.com
  goto gh                 Resolves alias "gh" and opens it
  goto                    Launches the interactive TUI
  goto add gh github.com  Adds an alias`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ArbitraryArgs,
		RunE:          runDefault,
	}

	pf := root.PersistentFlags()
	pf.StringVar(&flags.browser, "browser", "", "browser to use (default|chrome|firefox|safari|arc|brave|edge)")
	pf.BoolVar(&flags.noHTTPS, "no-https", false, "use http:// instead of https:// when prepending protocol")
	pf.BoolVar(&flags.dryRun, "dry-run", false, "print resolved URL without opening it")

	root.AddCommand(
		newAddCmd(),
		newRmCmd(),
		newLsCmd(),
		newEditCmd(),
		newSearchCmd(),
		newImportCmd(),
		newExportCmd(),
		newConfigCmd(),
		newCompletionCmd(),
		newVersionCmd(),
	)
	return root
}

// runDefault is the action when `goto` is invoked without a subcommand.
func runDefault(cmd *cobra.Command, args []string) error {
	cfg, st, err := loadState()
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return tui.Run(st, cfg)
	}
	target := strings.Join(args, " ")
	return resolveAndOpen(target, cfg, st, cmd)
}

func resolveAndOpen(target string, cfg config.Config, st *store.Store, cmd *cobra.Command) error {
	target = strings.TrimSpace(target)
	if target == "" {
		return fmt.Errorf("empty target")
	}

	// 1. explicit URL with protocol
	if hasProtocol(target) {
		return openURL(target, cfg, cmd)
	}

	// 2. exact alias
	if a, err := st.Get(target); err == nil {
		bumpHit(st, a)
		return openURL(a.URL, cfg, cmd)
	}

	// 3. fuzzy
	matches := alias.Rank(target, st.List(), cfg.FuzzyThreshold)
	switch {
	case len(matches) == 1:
		bumpHit(st, matches[0].Alias)
		return openURL(matches[0].Alias.URL, cfg, cmd)
	case len(matches) > 1 && matches[0].Score >= 0.9 && (len(matches) < 2 || matches[1].Score < 0.75):
		// strong clear winner
		bumpHit(st, matches[0].Alias)
		return openURL(matches[0].Alias.URL, cfg, cmd)
	case len(matches) > 1:
		return fmt.Errorf("ambiguous target %q; candidates: %s", target, candidateNames(matches))
	}

	// 4. looks like a URL -> normalize
	if urlx.LooksLikeURL(target) {
		return openURL(urlx.Normalize(target, flags.noHTTPS), cfg, cmd)
	}

	return fmt.Errorf("no alias found and %q is not a URL; try: goto search %q", target, target)
}

func candidateNames(matches []alias.Match) string {
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, m.Alias.Name)
	}
	return strings.Join(out, ", ")
}

func openURL(url string, cfg config.Config, cmd *cobra.Command) error {
	if flags.dryRun {
		fmt.Fprintln(cmd.OutOrStdout(), url)
		return nil
	}
	browser := flags.browser
	if browser == "" {
		browser = cfg.Browser
	}
	return urlx.Open(url, browser)
}

func bumpHit(st *store.Store, a alias.Alias) {
	a.HitCount++
	a.LastOpened = time.Now()
	st.Set(a)
	_ = st.Save() // best-effort; stderr warnings left to callers
}

func hasProtocol(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ':' && i > 0 {
			return true
		}
		isLower := c >= 'a' && c <= 'z'
		isUpper := c >= 'A' && c <= 'Z'
		isDigit := c >= '0' && c <= '9'
		isSym := c == '+' || c == '-' || c == '.'
		if !isLower && !isUpper && !isDigit && !isSym {
			return false
		}
	}
	return false
}

// loadState loads config + store, returning them ready to use.
func loadState() (config.Config, *store.Store, error) {
	cfgPath, err := config.ConfigPath()
	if err != nil {
		return config.Config{}, nil, err
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "goto: warning loading config: %v\n", err)
	}
	aliasPath, err := config.AliasesPath()
	if err != nil {
		return cfg, nil, err
	}
	st := store.New(aliasPath)
	if err := st.Load(); err != nil {
		return cfg, nil, fmt.Errorf("load aliases: %w", err)
	}
	return cfg, st, nil
}
