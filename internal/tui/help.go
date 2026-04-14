package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aaangelmartin/goto/internal/i18n"
)

func (m *model) updateHelp(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "esc", "q", "?":
			m.screen = screenList
		}
	}
	return m, nil
}

func (m *model) helpView() string {
	rows := [][2]string{
		{"↑/k, ↓/j", i18n.T("help_move")},
		{"g, G", i18n.T("help_jump")},
		{"enter", i18n.T("help_open")},
		{"/", i18n.T("help_filter")},
		{"esc", i18n.T("help_esc")},
		{"a", i18n.T("help_add")},
		{"e", i18n.T("help_edit")},
		{"d, x", i18n.T("help_delete")},
		{"y", i18n.T("help_yank")},
		{"t", i18n.T("help_tag")},
		{"?", i18n.T("help_toggle")},
		{"q, ctrl+c", i18n.T("help_quit")},
	}
	var b strings.Builder
	b.WriteString(m.theme.Title.Render(i18n.T("tui_help_title")))
	b.WriteString("\n\n")
	for _, r := range rows {
		b.WriteString(m.theme.Key.Render(padRight(r[0], 14)))
		b.WriteString("  ")
		b.WriteString(m.theme.Item.Render(r[1]))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(m.theme.Subtitle.Render(i18n.T("help_config")))
	b.WriteString(" ")
	b.WriteString(m.theme.Status.Render("goto config"))
	b.WriteString("   ")
	b.WriteString(m.theme.Subtitle.Render(i18n.T("help_theme")))
	b.WriteString(" ")
	b.WriteString(m.theme.Status.Render(m.theme.Name))
	return m.theme.BoxFocused.Width(m.innerWidth() - 2).Height(m.innerHeight()).Render(b.String())
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
