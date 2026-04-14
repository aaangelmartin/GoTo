# Changelog

All notable changes to this project are documented in this file. / Todos los cambios relevantes se documentan aquí.

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Versioning: [SemVer](https://semver.org/).

## [Unreleased]

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

[Unreleased]: https://github.com/aaangelmartin/GoTo/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/aaangelmartin/GoTo/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/aaangelmartin/GoTo/releases/tag/v0.1.0
