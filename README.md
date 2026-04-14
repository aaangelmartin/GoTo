<div align="center">

```
   ___  ___  ____  ____
  / __)/ _ \(_  _)/ _  )
 ( (_-\(_)  ) )(  )(_) )
  \___/\___/ (__) \____)
```

### Open any URL and manage link aliases straight from your terminal.

[![CI](https://github.com/aaangelmartin/GoTo/actions/workflows/ci.yml/badge.svg)](https://github.com/aaangelmartin/GoTo/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/aaangelmartin/GoTo?sort=semver&color=%23FF79C6)](https://github.com/aaangelmartin/GoTo/releases/latest)
[![License](https://img.shields.io/badge/license-Apache%202.0-8BE9FD)](./LICENSE)
[![Go Report](https://goreportcard.com/badge/github.com/aaangelmartin/goto)](https://goreportcard.com/report/github.com/aaangelmartin/goto)
[![Made with Charm](https://img.shields.io/badge/made%20with-%F0%9F%92%9C%20Charm-FF79C6)](https://charm.sh)

**English** · [Español](./README.es.md)

</div>

---

## Why `goto`?

`open google.com` fails on macOS because `open(1)` thinks `google.com` is a file.
You end up typing `open https://google.com` every time. Not anymore.

```sh
$ goto google.com           # opens https://google.com
$ goto gh                   # opens your saved alias 'gh'
$ goto                      # launches the interactive TUI
```

- **Zero-friction URLs** — auto-prepends `https://` when missing; respects `http://`, `mailto:`, `ssh://`, `file://`, and any other protocol.
- **Named aliases** — `gh`, `mail`, `jira`, `localhost3000`. Keep your muscle memory, stop typing URLs.
- **Fuzzy matching** — `goto githu` still hits `github`.
- **Beautiful TUI** — [Bubble Tea](https://github.com/charmbracelet/bubbletea) powered. Add, edit, delete, search, copy, filter by tag. Four built-in themes.
- **Bilingual** — English and Spanish interface auto-detected from `LANG`, overridable with `--lang`.
- **Cross-platform** — macOS, Linux, Windows. Single static binary, ~5 MB.
- **Private by default** — everything lives in `~/.local/share/goto/` and `~/.config/goto/`. No telemetry.

## Install

### Homebrew (macOS / Linux)

```sh
brew install aaangelmartin/tap/goto
```

### Go

```sh
go install github.com/aaangelmartin/goto/cmd/goto@latest
```

### Manual

Download a binary for your platform from the [latest release](https://github.com/aaangelmartin/GoTo/releases/latest), extract, and drop `goto` into your `$PATH`.

## Usage

### Open things

```sh
goto google.com                       # https://google.com
goto localhost:3000 --no-https        # http://localhost:3000
goto https://claude.ai                # passthrough
goto search "claude code"             # opens your configured search engine
```

### Manage aliases

```sh
goto add gh github.com/aaangelmartin --tag dev --desc "My GitHub"
goto add jira my-co.atlassian.net/jira --tag work
goto ls                               # pretty table
goto ls --tag work --json             # machine-readable
goto edit gh --url github.com/aaangelmartin/GoTo
goto rm jira -y                       # skip confirmation
goto export > aliases.json            # backup
goto import aliases.json --overwrite  # restore
```

### Launch the TUI

```sh
goto
```

<div align="center">

| key | action |
| --- | --- |
| `↑`/`k`, `↓`/`j` | move selection |
| `enter` | open in browser |
| `/` | live filter |
| `a` / `e` / `d` | add / edit / delete |
| `t` | toggle tag filter on selection |
| `y` | copy URL to clipboard |
| `?` | keyboard help |
| `q` | quit |

</div>

### Language

The interface follows `$LANG` automatically (any `es*` locale → Spanish, else English).
You can force a language with either:

```sh
goto --lang es ls        # one-off
export GOTO_LANG=es      # session-wide
```

### Shell completion

```sh
goto completion zsh  > "${fpath[1]}/_goto"   # zsh
goto completion bash > /etc/bash_completion.d/goto
goto completion fish > ~/.config/fish/completions/goto.fish
```

## Configuration

Config lives at `$XDG_CONFIG_HOME/goto/config.toml` (defaults to `~/.config/goto/config.toml`). Aliases live at `$XDG_DATA_HOME/goto/aliases.json`.

```toml
browser         = "default"          # default | chrome | firefox | safari | arc | brave | edge | opera | vivaldi
search_engine   = "https://www.google.com/search?q={q}"   # use DuckDuckGo, Kagi, etc.
theme           = "default"          # default | dracula | catppuccin | nord | tokyonight
confirm_delete  = true
fuzzy_threshold = 0.4                # lower = more lenient fuzzy matches
```

Override paths with environment variables:

```sh
export GOTO_CONFIG=$HOME/.dotfiles/goto.toml
export GOTO_ALIASES=$HOME/.dotfiles/aliases.json
export GOTO_LANG=es                  # or en
```

## Flags

| flag | description |
| --- | --- |
| `--browser <name>` | override browser for this call |
| `--no-https` | prepend `http://` instead of `https://` when no protocol is given |
| `--dry-run` | print the resolved URL without opening it |
| `--lang <en\|es>` | force interface language for this invocation |

## How resolution works

When you run `goto foo`, `goto` tries these in order:

1. **Protocol present?** (`foo://bar`, `mailto:a@b`) — open as-is.
2. **Exact alias match?** — open, bump hit counter.
3. **Fuzzy alias match** — if exactly one candidate (or a strong winner above `fuzzy_threshold`), open it.
4. **Looks like a URL?** (contains a dot, no whitespace) — normalize (`https://` prefix) and open.
5. **Nothing matched** — suggest `goto search "foo"`.

## Roadmap

- [ ] Sync backend (gist / git repo) for cross-device aliases
- [ ] Group / workspace filtering
- [ ] History view in TUI
- [ ] `goto recent` for MRU
- [ ] Browser profile selection (Chrome `--profile`)
- [ ] Scoop bucket for Windows
- [ ] Interactive disambiguation picker when multiple fuzzy candidates
- [ ] Additional languages (French, Portuguese, German, Japanese…)

## Development

```sh
git clone https://github.com/aaangelmartin/GoTo
cd GoTo
make build       # ./bin/goto
make test        # go test ./... with -race
make lint        # golangci-lint
make snapshot    # local goreleaser dry-run
```

When adding a new user-facing string, add it to **both** language maps in [`internal/i18n/catalog.go`](./internal/i18n/catalog.go). A parity test will fail otherwise.

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

[Apache 2.0](./LICENSE) © Ángel Martín — see [NOTICE](./NOTICE).

---

<div align="center">
Built with <a href="https://charm.sh">Charm</a> · <a href="https://github.com/spf13/cobra">Cobra</a> · <a href="https://goreleaser.com">goreleaser</a>
</div>
