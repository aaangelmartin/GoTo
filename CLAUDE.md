# Project rules for Claude

## Workflow
- Use Conventional Commits (`feat:`, `fix:`, `chore:`, `docs:`, `ci:`, `build:`, `refactor:`, `test:`).
- Group commits per large function/slice; don't commit after every micro-change.
- **Never push to `main`**. Always work on a feature branch and open a PR.
- Commit messages are in **English** (OSS standard).

## Bilingual rule
- Every user-facing string (CLI help, errors, TUI labels, README, CHANGELOG, CONTRIBUTING) must exist in **both English and Spanish**.
- CLI / TUI strings live in `internal/i18n/catalog.go` — add keys to both `EN` and `ES` maps. A parity test (`TestCatalogParity`) will fail if either is missing.
- `README.md` is English, `README.es.md` is Spanish, each links to the other at the top.
- `CONTRIBUTING.md` and `CHANGELOG.md` include both languages in-file.

## License
- Apache 2.0. Keep `NOTICE` updated if new contributors appear.
