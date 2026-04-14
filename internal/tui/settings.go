package tui

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aaangelmartin/goto/internal/config"
	"github.com/aaangelmartin/goto/internal/i18n"
)

// settingsRow kinds.
type settingsRowKind int

const (
	rowCycle    settingsRowKind = iota // left/right cycles a fixed choice list
	rowText                            // enter toggles edit mode on a textinput
	rowExtEntry                        // existing extension row (.pdf → Preview) editable + deletable
	rowExtAdd                          // "+ Add extension opener" action row
)

type settingsRow struct {
	kind    settingsRowKind
	label   string
	key     string   // config key this row writes (e.g. "language", ".pdf")
	choices []string // for rowCycle
}

type settingsModel struct {
	rows     []settingsRow
	cursor   int
	editing  bool
	input    textinput.Model
	addInput textinput.Model // used by rowExtAdd to type a new ".ext app" pair
}

// openSettings builds the row list from the current config and switches the
// screen. Called from the list screen when the user presses 'o'.
func (m *model) openSettings() {
	m.settings = newSettingsModel(m.cfg, m.theme)
	m.screen = screenSettings
}

var (
	langChoices = []string{"auto", "en", "es"}
	dirChoices  = []string{"shell", "finder"}
)

func newSettingsModel(cfg config.Config, th Theme) settingsModel {
	rows := []settingsRow{
		{kind: rowCycle, label: "language", key: "language", choices: langChoices},
		{kind: rowCycle, label: "theme", key: "theme", choices: themeChoices},
		{kind: rowCycle, label: "browser", key: "browser", choices: browserChoices},
		{kind: rowCycle, label: "default_action", key: "default_action", choices: actionChoices},
		{kind: rowCycle, label: "directory_mode", key: "directory_mode", choices: dirChoices},
		{kind: rowText, label: "terminal (ssh)", key: "terminal"},
		{kind: rowText, label: "search_engine", key: "search_engine"},
		// Per-type openers.
		{kind: rowText, label: "opener · url", key: "openers.url"},
		{kind: rowText, label: "opener · mailto", key: "openers.mailto"},
		{kind: rowText, label: "opener · ssh", key: "openers.ssh"},
		{kind: rowText, label: "opener · file", key: "openers.file"},
		{kind: rowText, label: "opener · directory", key: "openers.directory"},
	}
	// Append extension-scoped openers (sorted for stable ordering).
	extKeys := make([]string, 0)
	for k := range cfg.Openers {
		if strings.HasPrefix(k, ".") {
			extKeys = append(extKeys, k)
		}
	}
	sort.Strings(extKeys)
	for _, k := range extKeys {
		rows = append(rows, settingsRow{kind: rowExtEntry, label: "ext · " + k, key: k})
	}
	rows = append(rows, settingsRow{kind: rowExtAdd, label: i18n.T("settings_add_ext"), key: ""})

	input := textinput.New()
	input.Prompt = "› "
	input.CharLimit = 256
	input.Width = 40
	input.PromptStyle = lipgloss.NewStyle().Foreground(th.Accent)

	add := textinput.New()
	add.Prompt = "› "
	add.CharLimit = 64
	add.Width = 40
	add.Placeholder = i18n.T("settings_add_ext_placeholder")
	add.PromptStyle = lipgloss.NewStyle().Foreground(th.Accent)

	return settingsModel{rows: rows, input: input, addInput: add}
}

// rowValue resolves the current value of a row from cfg.
func (m *model) rowValue(r settingsRow) string {
	cfg := &m.cfg
	switch r.key {
	case "language":
		return emptyFallback(cfg.Language, "auto")
	case "theme":
		return emptyFallback(cfg.Theme, "default")
	case "browser":
		return emptyFallback(cfg.Browser, "default")
	case "default_action":
		return emptyFallback(cfg.DefaultAction, "auto")
	case "directory_mode":
		return emptyFallback(cfg.DirectoryMode, "shell")
	case "terminal":
		return cfg.Terminal
	case "search_engine":
		return cfg.SearchEngine
	}
	if strings.HasPrefix(r.key, "openers.") {
		return cfg.Openers[strings.TrimPrefix(r.key, "openers.")]
	}
	// Extension row.
	return cfg.Openers[r.key]
}

func emptyFallback(v, fb string) string {
	if v == "" {
		return fb
	}
	return v
}

func (m *model) setRowValue(r settingsRow, v string) {
	cfg := &m.cfg
	switch r.key {
	case "language":
		cfg.Language = v
		if v == "auto" {
			i18n.SetLang("")
		} else {
			i18n.SetLang(v)
		}
	case "theme":
		cfg.Theme = v
		m.theme = themeByName(v)
	case "browser":
		cfg.Browser = v
	case "default_action":
		cfg.DefaultAction = v
	case "directory_mode":
		cfg.DirectoryMode = v
	case "terminal":
		cfg.Terminal = v
	case "search_engine":
		cfg.SearchEngine = v
	default:
		if cfg.Openers == nil {
			cfg.Openers = map[string]string{}
		}
		if strings.HasPrefix(r.key, "openers.") {
			cfg.Openers[strings.TrimPrefix(r.key, "openers.")] = v
		} else {
			// extension row
			if v == "" {
				delete(cfg.Openers, r.key)
			} else {
				cfg.Openers[r.key] = v
			}
		}
	}
}

func (m *model) cycleRow(r settingsRow, dir int) {
	cur := m.rowValue(r)
	idx := 0
	for i, c := range r.choices {
		if c == cur {
			idx = i
			break
		}
	}
	idx = (idx + dir + len(r.choices)) % len(r.choices)
	m.setRowValue(r, r.choices[idx])
}

// persistSettings writes the working copy of cfg to disk.
func (m *model) persistSettings() error {
	path, err := config.ConfigPath()
	if err != nil {
		return err
	}
	return config.Save(path, m.cfg)
}

func (m *model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	// Textinput edit mode handling first.
	if m.settings.editing {
		return m.updateSettingsEditing(km)
	}

	switch km.String() {
	case "esc", "q":
		if err := m.persistSettings(); err != nil {
			m.setStatus(i18n.Tf("tui_status_err", err.Error()))
		} else {
			m.setStatus(i18n.T("tui_status_saved"))
		}
		m.screen = screenList
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.settings.cursor > 0 {
			m.settings.cursor--
		}
	case "down", "j":
		if m.settings.cursor < len(m.settings.rows)-1 {
			m.settings.cursor++
		}
	case "left", "h":
		r := m.settings.rows[m.settings.cursor]
		if r.kind == rowCycle {
			m.cycleRow(r, -1)
		}
	case "right", "l":
		r := m.settings.rows[m.settings.cursor]
		if r.kind == rowCycle {
			m.cycleRow(r, +1)
		}
	case "enter", " ":
		r := m.settings.rows[m.settings.cursor]
		switch r.kind {
		case rowCycle:
			m.cycleRow(r, +1)
		case rowText, rowExtEntry:
			m.settings.input.SetValue(m.rowValue(r))
			m.settings.input.Focus()
			m.settings.editing = true
			return m, textinput.Blink
		case rowExtAdd:
			m.settings.addInput.SetValue("")
			m.settings.addInput.Focus()
			m.settings.editing = true
			return m, textinput.Blink
		}
	case "d", "x":
		// Delete an extension opener when positioned on one.
		r := m.settings.rows[m.settings.cursor]
		if r.kind == rowExtEntry {
			delete(m.cfg.Openers, r.key)
			// Rebuild rows so the deleted ext disappears immediately.
			m.settings = newSettingsModel(m.cfg, m.theme)
			if m.settings.cursor >= len(m.settings.rows) {
				m.settings.cursor = len(m.settings.rows) - 1
			}
		}
	}
	return m, nil
}

func (m *model) updateSettingsEditing(km tea.KeyMsg) (tea.Model, tea.Cmd) {
	r := m.settings.rows[m.settings.cursor]
	switch km.String() {
	case "esc":
		m.settings.editing = false
		m.settings.input.Blur()
		m.settings.addInput.Blur()
		return m, nil
	case "enter":
		if r.kind == rowExtAdd {
			parts := strings.Fields(m.settings.addInput.Value())
			if len(parts) >= 2 {
				ext := parts[0]
				if !strings.HasPrefix(ext, ".") {
					ext = "." + ext
				}
				app := strings.Join(parts[1:], " ")
				if m.cfg.Openers == nil {
					m.cfg.Openers = map[string]string{}
				}
				m.cfg.Openers[strings.ToLower(ext)] = app
				// Rebuild rows so the new ext appears.
				m.settings = newSettingsModel(m.cfg, m.theme)
			}
		} else {
			m.setRowValue(r, strings.TrimSpace(m.settings.input.Value()))
		}
		m.settings.editing = false
		m.settings.input.Blur()
		m.settings.addInput.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	if r.kind == rowExtAdd {
		m.settings.addInput, cmd = m.settings.addInput.Update(km)
	} else {
		m.settings.input, cmd = m.settings.input.Update(km)
	}
	return m, cmd
}

func (m *model) settingsView() string {
	th := m.theme
	var b strings.Builder
	b.WriteString(th.Title.Render(i18n.T("settings_title")))
	b.WriteString("\n")
	b.WriteString(th.Desc.Render(i18n.T("settings_desc")))
	b.WriteString("\n\n")

	for i, r := range m.settings.rows {
		label := r.label
		val := m.rowValue(r)
		if val == "" {
			val = "(empty)"
		}
		var line string
		switch r.kind {
		case rowCycle:
			line = th.Subtitle.Render(label) + "  " + th.Item.Render("← "+val+" →")
		case rowText, rowExtEntry:
			line = th.Subtitle.Render(label) + "  " + th.Item.Render(val)
		case rowExtAdd:
			line = th.Key.Render("＋ " + label)
		}
		if i == m.settings.cursor {
			b.WriteString(th.ItemSel.Render("▶ " + line))
			if m.settings.editing {
				b.WriteString("\n  ")
				if r.kind == rowExtAdd {
					b.WriteString(m.settings.addInput.View())
				} else {
					b.WriteString(m.settings.input.View())
				}
			}
		} else {
			b.WriteString("  " + line)
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	if m.settings.editing {
		b.WriteString(th.Help.Render(i18n.T("settings_editing_hint")))
	} else {
		b.WriteString(th.Help.Render(i18n.T("settings_hint")))
	}
	return th.BoxFocused.Width(m.innerWidth() - 2).Height(m.innerHeight()).Render(b.String())
}
