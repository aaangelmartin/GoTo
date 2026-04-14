package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aaangelmartin/goto/internal/alias"
	"github.com/aaangelmartin/goto/internal/i18n"
	"github.com/aaangelmartin/goto/internal/store"
	"github.com/aaangelmartin/goto/internal/urlx"
)

// formTypes is the cycle order shown on the Type field. Order chosen to
// surface the most-used types first.
var formTypes = []alias.Type{
	alias.TypeAuto,
	alias.TypeURL,
	alias.TypeMailto,
	alias.TypeSSH,
	alias.TypeFile,
	alias.TypeDirectory,
	alias.TypeCommand,
}

type formModel struct {
	inputs  [4]textinput.Model // name, target, tags, desc
	typeIdx int                // index into formTypes; Type field is a selector, not a textinput
	focused int                // 0..3 = inputs, 4 = type selector
	errMsg  string
	editing string // non-empty means editing an existing alias (by name)
	theme   Theme
}

func (m *model) openForm(kind formMode) {
	m.formKind = kind
	m.form = newFormModel(m.theme)
	m.screen = screenForm
}

func newFormModel(th Theme) formModel {
	mk := func(placeholder string, width int) textinput.Model {
		ti := textinput.New()
		ti.Placeholder = placeholder
		ti.Width = width
		ti.CharLimit = 512
		ti.Prompt = "› "
		ti.PromptStyle = lipgloss.NewStyle().Foreground(th.Accent)
		ti.TextStyle = lipgloss.NewStyle().Foreground(th.FG)
		return ti
	}
	return formModel{
		inputs: [4]textinput.Model{
			mk(i18n.T("tui_placeholder_name"), 40),
			mk(i18n.T("tui_placeholder_url"), 60),
			mk(i18n.T("tui_placeholder_tags"), 40),
			mk(i18n.T("tui_placeholder_desc"), 60),
		},
		theme: th,
	}
}

func (f *formModel) focusFirst() tea.Cmd {
	f.focused = 0
	f.inputs[0].Focus()
	return textinput.Blink
}

// totalFields = 4 text inputs + 1 type selector.
const formFieldCount = 5

func (f *formModel) loadFrom(a alias.Alias) {
	f.editing = a.Name
	f.inputs[0].SetValue(a.Name)
	f.inputs[1].SetValue(a.Target)
	f.inputs[2].SetValue(strings.Join(a.Tags, ", "))
	f.inputs[3].SetValue(a.Description)
	// Position the type cycle on the alias's stored type (default to auto).
	f.typeIdx = 0
	for i, t := range formTypes {
		if t == a.Type {
			f.typeIdx = i
			break
		}
	}
}

func (m *model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "esc":
			m.screen = screenList
			return m, nil
		case "tab", "down":
			// On the Type field, Down cycles the type instead of moving focus
			// away — only Tab moves off the selector.
			if m.form.focused == 4 && km.String() == "down" {
				m.form.typeIdx = (m.form.typeIdx + 1) % len(formTypes)
				return m, nil
			}
			m.form.nextField()
			return m, nil
		case "shift+tab", "up":
			if m.form.focused == 4 && km.String() == "up" {
				m.form.typeIdx = (m.form.typeIdx - 1 + len(formTypes)) % len(formTypes)
				return m, nil
			}
			m.form.prevField()
			return m, nil
		case "enter":
			if m.form.focused == formFieldCount-1 {
				if err := m.form.submit(m.store); err != nil {
					m.form.errMsg = err.Error()
					return m, nil
				}
				m.refresh()
				m.setStatus(i18n.T("tui_status_saved"))
				m.screen = screenList
				return m, nil
			}
			m.form.nextField()
			return m, nil
		}
	}

	// Type selector isn't a textinput — ignore other keys when focused there.
	if m.form.focused == 4 {
		return m, nil
	}
	var cmd tea.Cmd
	m.form.inputs[m.form.focused], cmd = m.form.inputs[m.form.focused].Update(msg)
	return m, cmd
}

func (f *formModel) nextField() {
	if f.focused < 4 {
		f.inputs[f.focused].Blur()
	}
	f.focused = (f.focused + 1) % formFieldCount
	if f.focused < 4 {
		f.inputs[f.focused].Focus()
	}
}
func (f *formModel) prevField() {
	if f.focused < 4 {
		f.inputs[f.focused].Blur()
	}
	f.focused = (f.focused - 1 + formFieldCount) % formFieldCount
	if f.focused < 4 {
		f.inputs[f.focused].Focus()
	}
}

func (f *formModel) submit(st *store.Store) error {
	name := strings.TrimSpace(f.inputs[0].Value())
	rawTarget := strings.TrimSpace(f.inputs[1].Value())
	if name == "" {
		return errEmptyName
	}
	if rawTarget == "" {
		return errEmptyURL
	}
	if strings.ContainsAny(name, " \t/\\") {
		return errInvalidName
	}
	var tags []string
	for _, t := range strings.Split(f.inputs[2].Value(), ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}

	chosenType := formTypes[f.typeIdx]
	target := rawTarget
	// Only normalize when the user left the type on auto AND the detector
	// thinks it's a URL; otherwise preserve the raw string so paths, mail,
	// ssh short-form and commands survive the round-trip intact.
	if chosenType == alias.TypeAuto && alias.Detect(rawTarget) == alias.TypeURL {
		target = urlx.Normalize(rawTarget, false)
	} else if chosenType == alias.TypeURL {
		target = urlx.Normalize(rawTarget, false)
	}

	a := alias.Alias{
		Name:        name,
		Target:      target,
		Type:        chosenType,
		Tags:        tags,
		Description: strings.TrimSpace(f.inputs[3].Value()),
	}

	if f.editing != "" {
		orig, err := st.Get(f.editing)
		if err != nil {
			return err
		}
		a.CreatedAt = orig.CreatedAt
		a.HitCount = orig.HitCount
		a.LastOpened = orig.LastOpened
		if !strings.EqualFold(f.editing, name) {
			if err := st.Delete(f.editing); err != nil {
				return err
			}
		}
		st.Set(a)
	} else {
		a.CreatedAt = time.Now()
		if err := st.Put(a); err != nil {
			return err
		}
	}
	return st.Save()
}

type formErr string

func (e formErr) Error() string { return string(e) }

// These sentinels carry i18n keys; we translate when rendering.
const (
	errEmptyName   formErr = "err_empty_name"
	errEmptyURL    formErr = "err_empty_url"
	errInvalidName formErr = "err_invalid_name"
)

func (m *model) formView() string {
	labels := []string{
		i18n.T("tui_field_name"),
		i18n.T("tui_field_url"),
		i18n.T("tui_field_tags"),
		i18n.T("tui_field_desc"),
		i18n.T("tui_field_type"),
	}
	var b strings.Builder
	title := i18n.T("tui_form_add")
	if m.formKind == formEdit {
		title = i18n.T("tui_form_edit")
	}
	b.WriteString(m.theme.Title.Render(title))
	b.WriteString("\n\n")
	for i, l := range labels {
		label := m.theme.Subtitle.Render(l + ":")
		if i == m.form.focused {
			label = m.theme.Key.Render("› " + l + ":")
		} else {
			label = "  " + label
		}
		b.WriteString(label)
		b.WriteString("\n")

		if i == 4 {
			// Type selector row.
			t := formTypes[m.form.typeIdx]
			badge := m.theme.TypeBadge(t, strings.ToUpper(i18n.T("type_"+string(t))))
			if i == m.form.focused {
				b.WriteString("  ")
				b.WriteString(badge)
				b.WriteString("  ")
				b.WriteString(m.theme.Help.Render(i18n.T("tui_type_hint")))
			} else {
				b.WriteString("  ")
				b.WriteString(badge)
			}
		} else {
			b.WriteString("  ")
			b.WriteString(m.form.inputs[i].View())
		}
		b.WriteString("\n\n")
	}
	if m.form.errMsg != "" {
		b.WriteString(m.theme.Danger_.Render("✗ " + i18n.T(m.form.errMsg)))
		b.WriteString("\n")
	}
	return m.theme.BoxFocused.Width(m.innerWidth() - 2).Height(m.innerHeight()).Render(b.String())
}
