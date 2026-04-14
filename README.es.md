<div align="center">

```
   ___  ___  ____  ____
  / __)/ _ \(_  _)/ _  )
 ( (_-\(_)  ) )(  )(_) )
  \___/\___/ (__) \____)
```

### Abre cualquier URL y gestiona tus alias de enlaces desde la terminal.

[![CI](https://github.com/aaangelmartin/GoTo/actions/workflows/ci.yml/badge.svg)](https://github.com/aaangelmartin/GoTo/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/aaangelmartin/GoTo?sort=semver&color=%23FF79C6)](https://github.com/aaangelmartin/GoTo/releases/latest)
[![License](https://img.shields.io/badge/licencia-Apache%202.0-8BE9FD)](./LICENSE)
[![Go Report](https://goreportcard.com/badge/github.com/aaangelmartin/goto)](https://goreportcard.com/report/github.com/aaangelmartin/goto)
[![Made with Charm](https://img.shields.io/badge/made%20with-%F0%9F%92%9C%20Charm-FF79C6)](https://charm.sh)

[English](./README.md) · **Español**

</div>

---

## ¿Por qué `goto`?

`open google.com` falla en macOS porque `open(1)` interpreta `google.com` como un archivo. Acabas escribiendo `open https://google.com` todo el rato. Se acabó.

```sh
$ goto google.com           # abre https://google.com
$ goto gh                   # abre tu alias guardado 'gh'
$ goto                      # lanza la TUI interactiva
```

- **URLs sin fricción** — añade `https://` si falta; respeta `http://`, `mailto:`, `ssh://`, `file://` y cualquier otro protocolo.
- **Alias con nombre** — `gh`, `mail`, `jira`, `localhost3000`. Memoriza nombres, deja de teclear URLs.
- **Coincidencia difusa** — `goto githu` sigue resolviendo `github`.
- **TUI preciosa** — con [Bubble Tea](https://github.com/charmbracelet/bubbletea). Añade, edita, borra, busca, copia, filtra por etiqueta. Cuatro temas integrados.
- **Bilingüe** — interfaz en español e inglés, detectada automáticamente por `LANG`; se puede forzar con `--lang`.
- **Multiplataforma** — macOS, Linux, Windows. Binario estático único, ~5 MB.
- **Privado por defecto** — todo vive en `~/.local/share/goto/` y `~/.config/goto/`. Sin telemetría.

## Instalación

### Homebrew (macOS / Linux)

```sh
brew install aaangelmartin/tap/goto
```

### Go

```sh
go install github.com/aaangelmartin/goto/cmd/goto@latest
```

### Manual

Descarga el binario para tu plataforma desde la [última release](https://github.com/aaangelmartin/GoTo/releases/latest), descomprime y pon `goto` en tu `$PATH`.

## Uso

### Abrir cosas

```sh
goto google.com                       # https://google.com
goto localhost:3000 --no-https        # http://localhost:3000
goto https://claude.ai                # pass-through
goto search "claude code"             # abre tu buscador configurado
```

### Gestionar alias

```sh
goto add gh github.com/aaangelmartin --tag dev --desc "Mi GitHub"
goto add jira mi-empresa.atlassian.net/jira --tag curro
goto ls                               # tabla bonita
goto ls --tag curro --json            # salida para scripts
goto edit gh --url github.com/aaangelmartin/GoTo
goto rm jira -y                       # sin confirmación
goto export > alias.json              # copia de seguridad
goto import alias.json --overwrite    # restaurar
```

### Lanzar la TUI

```sh
goto
```

<div align="center">

| tecla | acción |
| --- | --- |
| `↑`/`k`, `↓`/`j` | mover selección |
| `enter` | abrir en navegador |
| `/` | filtro en vivo |
| `a` / `e` / `d` | añadir / editar / borrar |
| `t` | activar/desactivar filtro por etiqueta |
| `y` | copiar URL al portapapeles |
| `?` | ayuda de teclado |
| `q` | salir |

</div>

### Idioma

La interfaz sigue automáticamente tu variable `$LANG` (cualquier locale `es*` → español, en caso contrario inglés).
Puedes forzar un idioma con:

```sh
goto --lang es ls        # una sola vez
export GOTO_LANG=es      # para toda la sesión
```

### Autocompletado de shell

```sh
goto completion zsh  > "${fpath[1]}/_goto"   # zsh
goto completion bash > /etc/bash_completion.d/goto
goto completion fish > ~/.config/fish/completions/goto.fish
```

## Configuración

La config vive en `$XDG_CONFIG_HOME/goto/config.toml` (por defecto `~/.config/goto/config.toml`). Los alias viven en `$XDG_DATA_HOME/goto/aliases.json`.

```toml
browser         = "default"          # default | chrome | firefox | safari | arc | brave | edge | opera | vivaldi
search_engine   = "https://www.google.com/search?q={q}"   # usa DuckDuckGo, Kagi, etc.
theme           = "default"          # default | dracula | catppuccin | nord | tokyonight
confirm_delete  = true
fuzzy_threshold = 0.4                # más bajo = coincidencias difusas más tolerantes
```

Sobrescribe las rutas con variables de entorno:

```sh
export GOTO_CONFIG=$HOME/.dotfiles/goto.toml
export GOTO_ALIASES=$HOME/.dotfiles/alias.json
export GOTO_LANG=es                  # o en
```

## Flags

| flag | descripción |
| --- | --- |
| `--browser <nombre>` | anula el navegador para esta llamada |
| `--no-https` | antepone `http://` en vez de `https://` si no hay protocolo |
| `--dry-run` | imprime la URL resuelta sin abrirla |
| `--lang <en\|es>` | fuerza el idioma de la interfaz en esta invocación |

## Cómo funciona la resolución

Cuando ejecutas `goto foo`, `goto` prueba en este orden:

1. **¿Protocolo presente?** (`foo://bar`, `mailto:a@b`) — abrir tal cual.
2. **¿Alias exacto?** — abrir e incrementar contador.
3. **Coincidencia difusa** — si hay un único candidato (o un ganador claro por encima de `fuzzy_threshold`), lo abre.
4. **¿Parece una URL?** (tiene punto, no tiene espacios) — normaliza (añade `https://`) y abre.
5. **Nada coincide** — sugiere `goto search "foo"`.

## Roadmap

- [ ] Backend de sincronización (gist / repo git) para alias multi-dispositivo
- [ ] Filtrado por grupos / workspaces
- [ ] Vista de historial en la TUI
- [ ] `goto recent` para MRU
- [ ] Selección de perfil de navegador (Chrome `--profile`)
- [ ] Scoop bucket para Windows
- [ ] Selector interactivo de desambiguación para múltiples candidatos
- [ ] Más idiomas (francés, portugués, alemán, japonés…)

## Desarrollo

```sh
git clone https://github.com/aaangelmartin/GoTo
cd GoTo
make build       # ./bin/goto
make test        # go test ./... con -race
make lint        # golangci-lint
make snapshot    # prueba local de goreleaser
```

Al añadir una nueva cadena visible para el usuario, añádela a **ambos** mapas de idiomas en [`internal/i18n/catalog.go`](./internal/i18n/catalog.go). Hay un test de paridad que lo verifica.

Mira [CONTRIBUTING.md](./CONTRIBUTING.md).

## Licencia

[Apache 2.0](./LICENSE) © Ángel Martín — ver [NOTICE](./NOTICE).

---

<div align="center">
Hecho con <a href="https://charm.sh">Charm</a> · <a href="https://github.com/spf13/cobra">Cobra</a> · <a href="https://goreleaser.com">goreleaser</a>
</div>
