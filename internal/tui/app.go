// Package tui implements the interactive Bubble Tea TUI.
package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aaangelmartin/goto/internal/alias"
	"github.com/aaangelmartin/goto/internal/config"
	"github.com/aaangelmartin/goto/internal/store"
	"github.com/aaangelmartin/goto/internal/urlx"
)

// Run starts the TUI and blocks until exit.
func Run(st *store.Store, cfg config.Config) error {
	m := newModel(st, cfg)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}

type screen int

const (
	screenList screen = iota
	screenForm
	screenConfirm
	screenHelp
)

type formMode int

const (
	formAdd formMode = iota
	formEdit
)

type model struct {
	store *store.Store
	cfg   config.Config
	theme Theme
	keys  keyMap

	screen screen

	// list state
	items      []alias.Alias
	filter     string
	filterMode bool
	cursor     int
	offset     int
	tagFilter  string

	// form state
	form     formModel
	formKind formMode

	// confirm state
	confirmTarget alias.Alias
	confirmYes    bool

	// dimensions
	width  int
	height int

	// status line
	status    string
	statusExp time.Time
}

func newModel(st *store.Store, cfg config.Config) *model {
	th := themeByName(cfg.Theme)
	return &model{
		store:  st,
		cfg:    cfg,
		theme:  th,
		keys:   defaultKeys(),
		screen: screenList,
		items:  st.List(),
	}
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	switch m.screen {
	case screenList:
		return m.updateList(msg)
	case screenForm:
		return m.updateForm(msg)
	case screenConfirm:
		return m.updateConfirm(msg)
	case screenHelp:
		return m.updateHelp(msg)
	}
	return m, nil
}

func (m *model) View() string {
	header := m.headerView()
	body := ""
	footer := m.footerView()
	switch m.screen {
	case screenList:
		body = m.listView()
	case screenForm:
		body = m.formView()
	case screenConfirm:
		body = m.confirmView()
	case screenHelp:
		body = m.helpView()
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

// ----- helpers -----

func (m *model) setStatus(s string) {
	m.status = s
	m.statusExp = time.Now().Add(3 * time.Second)
}

func (m *model) filteredItems() []alias.Alias {
	q := strings.ToLower(strings.TrimSpace(m.filter))
	out := make([]alias.Alias, 0, len(m.items))
	for _, a := range m.items {
		if m.tagFilter != "" && !a.HasTag(m.tagFilter) {
			continue
		}
		if q == "" {
			out = append(out, a)
			continue
		}
		if strings.Contains(strings.ToLower(a.Name), q) ||
			strings.Contains(strings.ToLower(a.URL), q) ||
			strings.Contains(strings.ToLower(a.Description), q) {
			out = append(out, a)
			continue
		}
		for _, t := range a.Tags {
			if strings.Contains(strings.ToLower(t), q) {
				out = append(out, a)
				break
			}
		}
	}
	return out
}

func (m *model) refresh() {
	m.items = m.store.List()
	if m.cursor >= len(m.filteredItems()) {
		m.cursor = len(m.filteredItems()) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *model) headerView() string {
	title := m.theme.Title.Render(" goto ")
	sub := m.theme.Subtitle.Render(fmt.Sprintf("%d aliases", len(m.items)))
	filter := ""
	if m.filterMode || m.filter != "" {
		filter = m.theme.Status.Render(fmt.Sprintf("  /%s", m.filter))
	}
	if m.tagFilter != "" {
		filter += m.theme.Tag.Render(fmt.Sprintf("  [#%s]", m.tagFilter))
	}
	return title + "  " + sub + filter
}

func (m *model) footerView() string {
	if m.status != "" && time.Now().Before(m.statusExp) {
		return m.theme.Status.Render(" " + m.status)
	}
	switch m.screen {
	case screenList:
		return m.theme.Help.Render(" enter open · a add · e edit · d delete · / filter · t tag · y copy · ? help · q quit")
	case screenForm:
		return m.theme.Help.Render(" tab next · shift+tab prev · enter save · esc cancel")
	case screenConfirm:
		return m.theme.Help.Render(" y/n · enter confirm · esc cancel")
	case screenHelp:
		return m.theme.Help.Render(" esc back")
	}
	return ""
}

// openAlias invokes the configured browser and bumps hit counter.
func (m *model) openAlias(a alias.Alias) tea.Cmd {
	return func() tea.Msg {
		browser := m.cfg.Browser
		if err := urlx.Open(a.URL, browser); err != nil {
			return openedMsg{err: err}
		}
		a.HitCount++
		a.LastOpened = time.Now()
		m.store.Set(a)
		_ = m.store.Save()
		return openedMsg{alias: a}
	}
}

type openedMsg struct {
	alias alias.Alias
	err   error
}
