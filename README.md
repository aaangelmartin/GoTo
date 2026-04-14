<div align="center">

```
  ____     _____
 / ___| __|_   _|__
| |  _ / _ \| |/ _ \
| |_| | (_) | | (_) |
 \____|\___/|_|\___/
```

### Open anything from your terminal — URLs, files, directories, mail, SSH hosts.

[![CI](https://img.shields.io/github/actions/workflow/status/aaangelmartin/GoTo/ci.yml?branch=main&label=CI&style=for-the-badge&labelColor=000000&color=22D3A6)](https://github.com/aaangelmartin/GoTo/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/aaangelmartin/GoTo?sort=semver&style=for-the-badge&labelColor=000000&color=00B5E2)](https://github.com/aaangelmartin/GoTo/releases/latest)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-D22128?style=for-the-badge&labelColor=000000)](./LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=for-the-badge&labelColor=000000&logo=go&logoColor=white)](https://go.dev)
[![Made with Charm](https://img.shields.io/badge/made%20with-Charm-FF79C6?style=for-the-badge&labelColor=000000)](https://charm.sh)
[![Languages](https://img.shields.io/badge/lang-EN%20%C2%B7%20ES-00B5E2?style=for-the-badge&labelColor=000000)](./README.es.md)
[![Homebrew](https://img.shields.io/badge/brew-aaangelmartin%2Ftap%2Fgoto-F9C900?style=for-the-badge&labelColor=000000)](https://github.com/aaangelmartin/homebrew-tap)

**English** · [Español](./README.es.md)

</div>

---

## Why `goto`?

`open google.com` fails on macOS because `open(1)` thinks `google.com` is a file. You end up typing `open https://google.com` every time. `goto` fixes that — and keeps going: one command, any target.

```sh
$ goto google.com                    # URL          → https://google.com
$ goto me@example.com                # mail         → opens mail client
$ goto user@myserver --type ssh      # ssh          → opens terminal + ssh
$ goto ~/Downloads                   # directory    → cd (via shell wrapper)
$ goto ~/notes.md                    # file         → default app (or per-ext)
$ goto gh                            # saved alias  → fuzzy resolved
$ goto                               # no args      → interactive TUI
```

### What makes it different

- **First-run wizard** — the first time you launch `goto` with no config it drops you into a 30-second TUI that picks your language, theme, default browser and default action, then offers to create your first alias.
- **Zero-friction URLs** — auto-prepends `https://` when missing; respects every known scheme (`http`, `mailto`, `ssh`, `file`, `sftp`, `chrome`, …).
- **Everything is an alias** — URLs, emails, SSH hosts, files, directories, even raw shell commands. Each alias remembers its type, so `goto dev` can `cd`, `goto pdfs` can open Preview, and `goto gh` can open your browser.
- **Smart auto-detection** — `goto` inspects the target and picks the right opener: path exists on disk → file/dir; email-shaped → mailto; `ssh://` or `user@host` → SSH; bare domain → URL. Override with `--type` or `default_action` in config.
- **Fully configurable per type** — each type picks its own app (`[openers]` table in `config.toml`). Per-extension overrides too (`.pdf = "Preview"`, `.md = "cursor"`).
- **Directory aliases that actually `cd`** — source `goto shell-init zsh`/`bash`/`fish` once and `goto dev` changes your shell's working directory (just like the classic iridakos `goto`).
- **Fuzzy matching** — `goto githu` still hits `github`.
- **Tab completion** — `goto <tab>` suggests your aliases. `goto completion zsh` installs it.
- **Beautiful TUI** — [Bubble Tea](https://github.com/charmbracelet/bubbletea) powered. Type badges color each alias by kind. Live filter, preview pane, add/edit/delete, clipboard, tag filter, five themes, and a **permanent language toggle** (`L` key).
- **Bilingual** — English and Spanish auto-detected from `$LANG`, overridable with `--lang` (one-off) or `L` in the TUI (persisted).
- **Cross-platform** — macOS, Linux, Windows. Single static binary.
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

### Manage aliases (any type)

```sh
# URL (detected automatically)
goto add gh github.com/aaangelmartin --tag dev --desc "My GitHub"

# Email → opens your mail client
goto add ceo bossperson@example.com --type mailto

# SSH → opens a new terminal window with an ssh session
goto add prod user@prod.example.com --type ssh

# Directory → `goto dev` will `cd` (needs shell-init once)
goto add dev ~/Development

# File → opens with the default app (or per-extension override)
goto add notes ~/Documents/brain-dump.md

# Raw shell command
goto add deploy "make deploy" --type command

goto ls                               # pretty table with type badges
goto ls --tag work --json             # machine-readable
goto edit gh --url github.com/aaangelmartin/GoTo
goto rm jira -y                       # skip confirmation
goto export > aliases.json            # backup
goto import aliases.json --overwrite  # restore
```

### Shell integration (for directory aliases)

To make `goto <dir-alias>` actually `cd` in your current shell (because a child process can't `cd` its parent), source the wrapper once:

```sh
# zsh
echo 'eval "$(goto shell-init zsh)"' >> ~/.zshrc

# bash
echo 'eval "$(goto shell-init bash)"' >> ~/.bashrc

# fish
echo 'goto shell-init fish | source' >> ~/.config/fish/config.fish
```

URLs, mail, ssh, files, and commands still work without the wrapper — it only matters for directories.

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
search_engine   = "https://www.google.com/search?q={q}"   # DuckDuckGo, Kagi, Perplexity, etc.
theme           = "default"          # default | dracula | catppuccin | nord | tokyonight
language        = "auto"             # auto | en | es   (persisted from the TUI 'L' key)
confirm_delete  = true
fuzzy_threshold = 0.4                # lower = more lenient fuzzy matches

default_action  = "auto"             # auto | url | file | directory  (what an un-aliased arg becomes)
directory_mode  = "shell"            # shell (cd via wrapper) | finder
terminal        = ""                 # terminal emulator for SSH (empty = OS default)

[openers]
url       = "default"                # which browser for URLs
mailto    = "default"                # mail client
ssh       = ""                       # ssh command (usually empty)
file      = "default"                # default app for files
directory = ""                       # only used when directory_mode = "finder"
".pdf"    = "Preview"                # per-extension overrides
".md"     = "cursor"
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
| `--dry-run` | print the resolved target + type without opening it |
| `--lang <en\|es>` | force interface language for this invocation |
| `--type <kind>` (on `add`) | explicit type: `url`, `mailto`, `ssh`, `file`, `directory`, `command`, `auto` |

## How resolution works

When you run `goto foo`, `goto` tries these in order:

1. **Exact alias match?** — resolve the stored type and open. Bump hit counter.
2. **Fuzzy alias match** — if one clear winner above `fuzzy_threshold`, open it. Multiple ambiguous candidates → error listing them.
3. **Auto-detect type of the raw argument**:
   - starts with a known scheme (`http://`, `mailto:`, `ssh://`, `file://`, …) → that type
   - matches an email regex → `mailto`
   - path-looking (`/…`, `~/…`, `./…`, `C:\…`) → stat the filesystem: directory vs file
   - `user@host` with no TLD → `ssh`
   - contains a dot, no whitespace → `url` (prepend `https://`)
   - contains whitespace → `command` (shell-exec)
4. **`default_action`** in config can override the detection (`url`, `file`, `directory`, `auto`).

## Roadmap

- [x] Multi-type aliases (URL, mailto, ssh, file, directory, command)
- [x] Shell wrapper for directory aliases (`goto shell-init`)
- [x] Tab completion for aliases
- [x] Permanent language toggle from the TUI
- [ ] TUI autocomplete for paths/URLs inside add/edit forms
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
