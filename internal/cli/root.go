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
	"github.com/aaangelmartin/goto/internal/i18n"
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
	lang     string
}

var flags globalFlags

// Execute runs the root command.
func Execute() error {
	// Pre-parse the --lang / -lang flag so Cobra help text (built at command
	// creation time) is rendered in the requested language.
	preparseLang(os.Args[1:])
	return newRootCmd().Execute()
}

func preparseLang(args []string) {
	for i, a := range args {
		switch {
		case a == "--lang" || a == "-lang":
			if i+1 < len(args) {
				i18n.SetLang(args[i+1])
				return
			}
		case strings.HasPrefix(a, "--lang="):
			i18n.SetLang(strings.TrimPrefix(a, "--lang="))
			return
		case strings.HasPrefix(a, "-lang="):
			i18n.SetLang(strings.TrimPrefix(a, "-lang="))
			return
		}
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "goto [target]",
		Short:         i18n.T("short"),
		Long:          i18n.T("long"),
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ArbitraryArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flags.lang != "" {
				i18n.SetLang(flags.lang)
			}
		},
		RunE: runDefault,
	}

	pf := root.PersistentFlags()
	pf.StringVar(&flags.browser, "browser", "", i18n.T("flag_browser"))
	pf.BoolVar(&flags.noHTTPS, "no-https", false, i18n.T("flag_nohttps"))
	pf.BoolVar(&flags.dryRun, "dry-run", false, i18n.T("flag_dryrun"))
	pf.StringVar(&flags.lang, "lang", "", i18n.T("flag_lang"))

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
		return fmt.Errorf("%s", i18n.T("err_empty_target"))
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
		bumpHit(st, matches[0].Alias)
		return openURL(matches[0].Alias.URL, cfg, cmd)
	case len(matches) > 1:
		return fmt.Errorf(i18n.T("err_ambiguous"), target, candidateNames(matches))
	}

	// 4. looks like a URL -> normalize
	if urlx.LooksLikeURL(target) {
		return openURL(urlx.Normalize(target, flags.noHTTPS), cfg, cmd)
	}

	return fmt.Errorf(i18n.T("err_notfound"), target, target)
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
	_ = st.Save()
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
		fmt.Fprintf(os.Stderr, i18n.T("warn_config_load"), err)
	}
	aliasPath, err := config.AliasesPath()
	if err != nil {
		return cfg, nil, err
	}
	st := store.New(aliasPath)
	if err := st.Load(); err != nil {
		return cfg, nil, fmt.Errorf(i18n.T("err_load_aliases"), err)
	}
	return cfg, st, nil
}
