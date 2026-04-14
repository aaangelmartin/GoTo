# Changelog

All notable changes to this project are documented in this file. / Todos los cambios relevantes se documentan aquí.

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Versioning: [SemVer](https://semver.org/).

## [Unreleased]

## [0.5.0] — 2026-04-14

### Added / Añadido
- **In-TUI Settings screen** — press `o` to open a full Settings view with cycling rows (language, theme, browser, default action, directory mode), text rows (terminal, search engine, per-type openers) and per-extension overrides with add/edit/delete. Replaces the `$EDITOR` shell-out that never worked inside Bubble Tea's altscreen.
- **Pantalla de Ajustes en la TUI** — pulsa `o` para abrir los ajustes dentro del TUI.
- **Explicit Type field** in the Add/Edit form — 5th row with a colored badge, cycled with `↑/↓`, covering auto · url · mailto · ssh · file · directory · command. Placeholders now advertise every supported shape.
- **Campo Type explícito** en el formulario — deja claro que se puede guardar URLs, rutas, mails, ssh y comandos, no solo URLs.
- **`goto opener list|set|unset`** and **`goto completion install`** for zero-friction setup of per-type apps and shell completion.

### Fixed / Arreglado
- **Themes look different now**: each theme keeps its signature accent (Dracula pink, Catppuccin mauve, Nord frost, Tokyo Night purple); `#00B5E2` stays as the "default" theme. The wizard previews the theme live as you scroll options.
- **Los temas se distinguen** — cada uno con su acento propio; preview en vivo.
- **`default_action` no longer restricts** — it used to override the type detector unconditionally, so picking "url" made paths and emails open as URLs. It now only tie-breaks bare ambiguous tokens.
- **`default_action` ya no restringe** — solo desempata palabras ambiguas.

### Changed / Cambiado
- `make build` ad-hoc-signs the binary on macOS to avoid Gatekeeper SIGKILL on local Apple Silicon builds.
- `make build` firma ad-hoc el binario en macOS para evitar SIGKILL de Gatekeeper.

## [0.4.0] — 2026-04-14

### Added / Añadido
- **First-run onboarding** — launching `goto` for the first time with no config/aliases opens a bilingual TUI wizard (welcome → language → theme → browser → default action → first alias → done) that saves `config.toml` so it never runs twice.
- **Asistente inicial** — al lanzar `goto` por primera vez sin config/aliases se abre un asistente TUI bilingüe (bienvenida → idioma → tema → navegador → acción por defecto → primer alias → fin) que guarda `config.toml` y no se vuelve a mostrar.

### Changed / Cambiado
- Shields.io badges repainted in the `for-the-badge` style with `labelColor=000000` and purpose-specific accents (`#00B5E2` project, `#D22128` Apache license, `#22D3A6` CI, `#F9C900` Homebrew, `#FF79C6` Charm, `#00ADD8` Go).
- Badges renovados con estilo `for-the-badge` + `labelColor=000000` y acentos por propósito.
- `CLAUDE.md` and the `.claude/` workspace directory are now git-ignored.

## [0.3.0] — 2026-04-14

### Added / Añadido
- **Multi-type aliases** — new `Type` field on each alias: `url`, `mailto`, `ssh`, `file`, `directory`, `command`, or `auto` (default, detected at open time). Backwards-compatible migration from the legacy `{"url": "..."}` shape.
- **Alias multi-tipo** — nuevo campo `Type`: `url`, `mailto`, `ssh`, `file`, `directory`, `command` o `auto` (por defecto). Migración retrocompatible.
- Auto-detector inspects the string (scheme, email regex, filesystem stat, user@host) and resolves `auto` at open time.
- New `opener` package dispatches per type: URLs via browser; mailto via mail client; SSH via terminal; files via OS default (or per-extension override); directories via shell wrapper (`cd`).
- `goto shell-init [bash|zsh|fish]` emits a shell wrapper that eval's a `cd` directive when a directory alias is opened.
- `goto add --type <kind>` lets you fix the type explicitly when auto-detection picks the wrong lane.
- TUI: type badges colored by kind; `L` key cycles the interface language and **persists the choice** to `config.toml` (`language` field).
- Tab completion for alias names on the root command (works once `goto completion zsh|bash|fish` is installed).

### Changed / Cambiado
- Color palette recolored around `#00B5E2` across every theme.
- Paleta recoloreada con `#00B5E2` como acento principal en todos los temas.
- Config gains `[openers]`, `default_action`, `directory_mode`, `terminal`, `language` fields.
- README banner regenerated with figlet's `standard` font; shields use the `#00B5E2` accent.

## [0.2.0] — 2026-04-14

### Added / Añadido
- Bilingual interface (English + Spanish) with auto-detection via `$LANG`/`$LC_ALL`/`GOTO_LANG` and a global `--lang {en|es}` flag.
- Interfaz bilingüe (inglés + español) con autodetección vía `$LANG`/`$LC_ALL`/`GOTO_LANG` y flag global `--lang {en|es}`.
- `internal/i18n` package with catalog parity test.
- Bilingual `README.md` + `README.es.md`; bilingual `CONTRIBUTING.md`.
- `NOTICE` file with Apache 2.0 attribution.

### Changed / Cambiado
- License switched from **MIT** to **Apache 2.0** for explicit patent grant.
- Licencia cambiada de **MIT** a **Apache 2.0** para incluir cesión explícita de patentes.
- All CLI and TUI user-facing strings go through the i18n catalog.

## [0.1.0] — 2026-04-14

### Added
- First public release.
- `goto <target>` opens a URL (auto-prepends `https://`) or resolves an alias.
- Subcommands: `add`, `rm`, `ls`, `edit`, `search`, `import`, `export`,
  `config`, `completion`, `version`.
- Interactive Bubble Tea TUI with live filter, preview pane, add/edit/delete
  forms, confirm modal, help overlay, clipboard yank and tag filtering.
- Five themes: `default`, `dracula`, `catppuccin`, `nord`, `tokyonight`.
- TOML configuration (XDG-compliant) with `GOTO_CONFIG` / `GOTO_ALIASES`
  overrides.
- Cross-platform opener (macOS, Linux, Windows) with per-platform browser
  aliases.
- Fuzzy alias resolver (exact > prefix > substring > subsequence-density).
- Homebrew tap and multi-arch GitHub Releases via goreleaser.

[Unreleased]: https://github.com/aaangelmartin/GoTo/compare/v0.5.0...HEAD
[0.5.0]: https://github.com/aaangelmartin/GoTo/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/aaangelmartin/GoTo/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/aaangelmartin/GoTo/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/aaangelmartin/GoTo/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/aaangelmartin/GoTo/releases/tag/v0.1.0
