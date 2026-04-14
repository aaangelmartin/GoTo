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
