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

type formModel struct {
	inputs  [4]textinput.Model // name, url, tags, desc
	focused int
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

func (f *formModel) loadFrom(a alias.Alias) {
	f.editing = a.Name
	f.inputs[0].SetValue(a.Name)
	f.inputs[1].SetValue(a.URL)
	f.inputs[2].SetValue(strings.Join(a.Tags, ", "))
	f.inputs[3].SetValue(a.Description)
}

func (m *model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.screen = screenList
			return m, nil
		case "tab", "down":
			m.form.nextField()
		case "shift+tab", "up":
			m.form.prevField()
		case "enter":
			if m.form.focused == len(m.form.inputs)-1 {
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
		}
	}

	var cmd tea.Cmd
	m.form.inputs[m.form.focused], cmd = m.form.inputs[m.form.focused].Update(msg)
	return m, cmd
}

func (f *formModel) nextField() {
	f.inputs[f.focused].Blur()
	f.focused = (f.focused + 1) % len(f.inputs)
	f.inputs[f.focused].Focus()
}
func (f *formModel) prevField() {
	f.inputs[f.focused].Blur()
	f.focused = (f.focused - 1 + len(f.inputs)) % len(f.inputs)
	f.inputs[f.focused].Focus()
}

func (f *formModel) submit(st *store.Store) error {
	name := strings.TrimSpace(f.inputs[0].Value())
	rawURL := strings.TrimSpace(f.inputs[1].Value())
	if name == "" {
		return errEmptyName
	}
	if rawURL == "" {
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

	a := alias.Alias{
		Name:        name,
		URL:         urlx.Normalize(rawURL, false),
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
		b.WriteString("  ")
		b.WriteString(m.form.inputs[i].View())
		b.WriteString("\n\n")
	}
	if m.form.errMsg != "" {
		b.WriteString(m.theme.Danger_.Render("✗ " + i18n.T(m.form.errMsg)))
		b.WriteString("\n")
	}
	return m.theme.BoxFocused.Width(m.innerWidth() - 2).Height(m.innerHeight()).Render(b.String())
}
