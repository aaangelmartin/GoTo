package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aaangelmartin/goto/internal/alias"
	"github.com/aaangelmartin/goto/internal/i18n"
)

func (m *model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case openedMsg:
		if msg.err != nil {
			m.setStatus(i18n.Tf("tui_status_err", msg.err.Error()))
		} else {
			m.setStatus(i18n.Tf("tui_status_opened", msg.alias.Name))
			m.refresh()
		}
		return m, nil
	case tea.KeyMsg:
		if m.filterMode {
			return m.handleFilterKey(msg)
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.screen = screenHelp
			return m, nil
		case "j", "down":
			filtered := m.filteredItems()
			if m.cursor < len(filtered)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "g":
			m.cursor = 0
		case "G":
			if f := m.filteredItems(); len(f) > 0 {
				m.cursor = len(f) - 1
			}
		case "/":
			m.filterMode = true
		case "a":
			m.openForm(formAdd)
			return m, m.form.focusFirst()
		case "e":
			filtered := m.filteredItems()
			if m.cursor < len(filtered) {
				m.openForm(formEdit)
				m.form.loadFrom(filtered[m.cursor])
				return m, m.form.focusFirst()
			}
		case "d", "x":
			filtered := m.filteredItems()
			if m.cursor < len(filtered) {
				m.confirmTarget = filtered[m.cursor]
				m.screen = screenConfirm
			}
		case "y":
			filtered := m.filteredItems()
			if m.cursor < len(filtered) {
				if err := copyClipboard(filtered[m.cursor].URL); err != nil {
					m.setStatus(i18n.Tf("tui_status_copyfail", err.Error()))
				} else {
					m.setStatus(i18n.Tf("tui_status_copied", filtered[m.cursor].URL))
				}
			}
		case "t":
			filtered := m.filteredItems()
			if m.tagFilter != "" {
				m.tagFilter = ""
				m.setStatus(i18n.T("tui_status_tag_clear"))
			} else if m.cursor < len(filtered) && len(filtered[m.cursor].Tags) > 0 {
				m.tagFilter = filtered[m.cursor].Tags[0]
				m.cursor = 0
				m.setStatus(i18n.Tf("tui_status_tag_set", m.tagFilter))
			}
		case "enter":
			filtered := m.filteredItems()
			if m.cursor < len(filtered) {
				return m, m.openAlias(filtered[m.cursor])
			}
		case "esc":
			if m.filter != "" || m.tagFilter != "" {
				m.filter = ""
				m.tagFilter = ""
				m.cursor = 0
			}
		}
	}
	return m, nil
}

func (m *model) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.filterMode = false
		m.filter = ""
		m.cursor = 0
	case "enter":
		m.filterMode = false
	case "backspace":
		if len(m.filter) > 0 {
			m.filter = m.filter[:len(m.filter)-1]
			m.cursor = 0
		}
	default:
		if len(msg.Runes) > 0 {
			m.filter += string(msg.Runes)
			m.cursor = 0
		}
	}
	return m, nil
}

func (m *model) innerHeight() int {
	if m.height == 0 {
		return 20
	}
	return max(5, m.height-3) // header + footer
}

func (m *model) innerWidth() int {
	if m.width == 0 {
		return 80
	}
	return m.width
}

func (m *model) leftWidth() int {
	w := m.innerWidth()
	return max(20, w*2/5)
}
func (m *model) rightWidth() int {
	return max(20, m.innerWidth()-m.leftWidth())
}

func (m *model) listView() string {
	if len(m.items) == 0 {
		return m.theme.Box.
			Width(m.innerWidth()-2).
			Height(m.innerHeight()).
			Align(lipgloss.Center, lipgloss.Center).
			Render(m.theme.Empty.Render(i18n.T("tui_empty")))
	}

	filtered := m.filteredItems()
	if len(filtered) == 0 {
		return m.theme.Box.
			Width(m.innerWidth()-2).
			Height(m.innerHeight()).
			Align(lipgloss.Center, lipgloss.Center).
			Render(m.theme.Empty.Render(i18n.T("tui_no_matches")))
	}

	if m.cursor >= len(filtered) {
		m.cursor = len(filtered) - 1
	}

	leftW := m.leftWidth()
	rightW := m.rightWidth()
	h := m.innerHeight()

	left := m.renderList(filtered, leftW, h)
	right := m.renderPreview(filtered[m.cursor], rightW, h)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m *model) renderList(items []alias.Alias, width, height int) string {
	// scroll window
	innerH := height - 2 // borders
	if innerH < 1 {
		innerH = 1
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+innerH {
		m.offset = m.cursor - innerH + 1
	}
	end := m.offset + innerH
	if end > len(items) {
		end = len(items)
	}

	var b strings.Builder
	for i := m.offset; i < end; i++ {
		a := items[i]
		line := fmt.Sprintf("%-20s  %s", truncate(a.Name, 20), truncate(a.URL, width-24))
		if i == m.cursor {
			line = "▶ " + line
			b.WriteString(m.theme.ItemSel.Width(width - 2).Render(line))
		} else {
			line = "  " + line
			b.WriteString(m.theme.Item.Render(line))
		}
		b.WriteString("\n")
	}
	return m.theme.BoxFocused.Width(width - 2).Height(height).Render(b.String())
}

func (m *model) renderPreview(a alias.Alias, width, height int) string {
	var b strings.Builder
	b.WriteString(m.theme.Title.Render(a.Name))
	b.WriteString("\n")
	b.WriteString(m.theme.URL.Render(a.URL))
	b.WriteString("\n\n")
	if len(a.Tags) > 0 {
		tags := make([]string, 0, len(a.Tags))
		for _, t := range a.Tags {
			tags = append(tags, m.theme.Tag.Render("#"+t))
		}
		b.WriteString(strings.Join(tags, " "))
		b.WriteString("\n\n")
	}
	if a.Description != "" {
		b.WriteString(m.theme.Desc.Render(wrap(a.Description, width-4)))
		b.WriteString("\n\n")
	}
	b.WriteString(m.theme.Status.Render(i18n.Tf("tui_opens", a.HitCount)))
	if !a.LastOpened.IsZero() {
		b.WriteString(m.theme.Status.Render("  ·  " + i18n.Tf("tui_last", a.LastOpened.Format("2006-01-02 15:04"))))
	}
	if !a.CreatedAt.IsZero() {
		b.WriteString("\n")
		b.WriteString(m.theme.Status.Render(i18n.Tf("tui_created", a.CreatedAt.Format("2006-01-02"))))
	}
	return m.theme.Box.Width(width - 2).Height(height).Render(b.String())
}

func truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if len(s) <= n {
		return s
	}
	if n < 3 {
		return s[:n]
	}
	return s[:n-1] + "…"
}

func wrap(s string, width int) string {
	if width <= 0 {
		return s
	}
	var b strings.Builder
	line := 0
	words := strings.Fields(s)
	for i, w := range words {
		if line+len(w)+1 > width {
			b.WriteString("\n")
			line = 0
		} else if i > 0 {
			b.WriteString(" ")
			line++
		}
		b.WriteString(w)
		line += len(w)
	}
	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
