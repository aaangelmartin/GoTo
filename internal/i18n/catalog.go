package i18n

// catalog holds every user-facing string in both supported languages.
// Catálogo con todas las cadenas visibles para el usuario en ambos idiomas.
//
// Keys are lowercase snake_case and should stay stable — they are referenced
// from CLI and TUI code. When adding a new string, add it to BOTH maps.
var catalog = map[Lang]map[string]string{
	EN: {
		// CLI — root
		"short":        "Open URLs and manage link aliases from the terminal",
		"long":         "goto opens URLs in your browser (auto-prepends https:// if missing)\nand lets you manage personal link aliases with a beautiful TUI.\n\n  goto google.com         Opens https://google.com\n  goto gh                 Resolves alias \"gh\" and opens it\n  goto                    Launches the interactive TUI\n  goto add gh github.com  Adds an alias",
		"flag_browser": "browser to use (default|chrome|firefox|safari|arc|brave|edge)",
		"flag_nohttps": "use http:// instead of https:// when prepending protocol",
		"flag_dryrun":  "print resolved URL without opening it",
		"flag_lang":    "UI language (en|es); overrides GOTO_LANG",

		// CLI — add
		"add_short":        "Add a new alias",
		"add_tag":          "comma-separated tags",
		"add_desc":         "description",
		"add_type":         "explicit type (url|mailto|ssh|file|directory|command); default auto-detects",
		"added":            "added: %s -> %s\n",
		"err_name_invalid": "invalid alias name %q: no whitespace or slashes",

		// CLI — rm
		"rm_short":   "Remove an alias",
		"rm_yes":     "skip confirmation",
		"rm_confirm": "Delete %q (%s)? [y/N] ",
		"rm_aborted": "aborted",
		"removed":    "removed: %s\n",

		// CLI — ls
		"ls_short":    "List all aliases",
		"ls_empty":    "no aliases yet — try: goto add <name> <url>",
		"ls_tag":      "filter by tag",
		"ls_json":     "output as JSON",
		"ls_col_name": "NAME",
		"ls_col_url":  "URL",
		"ls_col_tags": "TAGS",
		"ls_col_hits": "HITS",

		// CLI — edit
		"edit_short":  "Edit an existing alias",
		"edit_url":    "new URL",
		"edit_desc":   "new description",
		"edit_tag":    "replace tags (comma-separated)",
		"edit_name":   "rename the alias",
		"edit_exists": "alias %q already exists",
		"updated":     "updated: %s -> %s\n",

		// CLI — search
		"search_short":      "Open a web search for the given query",
		"err_search_config": "search_engine in config must contain {q}",

		// CLI — import/export
		"import_short":     "Import aliases from a JSON file",
		"import_overwrite": "overwrite existing aliases instead of skipping",
		"imported":         "imported: %d added, %d updated, %d skipped\n",
		"export_short":     "Export aliases as JSON (stdout if no file given)",

		// CLI — config
		"config_short":      "Open the config file in $EDITOR (or print its path)",
		"config_path_short": "Print the config file path",

		// CLI — completion
		"completion_short": "Generate shell completion scripts",

		// CLI — version
		"version_short": "Print version, commit and build date",
		"version_line":  "goto %s (commit %s, built %s)\n",

		// Errors (default action)
		"err_empty_target": "empty target",
		"err_ambiguous":    "ambiguous target %q; candidates: %s",
		"err_notfound":     "no alias found and %q is not a URL; try: goto search %q",
		"warn_config_load": "goto: warning loading config: %v\n",
		"err_load_aliases": "load aliases: %w",
		"err_parse":        "parse %s: %w",

		// TUI — labels
		"tui_aliases_count":    "%d aliases",
		"tui_status_opened":    "opened %s",
		"tui_status_copied":    "copied %s",
		"tui_status_copyfail":  "copy failed: %s",
		"tui_status_deleted":   "deleted %s",
		"tui_status_saved":     "saved",
		"tui_status_delfail":   "delete failed: %s",
		"tui_status_tag_set":   "filtering by #%s",
		"tui_status_tag_clear": "cleared tag filter",
		"tui_status_err":       "error: %s",
		"tui_empty":            "no aliases yet\n\npress  a  to add your first one",
		"tui_no_matches":       "no matches",
		"tui_help_list":        " enter open · a add · e edit · d delete · / filter · t tag · y copy · ? help · q quit",
		"tui_help_form":        " tab next · shift+tab prev · enter save · esc cancel",
		"tui_help_confirm":     " y/n · enter confirm · esc cancel",
		"tui_help_back":        " esc back",
		"tui_form_add":         "Add alias",
		"tui_form_edit":        "Edit alias",
		"tui_field_name":       "name",
		"tui_field_url":        "url",
		"tui_field_tags":       "tags",
		"tui_field_desc":       "description",
		"tui_confirm_delete":   "Delete alias?",
		"tui_confirm_yes":      "[ Yes ]",
		"tui_confirm_no":       "[ No ]",
		"tui_help_title":       "keyboard",
		"tui_opens":            "%d opens",
		"tui_last":             "last %s",
		"tui_created":          "created %s",
		"tui_placeholder_name": "alias name (e.g. gh)",
		"tui_placeholder_url":  "url (e.g. github.com)",
		"tui_placeholder_tags": "tags, comma separated",
		"tui_placeholder_desc": "description (optional)",
		"err_empty_name":       "name is required",
		"err_empty_url":        "url is required",
		"err_invalid_name":     "name must not contain spaces or slashes",

		// v0.3 — multi-type targets + shell init + language switch
		"shellinit_short":           "Emit a shell wrapper that lets directory aliases cd the parent shell (source into .zshrc/.bashrc/config.fish)",
		"err_shell_wrapper_missing": "goto: shell wrapper not installed — run: eval \"$(goto shell-init zsh)\"",
		"type_url":                  "url",
		"type_mailto":               "mail",
		"type_ssh":                  "ssh",
		"type_file":                 "file",
		"type_directory":            "dir",
		"type_command":              "cmd",
		"type_auto":                 "auto",
		"tui_lang_switched":         "language set to %s (permanent)",
		"help_lang":                 "cycle interface language (persisted)",

		// First-run onboarding
		"onb_welcome_title":    "Welcome to goto ✦",
		"onb_welcome_body":     "You're about to set up goto — a single command for URLs, files, directories, mail and SSH hosts.\nThis 30-second wizard picks your defaults (everything is editable later in ~/.config/goto/config.toml).",
		"onb_welcome_tip":      "press enter to begin  ·  esc to skip with defaults",
		"onb_lang_title":       "Language",
		"onb_lang_desc":        "Choose the interface language. 'auto' follows $LANG.",
		"onb_theme_title":      "Theme",
		"onb_theme_desc":       "All themes use the same accent (#00B5E2); pick the background palette you like.",
		"onb_browser_title":    "Default browser",
		"onb_browser_desc":     "Which browser opens URL aliases. 'default' follows your OS preference.",
		"onb_action_title":     "Default action",
		"onb_action_desc":      "What happens when you run `goto X` and X is not an alias: auto-detects type (recommended), or force url/file/directory.",
		"onb_firstalias_title": "Add your first alias (optional)",
		"onb_firstalias_desc":  "Try something like `gh → github.com/aaangelmartin`. You can always add more with `goto add` or in the TUI.",
		"onb_firstalias_skip":  "tab switches fields · enter saves (or skips if empty) · esc skips",
		"onb_done_title":       "You're all set 🎉",
		"onb_done_body":        "Tip: for directory aliases to cd your shell, source the wrapper once:\n    eval \"$(goto shell-init zsh)\"\n(replace zsh with bash or fish)",
		"onb_done_tip":         "press enter to start",
		"onb_footer":           " ↑/↓ choose · enter next · ← back · esc skip wizard · q quit",

		// TUI — help rows
		"help_move":   "move selection",
		"help_jump":   "jump to top / bottom",
		"help_open":   "open alias in browser",
		"help_filter": "filter list (type to narrow)",
		"help_esc":    "clear filter / go back",
		"help_add":    "add a new alias",
		"help_edit":   "edit selected alias",
		"help_delete": "delete selected alias",
		"help_yank":   "copy URL to clipboard",
		"help_tag":    "filter by selected alias' first tag (toggle)",
		"help_toggle": "toggle this help",
		"help_quit":   "quit",
		"help_config": "config:",
		"help_theme":  "theme:",
	},
	ES: {
		// CLI — root
		"short":        "Abre URLs y gestiona alias de enlaces desde la terminal",
		"long":         "goto abre URLs en tu navegador (añade https:// si falta)\ny te permite gestionar alias de enlaces con una TUI bonita.\n\n  goto google.com         Abre https://google.com\n  goto gh                 Resuelve el alias \"gh\" y lo abre\n  goto                    Lanza la TUI interactiva\n  goto add gh github.com  Añade un alias",
		"flag_browser": "navegador a usar (default|chrome|firefox|safari|arc|brave|edge)",
		"flag_nohttps": "usa http:// en lugar de https:// al anteponer protocolo",
		"flag_dryrun":  "muestra la URL resuelta sin abrirla",
		"flag_lang":    "idioma de la interfaz (en|es); anula GOTO_LANG",

		// CLI — add
		"add_short":        "Añade un nuevo alias",
		"add_tag":          "etiquetas separadas por comas",
		"add_type":         "tipo explícito (url|mailto|ssh|file|directory|command); por defecto autodetecta",
		"add_desc":         "descripción",
		"added":            "añadido: %s -> %s\n",
		"err_name_invalid": "nombre de alias inválido %q: sin espacios ni barras",

		// CLI — rm
		"rm_short":   "Elimina un alias",
		"rm_yes":     "omitir confirmación",
		"rm_confirm": "¿Borrar %q (%s)? [s/N] ",
		"rm_aborted": "cancelado",
		"removed":    "eliminado: %s\n",

		// CLI — ls
		"ls_short":    "Lista todos los alias",
		"ls_empty":    "aún no hay alias — prueba: goto add <nombre> <url>",
		"ls_tag":      "filtrar por etiqueta",
		"ls_json":     "salida en JSON",
		"ls_col_name": "NOMBRE",
		"ls_col_url":  "URL",
		"ls_col_tags": "ETIQUETAS",
		"ls_col_hits": "USOS",

		// CLI — edit
		"edit_short":  "Edita un alias existente",
		"edit_url":    "nueva URL",
		"edit_desc":   "nueva descripción",
		"edit_tag":    "reemplaza etiquetas (separadas por comas)",
		"edit_name":   "renombra el alias",
		"edit_exists": "el alias %q ya existe",
		"updated":     "actualizado: %s -> %s\n",

		// CLI — search
		"search_short":      "Abre una búsqueda web con la consulta indicada",
		"err_search_config": "search_engine en config debe contener {q}",

		// CLI — import/export
		"import_short":     "Importa alias desde un fichero JSON",
		"import_overwrite": "sobreescribe alias existentes en vez de omitirlos",
		"imported":         "importado: %d añadidos, %d actualizados, %d omitidos\n",
		"export_short":     "Exporta alias como JSON (stdout si no se da fichero)",

		// CLI — config
		"config_short":      "Abre el fichero de config en $EDITOR (o imprime su ruta)",
		"config_path_short": "Imprime la ruta del fichero de config",

		// CLI — completion
		"completion_short": "Genera scripts de autocompletado para tu shell",

		// CLI — version
		"version_short": "Muestra versión, commit y fecha de compilación",
		"version_line":  "goto %s (commit %s, compilado %s)\n",

		// Errors (default action)
		"err_empty_target": "destino vacío",
		"err_ambiguous":    "destino ambiguo %q; candidatos: %s",
		"err_notfound":     "no se encontró alias y %q no es una URL; prueba: goto search %q",
		"warn_config_load": "goto: aviso al cargar config: %v\n",
		"err_load_aliases": "cargar alias: %w",
		"err_parse":        "parsear %s: %w",

		// TUI — labels
		"tui_aliases_count":    "%d alias",
		"tui_status_opened":    "abierto %s",
		"tui_status_copied":    "copiado %s",
		"tui_status_copyfail":  "fallo al copiar: %s",
		"tui_status_deleted":   "eliminado %s",
		"tui_status_saved":     "guardado",
		"tui_status_delfail":   "fallo al borrar: %s",
		"tui_status_tag_set":   "filtrando por #%s",
		"tui_status_tag_clear": "filtro de etiqueta limpiado",
		"tui_status_err":       "error: %s",
		"tui_empty":            "aún no hay alias\n\npulsa  a  para añadir el primero",
		"tui_no_matches":       "sin coincidencias",
		"tui_help_list":        " enter abrir · a añadir · e editar · d borrar · / filtrar · t etiqueta · y copiar · ? ayuda · q salir",
		"tui_help_form":        " tab siguiente · shift+tab anterior · enter guardar · esc cancelar",
		"tui_help_confirm":     " s/n · enter confirmar · esc cancelar",
		"tui_help_back":        " esc volver",
		"tui_form_add":         "Añadir alias",
		"tui_form_edit":        "Editar alias",
		"tui_field_name":       "nombre",
		"tui_field_url":        "url",
		"tui_field_tags":       "etiquetas",
		"tui_field_desc":       "descripción",
		"tui_confirm_delete":   "¿Borrar alias?",
		"tui_confirm_yes":      "[ Sí ]",
		"tui_confirm_no":       "[ No ]",
		"tui_help_title":       "teclado",
		"tui_opens":            "%d aperturas",
		"tui_last":             "última %s",
		"tui_created":          "creado %s",
		"tui_placeholder_name": "nombre del alias (p.ej. gh)",
		"tui_placeholder_url":  "url (p.ej. github.com)",
		"tui_placeholder_tags": "etiquetas separadas por comas",
		"tui_placeholder_desc": "descripción (opcional)",
		"err_empty_name":       "el nombre es obligatorio",
		"err_empty_url":        "la url es obligatoria",
		"err_invalid_name":     "el nombre no puede llevar espacios ni barras",

		// v0.3 — destinos multi-tipo + shell init + cambio de idioma
		"shellinit_short":           "Emite un wrapper de shell para que los alias de directorio cambien de directorio en el shell padre (sourcéalo en .zshrc/.bashrc/config.fish)",
		"err_shell_wrapper_missing": "goto: wrapper de shell no instalado — ejecuta: eval \"$(goto shell-init zsh)\"",
		"type_url":                  "url",
		"type_mailto":               "mail",
		"type_ssh":                  "ssh",
		"type_file":                 "archivo",
		"type_directory":            "dir",
		"type_command":              "cmd",
		"type_auto":                 "auto",
		"tui_lang_switched":         "idioma establecido a %s (permanente)",
		"help_lang":                 "cambiar idioma (se guarda)",

		// Onboarding inicial
		"onb_welcome_title":    "Bienvenido a goto ✦",
		"onb_welcome_body":     "Vas a configurar goto — un único comando para URLs, archivos, directorios, mail y hosts SSH.\nEste asistente de 30 segundos elige tus valores por defecto (todo editable luego en ~/.config/goto/config.toml).",
		"onb_welcome_tip":      "pulsa enter para empezar  ·  esc para saltar con los valores por defecto",
		"onb_lang_title":       "Idioma",
		"onb_lang_desc":        "Elige el idioma de la interfaz. 'auto' sigue $LANG.",
		"onb_theme_title":      "Tema",
		"onb_theme_desc":       "Todos los temas usan el mismo acento (#00B5E2); elige la paleta de fondo que prefieras.",
		"onb_browser_title":    "Navegador por defecto",
		"onb_browser_desc":     "Qué navegador abrirá los alias de tipo URL. 'default' respeta tu preferencia del SO.",
		"onb_action_title":     "Acción por defecto",
		"onb_action_desc":      "Qué pasa al ejecutar `goto X` si X no es un alias: autodetecta el tipo (recomendado), o fuerza url/file/directory.",
		"onb_firstalias_title": "Añade tu primer alias (opcional)",
		"onb_firstalias_desc":  "Prueba algo como `gh → github.com/aaangelmartin`. Siempre puedes añadir más con `goto add` o desde la TUI.",
		"onb_firstalias_skip":  "tab cambia de campo · enter guarda (o salta si está vacío) · esc salta",
		"onb_done_title":       "¡Listo! 🎉",
		"onb_done_body":        "Tip: para que los alias de directorio cambien tu shell, sourcea el wrapper una vez:\n    eval \"$(goto shell-init zsh)\"\n(sustituye zsh por bash o fish)",
		"onb_done_tip":         "pulsa enter para empezar",
		"onb_footer":           " ↑/↓ elegir · enter siguiente · ← volver · esc saltar asistente · q salir",

		// TUI — help rows
		"help_move":   "mover selección",
		"help_jump":   "ir al inicio / final",
		"help_open":   "abrir alias en el navegador",
		"help_filter": "filtrar lista (teclea para acotar)",
		"help_esc":    "limpiar filtro / volver",
		"help_add":    "añadir un alias",
		"help_edit":   "editar alias seleccionado",
		"help_delete": "borrar alias seleccionado",
		"help_yank":   "copiar URL al portapapeles",
		"help_tag":    "filtrar por la primera etiqueta del alias (toggle)",
		"help_toggle": "mostrar/ocultar esta ayuda",
		"help_quit":   "salir",
		"help_config": "config:",
		"help_theme":  "tema:",
	},
}
