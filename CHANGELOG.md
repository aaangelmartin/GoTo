# Changelog

All notable changes to this project are documented in this file. The format is
based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this
project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.1.0] — 2026-04-14

### Added
- First public release.
- `goto <target>` opens a URL (auto-prepends `https://`) or resolves an alias.
- Subcommands: `add`, `rm`, `ls`, `edit`, `search`, `import`, `export`,
  `config`, `completion`, `version`.
- Interactive Bubble Tea TUI with live filter, preview pane, add/edit/delete
  forms, confirm modal, help overlay, clipboard yank and tag filtering.
- Four themes: `default`, `dracula`, `catppuccin`, `nord`, `tokyonight`.
- TOML configuration (XDG-compliant) with `GOTO_CONFIG` / `GOTO_ALIASES`
  overrides.
- Cross-platform opener (macOS, Linux, Windows) with per-platform browser
  aliases.
- Fuzzy alias resolver (exact > prefix > substring > subsequence-density).
- Homebrew tap and multi-arch GitHub Releases via goreleaser.

[Unreleased]: https://github.com/aaangelmartin/GoTo/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/aaangelmartin/GoTo/releases/tag/v0.1.0
