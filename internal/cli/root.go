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
	"github.com/aaangelmartin/goto/internal/opener"
	"github.com/aaangelmartin/goto/internal/store"
	"github.com/aaangelmartin/goto/internal/tui"
	"github.com/aaangelmartin/goto/internal/urlx"
)

type globalFlags struct {
	browser   string
	noHTTPS   bool
	dryRun    bool
	yes       bool
	useJSON   bool
	tagFlag   string
	descFlag  string
	typeFlag  string
	lang      string
	shellExec bool
}

func parseType(s string) alias.Type {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "auto":
		return alias.TypeAuto
	case "url", "web", "link":
		return alias.TypeURL
	case "mail", "mailto", "email":
		return alias.TypeMailto
	case "ssh":
		return alias.TypeSSH
	case "file":
		return alias.TypeFile
	case "dir", "directory", "folder":
		return alias.TypeDirectory
	case "cmd", "command":
		return alias.TypeCommand
	}
	return alias.TypeAuto
}

var flags globalFlags

// Execute runs the root command.
func Execute() error {
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
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completeAliasNames(toComplete), cobra.ShellCompDirectiveNoFileComp
		},
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
	// Internal flag used by the shell wrapper; hidden from help.
	pf.BoolVar(&flags.shellExec, "_shell-exec", false, "internal: emit shell eval output for directory aliases")
	_ = pf.MarkHidden("_shell-exec")

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
		newShellInitCmd(),
		newVersionCmd(),
	)
	return root
}

func completeAliasNames(prefix string) []string {
	_, st, err := loadState()
	if err != nil {
		return nil
	}
	var out []string
	for _, a := range st.List() {
		if prefix == "" || strings.HasPrefix(a.Name, prefix) {
			out = append(out, a.Name)
		}
	}
	return out
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
	return resolveAndDispatch(target, cfg, st, cmd)
}

func resolveAndDispatch(target string, cfg config.Config, st *store.Store, cmd *cobra.Command) error {
	target = strings.TrimSpace(target)
	if target == "" {
		return fmt.Errorf("%s", i18n.T("err_empty_target"))
	}

	// 1. exact alias match (preferred over protocol guessing so an alias
	//    named "ssh" can still win even if it looks like a URL).
	if a, err := st.Get(target); err == nil {
		return openAliasWithHit(a, cfg, st, cmd)
	}

	// 2. fuzzy alias match.
	matches := alias.Rank(target, st.List(), cfg.FuzzyThreshold)
	switch {
	case len(matches) == 1:
		return openAliasWithHit(matches[0].Alias, cfg, st, cmd)
	case len(matches) > 1 && matches[0].Score >= 0.9 && matches[1].Score < 0.75:
		return openAliasWithHit(matches[0].Alias, cfg, st, cmd)
	case len(matches) > 1:
		return fmt.Errorf(i18n.T("err_ambiguous"), target, candidateNames(matches))
	}

	// 3. no alias matched — treat the raw argument as a target.
	adhoc := alias.Alias{Target: target, Type: alias.TypeAuto}
	t := alias.Detect(target)
	// Apply default_action override: if user forces "url", skip path detection.
	switch cfg.DefaultAction {
	case "url", "link":
		t = alias.TypeURL
	case "file":
		t = alias.TypeFile
	case "directory":
		t = alias.TypeDirectory
	}
	if t == alias.TypeURL {
		adhoc.Target = urlx.Normalize(target, flags.noHTTPS)
	}
	return dispatchOpen(adhoc, t, cfg, cmd)
}

func openAliasWithHit(a alias.Alias, cfg config.Config, st *store.Store, cmd *cobra.Command) error {
	t := alias.Resolve(a)
	if err := dispatchOpen(a, t, cfg, cmd); err != nil {
		return err
	}
	bumpHit(st, a)
	return nil
}

func dispatchOpen(a alias.Alias, t alias.Type, cfg config.Config, cmd *cobra.Command) error {
	if flags.dryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s\n", t, a.Target)
		return nil
	}
	ocfg := opener.Config{
		Browser:       firstNonEmpty(flags.browser, cfg.Browser),
		Openers:       cfg.Openers,
		Terminal:      cfg.Terminal,
		DirectoryMode: cfg.DirectoryMode,
	}
	res, err := opener.Open(a, t, ocfg)
	if err != nil {
		return err
	}
	if res.ShellScript != "" {
		if flags.shellExec {
			fmt.Fprintln(cmd.OutOrStdout(), res.ShellScript)
			return nil
		}
		// Not running inside the wrapper: print a hint so users install it.
		fmt.Fprintln(cmd.ErrOrStderr(), i18n.T("err_shell_wrapper_missing"))
		fmt.Fprintln(cmd.OutOrStdout(), res.ShellScript)
	}
	return nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func candidateNames(matches []alias.Match) string {
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, m.Alias.Name)
	}
	return strings.Join(out, ", ")
}

func bumpHit(st *store.Store, a alias.Alias) {
	a.HitCount++
	a.LastOpened = time.Now()
	st.Set(a)
	_ = st.Save()
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
	// Honor a persisted language choice only when the caller did not already
	// override via --lang (flags.lang) or env (pre-parse in Execute).
	if cfg.Language != "" && cfg.Language != "auto" && flags.lang == "" && os.Getenv("GOTO_LANG") == "" {
		i18n.SetLang(cfg.Language)
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
