# Contributing

Thanks for wanting to contribute to **goto**! This project is small, opinionated, and friendly. The guidelines below keep it tidy.

## Getting set up

```sh
git clone https://github.com/aaangelmartin/GoTo
cd GoTo
make build    # produces ./bin/goto
make test     # runs tests with -race
make lint     # runs golangci-lint
```

Required:
- Go ≥ 1.22
- `golangci-lint` (via `brew install golangci-lint`)
- `vhs` (optional, only to regenerate the demo GIF: `brew install vhs`)

## Branch & commit workflow

- Work on a dedicated branch (`feat/<slug>`, `fix/<slug>`, `chore/<slug>`).
- **Never push directly to `main`** — open a pull request.
- Use [Conventional Commits](https://www.conventionalcommits.org) (`feat:`, `fix:`, `chore:`, `docs:`, `ci:`, `build:`, `refactor:`, `test:`).
- Group related work into one commit; avoid noise commits like `wip`. Rebase/squash locally before opening the PR if needed.

## Pull requests

- Describe _why_ (link the issue if applicable), not just what.
- Include test coverage for new logic in `internal/**`.
- Run `make lint test` locally and make sure both are green.
- Keep PRs focused. If a refactor and a feature are entangled, split them.

## Issues

Before filing a bug, try:

```sh
goto version
goto --dry-run <target>
```

…and include that output. For feature requests, describe the user-facing pain and a concrete example before proposing an implementation.

## Releases

Releases are cut by tagging `vX.Y.Z` on `main`:

```sh
git tag v0.2.0
git push origin v0.2.0
```

GitHub Actions builds binaries for macOS (arm64/amd64), Linux (arm64/amd64) and Windows (amd64) with `goreleaser`, publishes them to the release, and updates the Homebrew tap automatically.
