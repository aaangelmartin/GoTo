<div align="center">

```
  ____     _____
 / ___| __|_   _|__
| |  _ / _ \| |/ _ \
| |_| | (_) | | (_) |
 \____|\___/|_|\___/
```

### Abre cualquier cosa desde tu terminal — URLs, archivos, directorios, mail, SSH.

[![CI](https://img.shields.io/github/actions/workflow/status/aaangelmartin/GoTo/ci.yml?branch=main&label=CI&style=for-the-badge&labelColor=000000&color=22D3A6)](https://github.com/aaangelmartin/GoTo/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/aaangelmartin/GoTo?sort=semver&style=for-the-badge&labelColor=000000&color=00B5E2)](https://github.com/aaangelmartin/GoTo/releases/latest)
[![Licencia: Apache 2.0](https://img.shields.io/badge/Licencia-Apache%202.0-D22128?style=for-the-badge&labelColor=000000)](./LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=for-the-badge&labelColor=000000&logo=go&logoColor=white)](https://go.dev)
[![Made with Charm](https://img.shields.io/badge/made%20with-Charm-FF79C6?style=for-the-badge&labelColor=000000)](https://charm.sh)
[![Idiomas](https://img.shields.io/badge/lang-EN%20%C2%B7%20ES-00B5E2?style=for-the-badge&labelColor=000000)](./README.md)
[![Homebrew](https://img.shields.io/badge/brew-aaangelmartin%2Ftap%2Fgoto-F9C900?style=for-the-badge&labelColor=000000)](https://github.com/aaangelmartin/homebrew-tap)

[English](./README.md) · **Español**

</div>

---

## ¿Por qué `goto`?

`open google.com` falla en macOS porque `open(1)` interpreta `google.com` como archivo. Acabas escribiendo `open https://google.com` todo el rato. `goto` lo arregla — y va más allá: **un solo comando para cualquier destino**.

```sh
$ goto google.com                    # URL         → https://google.com
$ goto yo@ejemplo.com                # mail        → abre cliente de correo
$ goto user@miserver --type ssh      # ssh         → abre terminal + ssh
$ goto ~/Descargas                   # directorio  → cd (requiere wrapper)
$ goto ~/notas.md                    # archivo     → app por defecto
$ goto gh                            # alias       → resuelto por fuzzy
$ goto                               # sin args    → TUI interactiva
```

### Qué lo hace distinto

- **Asistente en el primer arranque** — la primera vez que lanzas `goto` sin config, entras en una TUI de 30 segundos que elige idioma, tema, navegador y acción por defecto, y te ofrece crear tu primer alias.
- **URLs sin fricción** — añade `https://` si falta; respeta todos los esquemas conocidos (`http`, `mailto`, `ssh`, `file`, `sftp`, `chrome`, …).
- **Todo es un alias** — URLs, emails, hosts SSH, archivos, directorios, incluso comandos de shell. Cada alias recuerda su tipo: `goto dev` hace `cd`, `goto pdfs` abre Preview, `goto gh` abre el navegador.
- **Autodetección inteligente** — `goto` inspecciona el destino y elige el opener adecuado: existe en disco → archivo/dir; forma de email → mailto; `ssh://` o `user@host` → SSH; dominio simple → URL. Se puede anular con `--type` o con `default_action` en config.
- **Configurable por tipo** — cada tipo elige su app (`[openers]` en `config.toml`). También por extensión (`.pdf = "Preview"`, `.md = "cursor"`).
- **Alias de directorio que hacen `cd`** — sourcea `goto shell-init zsh`/`bash`/`fish` una vez y `goto dev` cambia el directorio de tu shell (como el `goto` clásico de iridakos).
- **Coincidencia difusa** — `goto githu` sigue resolviendo `github`.
- **Autocompletado con tab** — `goto <tab>` sugiere tus alias. `goto completion zsh` lo instala.
- **TUI preciosa** — con [Bubble Tea](https://github.com/charmbracelet/bubbletea). Badges coloreados por tipo. Filtro en vivo, panel de preview, añadir/editar/borrar, portapapeles, filtro por etiqueta, 5 temas y **cambio de idioma permanente** (tecla `L`).
- **Bilingüe** — inglés y español autodetectados por `$LANG`, anulables con `--lang` (puntual) o `L` en la TUI (persistente).
- **Multiplataforma** — macOS, Linux, Windows. Binario estático único.
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

### Gestionar alias (de cualquier tipo)

```sh
# URL (detectado automáticamente)
goto add gh github.com/aaangelmartin --tag dev --desc "Mi GitHub"

# Email → abre el cliente de correo
goto add jefe jefe@miempresa.com --type mailto

# SSH → abre un terminal nuevo con ssh
goto add prod user@prod.miempresa.com --type ssh

# Directorio → `goto dev` hará `cd` (requiere shell-init una vez)
goto add dev ~/Desarrollo

# Archivo → abre con la app por defecto (o override por extensión)
goto add notas ~/Documentos/ideas.md

# Comando de shell
goto add deploy "make deploy" --type command

goto ls                               # tabla con badges por tipo
goto ls --tag curro --json            # salida para scripts
goto edit gh --url github.com/aaangelmartin/GoTo
goto rm jira -y                       # sin confirmación
goto export > alias.json              # copia de seguridad
goto import alias.json --overwrite    # restaurar
```

### Integración con shell (para alias de directorio)

Para que `goto <alias-de-dir>` haga `cd` de verdad en tu shell (un proceso hijo no puede `cd` en el padre), sourcea el wrapper una vez:

```sh
# zsh
echo 'eval "$(goto shell-init zsh)"' >> ~/.zshrc

# bash
echo 'eval "$(goto shell-init bash)"' >> ~/.bashrc

# fish
echo 'goto shell-init fish | source' >> ~/.config/fish/config.fish
```

URLs, mail, ssh, archivos y comandos funcionan sin el wrapper — solo importa para directorios.

### Lanzar la TUI

```sh
goto
```

<div align="center">

| tecla | acción |
| --- | --- |
| `↑`/`k`, `↓`/`j` | mover selección |
| `enter` | abrir |
| `/` | filtro en vivo |
| `a` / `e` / `d` | añadir / editar / borrar |
| `t` | filtro por etiqueta |
| `y` | copiar URL al portapapeles |
| `L` | cambiar idioma (permanente) |
| `?` | ayuda de teclado |
| `q` | salir |

</div>

### Idioma

La interfaz sigue tu variable `$LANG` automáticamente (cualquier locale `es*` → español, en caso contrario inglés).
Puedes forzar un idioma con:

```sh
goto --lang es ls              # una sola vez
export GOTO_LANG=es            # por sesión
# o pulsa `L` en la TUI        # persistente (se guarda en config.toml)
```

### Autocompletado de shell

```sh
goto completion zsh  > "${fpath[1]}/_goto"
goto completion bash > /etc/bash_completion.d/goto
goto completion fish > ~/.config/fish/completions/goto.fish
```

Tras instalarlo, `goto <tab>` autocompletará los nombres de tus alias.

## Configuración

La config vive en `$XDG_CONFIG_HOME/goto/config.toml` (por defecto `~/.config/goto/config.toml`). Los alias viven en `$XDG_DATA_HOME/goto/aliases.json`.

```toml
browser         = "default"          # default | chrome | firefox | safari | arc | brave | edge | opera | vivaldi
search_engine   = "https://www.google.com/search?q={q}"
theme           = "default"          # default | dracula | catppuccin | nord | tokyonight
language        = "auto"             # auto | en | es   (se guarda desde la tecla L de la TUI)
confirm_delete  = true
fuzzy_threshold = 0.4

default_action  = "auto"             # auto | url | file | directory   (qué hacer con un arg sin alias)
directory_mode  = "shell"            # shell (cd via wrapper) | finder
terminal        = ""                 # terminal para SSH (vacío = por defecto del SO)

[openers]
url       = "default"                # navegador para URLs
mailto    = "default"                # cliente de correo
ssh       = ""                       # comando ssh (normalmente vacío)
file      = "default"                # app por defecto para archivos
directory = ""                       # solo si directory_mode = "finder"
".pdf"    = "Preview"                # override por extensión
".md"     = "cursor"
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
| `--browser <nombre>` | anula el navegador en esta llamada |
| `--no-https` | antepone `http://` en vez de `https://` si no hay protocolo |
| `--dry-run` | imprime el destino resuelto + tipo sin abrirlo |
| `--lang <en\|es>` | fuerza el idioma de la interfaz |
| `--type <tipo>` (en `add`) | tipo explícito: `url`, `mailto`, `ssh`, `file`, `directory`, `command`, `auto` |

## Cómo funciona la resolución

Cuando ejecutas `goto foo`, `goto` prueba en este orden:

1. **¿Alias exacto?** — resuelve el tipo guardado y abre. Incrementa contador.
2. **Coincidencia difusa** — si hay un ganador claro por encima de `fuzzy_threshold`, lo abre. Si hay ambigüedad, error con los candidatos.
3. **Autodetecta el tipo del argumento**:
   - empieza por un esquema conocido (`http://`, `mailto:`, `ssh://`, `file://`, …) → ese tipo
   - encaja con regex de email → `mailto`
   - pinta de ruta (`/…`, `~/…`, `./…`, `C:\…`) → stat al disco: directorio vs archivo
   - `user@host` sin TLD → `ssh`
   - tiene punto y no tiene espacios → `url` (antepone `https://`)
   - tiene espacios → `command` (se ejecuta en shell)
4. **`default_action`** en la config puede forzar la detección (`url`, `file`, `directory`, `auto`).

## Roadmap

- [x] Alias multi-tipo (URL, mailto, ssh, archivo, directorio, comando)
- [x] Wrapper de shell para alias de directorio (`goto shell-init`)
- [x] Autocompletado con tab para alias
- [x] Cambio de idioma permanente desde la TUI
- [ ] Autocompletado de rutas/URLs dentro de los formularios de la TUI
- [ ] Backend de sincronización (gist / repo git) para alias multi-dispositivo
- [ ] Filtrado por grupos / workspaces
- [ ] Vista de historial en la TUI
- [ ] `goto recent` para MRU
- [ ] Selección de perfil de navegador (Chrome `--profile`)
- [ ] Scoop bucket para Windows
- [ ] Selector interactivo de desambiguación
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

Ver [CONTRIBUTING.md](./CONTRIBUTING.md).

## Licencia

[Apache 2.0](./LICENSE) © Ángel Martín — ver [NOTICE](./NOTICE).

---

<div align="center">
Hecho con <a href="https://charm.sh">Charm</a> · <a href="https://github.com/spf13/cobra">Cobra</a> · <a href="https://goreleaser.com">goreleaser</a>
</div>
