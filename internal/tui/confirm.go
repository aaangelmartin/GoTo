package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aaangelmartin/goto/internal/i18n"
)

func (m *model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "y", "Y", "s", "S":
			m.confirmDelete()
			m.screen = screenList
		case "n", "N", "esc":
			m.screen = screenList
		case "enter":
			if m.confirmYes {
				m.confirmDelete()
			}
			m.screen = screenList
		case "left", "right", "tab":
			m.confirmYes = !m.confirmYes
		}
	}
	return m, nil
}

func (m *model) confirmDelete() {
	if err := m.store.Delete(m.confirmTarget.Name); err != nil {
		m.setStatus(i18n.Tf("tui_status_delfail", err.Error()))
		return
	}
	_ = m.store.Save()
	m.setStatus(i18n.Tf("tui_status_deleted", m.confirmTarget.Name))
	m.refresh()
}

func (m *model) confirmView() string {
	title := m.theme.Danger_.Render(i18n.T("tui_confirm_delete"))
	body := fmt.Sprintf("%s\n%s\n\n",
		m.theme.Title.Render(m.confirmTarget.Name),
		m.theme.URL.Render(m.confirmTarget.Target),
	)
	yes := i18n.T("tui_confirm_yes")
	no := i18n.T("tui_confirm_no")
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
