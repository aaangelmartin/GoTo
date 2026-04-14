package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
		{"↑/k, ↓/j", "move selection"},
		{"g, G", "jump to top / bottom"},
		{"enter", "open alias in browser"},
		{"/", "filter list (type to narrow)"},
		{"esc", "clear filter / go back"},
		{"a", "add a new alias"},
		{"e", "edit selected alias"},
		{"d, x", "delete selected alias"},
		{"y", "copy URL to clipboard"},
		{"t", "filter by selected alias' first tag (toggle)"},
		{"?", "toggle this help"},
		{"q, ctrl+c", "quit"},
	}
	var b strings.Builder
	b.WriteString(m.theme.Title.Render("keyboard"))
	b.WriteString("\n\n")
	for _, r := range rows {
		b.WriteString(m.theme.Key.Render(padRight(r[0], 14)))
		b.WriteString("  ")
		b.WriteString(m.theme.Item.Render(r[1]))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(m.theme.Subtitle.Render("config:"))
	b.WriteString(" ")
	b.WriteString(m.theme.Status.Render("goto config"))
	b.WriteString("   ")
	b.WriteString(m.theme.Subtitle.Render("theme:"))
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
