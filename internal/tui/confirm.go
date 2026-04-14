package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "y", "Y":
			if err := m.store.Delete(m.confirmTarget.Name); err != nil {
				m.setStatus("delete failed: " + err.Error())
			} else {
				_ = m.store.Save()
				m.setStatus("deleted " + m.confirmTarget.Name)
				m.refresh()
			}
			m.screen = screenList
		case "n", "N", "esc":
			m.screen = screenList
		case "enter":
			if m.confirmYes {
				if err := m.store.Delete(m.confirmTarget.Name); err != nil {
					m.setStatus("delete failed: " + err.Error())
				} else {
					_ = m.store.Save()
					m.setStatus("deleted " + m.confirmTarget.Name)
					m.refresh()
				}
			}
			m.screen = screenList
		case "left", "right", "tab":
			m.confirmYes = !m.confirmYes
		}
	}
	return m, nil
}

func (m *model) confirmView() string {
	title := m.theme.Danger_.Render("Delete alias?")
	body := fmt.Sprintf("%s\n%s\n\n",
		m.theme.Title.Render(m.confirmTarget.Name),
		m.theme.URL.Render(m.confirmTarget.URL),
	)
	yes := "[ Yes ]"
	no := "[ No ]"
	if m.confirmYes {
		yes = m.theme.ItemSel.Render(yes)
		no = m.theme.Item.Render(no)
	} else {
		yes = m.theme.Item.Render(yes)
		no = m.theme.ItemSel.Render(no)
	}
	return m.theme.BoxFocused.
		Width(m.innerWidth()-2).
		Height(m.innerHeight()).
		Align(lipgloss.Center, lipgloss.Center).
		Render(title + "\n\n" + body + yes + "   " + no)
}
