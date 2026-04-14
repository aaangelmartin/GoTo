# Contributing · Contribuir

[English](#english) · [Español](#español)

---

<a id="english"></a>

## English

Thanks for wanting to contribute to **goto**! This project is small, opinionated, and friendly. The guidelines below keep it tidy.

### Getting set up

```sh
git clone https://github.com/aaangelmartin/GoTo
cd GoTo
make build    # produces ./bin/goto
make test     # runs tests with -race
make lint     # runs golangci-lint
```

Required:
- Go ≥ 1.24
- `golangci-lint` (`brew install golangci-lint`)
- `vhs` (optional, only to regenerate the demo GIF: `brew install vhs`)

### Branch & commit workflow

- Work on a dedicated branch (`feat/<slug>`, `fix/<slug>`, `chore/<slug>`).
- **Never push directly to `main`** — open a pull request.
- Use [Conventional Commits](https://www.conventionalcommits.org) (`feat:`, `fix:`, `chore:`, `docs:`, `ci:`, `build:`, `refactor:`, `test:`).
- Commit messages are written in **English** (standard for OSS); code comments and UI strings live in both languages via `internal/i18n/catalog.go`.
- Group related work into one commit; avoid noise like `wip`.

### Adding a user-facing string

All UI strings go through the i18n catalog:

1. Add the key + English value to `catalog[EN]` in [`internal/i18n/catalog.go`](./internal/i18n/catalog.go).
2. Add the Spanish value to `catalog[ES]` (a parity test enforces both).
3. Use `i18n.T("key")` or `i18n.Tf("key", args...)` at the call site.

### Pull requests

- Describe _why_ (link the issue), not just what.
- Include test coverage for new logic in `internal/**`.
- Run `make lint test` locally and make sure both are green.

### Releases

Releases are cut by tagging `vX.Y.Z` on `main`:

```sh
git tag v0.2.0
git push origin v0.2.0
```

GitHub Actions builds binaries for macOS (arm64/amd64), Linux (arm64/amd64) and Windows (amd64) with `goreleaser`, publishes them to the release, and updates the Homebrew tap automatically.

---

<a id="español"></a>

## Español

¡Gracias por querer contribuir a **goto**! El proyecto es pequeño, opinado y amable. Estas normas lo mantienen ordenado.

### Puesta en marcha

```sh
git clone https://github.com/aaangelmartin/GoTo
cd GoTo
make build    # genera ./bin/goto
make test     # tests con -race
make lint     # golangci-lint
```

Necesario:
- Go ≥ 1.24
- `golangci-lint` (`brew install golangci-lint`)
- `vhs` (opcional, solo para regenerar el GIF de demo: `brew install vhs`)

### Ramas y commits

- Trabaja en una rama dedicada (`feat/<slug>`, `fix/<slug>`, `chore/<slug>`).
- **Nunca pushees directamente a `main`** — abre un pull request.
- Usa [Conventional Commits](https://www.conventionalcommits.org) (`feat:`, `fix:`, `chore:`, `docs:`, `ci:`, `build:`, `refactor:`, `test:`).
- Los mensajes de commit van en **inglés** (estándar OSS); los comentarios del código y las cadenas de UI están en ambos idiomas vía `internal/i18n/catalog.go`.
- Agrupa trabajo relacionado en un commit; evita ruido tipo `wip`.

### Añadir una cadena nueva visible para el usuario

Todas las cadenas de UI pasan por el catálogo i18n:

1. Añade la clave + valor en inglés en `catalog[EN]` de [`internal/i18n/catalog.go`](./internal/i18n/catalog.go).
2. Añade el valor en español en `catalog[ES]` (un test de paridad verifica ambos).
3. Usa `i18n.T("clave")` o `i18n.Tf("clave", args...)` en el punto de uso.

### Pull requests

- Describe el _por qué_ (enlaza la issue), no solo el qué.
- Incluye tests para nueva lógica en `internal/**`.
- Ejecuta `make lint test` en local y asegúrate de que pasan ambos.

### Releases

Las releases se cortan etiquetando `vX.Y.Z` en `main`:

```sh
git tag v0.2.0
git push origin v0.2.0
```

GitHub Actions compila binarios para macOS (arm64/amd64), Linux (arm64/amd64) y Windows (amd64) con `goreleaser`, los publica en la release y actualiza el tap de Homebrew automáticamente.
