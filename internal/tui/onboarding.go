package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aaangelmartin/goto/internal/alias"
	"github.com/aaangelmartin/goto/internal/config"
	"github.com/aaangelmartin/goto/internal/i18n"
)

// onboardStep enumerates the pages of the first-run wizard. Keep this list
// short — anything fancier belongs in the full TUI, not in the welcome flow.
type onboardStep int

const (
	stepWelcome onboardStep = iota
	stepLanguage
	stepTheme
	stepBrowser
	stepAction
	stepFirstAlias
	stepDone
)

type onboardModel struct {
	step        onboardStep
	choiceIdx   int
	aliasInputs [2]textinput.Model // name, target
	aliasFocus  int
	theme       Theme
}

func newOnboardModel(th Theme) onboardModel {
	mk := func(ph string, w int) textinput.Model {
		ti := textinput.New()
		ti.Placeholder = ph
		ti.Width = w
		ti.CharLimit = 256
		ti.Prompt = "› "
		ti.PromptStyle = lipgloss.NewStyle().Foreground(th.Accent)
		return ti
	}
	return onboardModel{
		theme: th,
		aliasInputs: [2]textinput.Model{
			mk("gh", 30),
			mk("github.com/aaangelmartin", 60),
		},
	}
}

// shouldOnboard returns true when the installation looks brand new (no
// config file on disk AND no aliases yet). We check the file, not cfg,
// because Load() seeds defaults silently.
func shouldOnboard(store_ interface{ Len() int }) bool {
	cfgPath, err := config.ConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(cfgPath)
	noConfig := os.IsNotExist(err)
	noAliases := store_ == nil || store_.Len() == 0
	return noConfig && noAliases
}

// languageChoices / themeChoices / ... are the option lists rendered per step.
var (
	languageChoices = []string{"auto", "en", "es"}
	themeChoices    = []string{"default", "dracula", "catppuccin", "nord", "tokyonight"}
	browserChoices  = []string{"default", "chrome", "firefox", "safari", "arc", "brave", "edge"}
	actionChoices   = []string{"auto", "url", "file", "directory"}
)

func (m *model) stepOptions() []string {
	switch m.onboard.step {
	case stepLanguage:
		return languageChoices
	case stepTheme:
		return themeChoices
	case stepBrowser:
		return browserChoices
	case stepAction:
		return actionChoices
	}
	return nil
}

func (m *model) updateOnboard(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// First-alias step has text inputs, handle separately.
		if m.onboard.step == stepFirstAlias {
			return m.updateOnboardFirstAlias(msg)
		}
		switch msg.String() {
		case "esc":
			// Allow skipping the wizard entirely; caller saves whatever
			// defaults are in memory so we don't loop next launch.
			m.finishOnboarding()
			m.screen = screenList
			m.refresh()
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.onboard.choiceIdx > 0 {
				m.onboard.choiceIdx--
			}
		case "down", "j":
			if opts := m.stepOptions(); m.onboard.choiceIdx < len(opts)-1 {
				m.onboard.choiceIdx++
			}
		case "enter", "right", " ", "tab":
			m.advanceOnboardStep()
		case "left":
			m.retreatOnboardStep()
		}
	}
	return m, nil
}

func (m *model) updateOnboardFirstAlias(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.finishOnboarding()
		m.screen = screenList
		m.refresh()
		return m, nil
	case "tab", "down":
		m.onboard.aliasFocus = (m.onboard.aliasFocus + 1) % 2
		m.onboard.aliasInputs[0].Blur()
		m.onboard.aliasInputs[1].Blur()
		m.onboard.aliasInputs[m.onboard.aliasFocus].Focus()
	case "shift+tab", "up":
		m.onboard.aliasFocus = (m.onboard.aliasFocus + 1) % 2
		m.onboard.aliasInputs[0].Blur()
		m.onboard.aliasInputs[1].Blur()
		m.onboard.aliasInputs[m.onboard.aliasFocus].Focus()
	case "enter":
		name := strings.TrimSpace(m.onboard.aliasInputs[0].Value())
		target := strings.TrimSpace(m.onboard.aliasInputs[1].Value())
		if name != "" && target != "" {
			a := alias.Alias{Name: name, Target: target, Type: alias.TypeAuto}
			_ = m.store.Put(a)
			_ = m.store.Save()
		}
		m.finishOnboarding()
		m.screen = screenList
		m.refresh()
		return m, nil
	}
	var cmd tea.Cmd
	m.onboard.aliasInputs[m.onboard.aliasFocus], cmd = m.onboard.aliasInputs[m.onboard.aliasFocus].Update(msg)
	return m, cmd
}

func (m *model) advanceOnboardStep() {
	// Commit the current step's choice into cfg before moving on.
	switch m.onboard.step {
	case stepLanguage:
		lang := languageChoices[m.onboard.choiceIdx]
		m.cfg.Language = lang
		if lang != "auto" {
			i18n.SetLang(lang)
		}
	case stepTheme:
		m.cfg.Theme = themeChoices[m.onboard.choiceIdx]
		m.theme = themeByName(m.cfg.Theme)
	case stepBrowser:
		m.cfg.Browser = browserChoices[m.onboard.choiceIdx]
	case stepAction:
		m.cfg.DefaultAction = actionChoices[m.onboard.choiceIdx]
	case stepFirstAlias, stepDone:
		// nothing
	}
	m.onboard.choiceIdx = 0
	m.onboard.step++
	if m.onboard.step == stepFirstAlias {
		m.onboard.aliasInputs[0].Focus()
	}
	if m.onboard.step > stepDone {
		m.finishOnboarding()
		m.screen = screenList
		m.refresh()
	}
}

func (m *model) retreatOnboardStep() {
	if m.onboard.step > stepWelcome {
		m.onboard.step--
		m.onboard.choiceIdx = 0
	}
}

// finishOnboarding persists the chosen settings so the wizard doesn't run again.
func (m *model) finishOnboarding() {
	if path, err := config.ConfigPath(); err == nil {
		_ = config.Save(path, m.cfg)
	}
}

func (m *model) onboardView() string {
	th := m.theme
	w := m.innerWidth() - 2
	h := m.innerHeight()

	var body string
	switch m.onboard.step {
	case stepWelcome:
		body = renderWelcome(th)
	case stepLanguage:
		body = renderChoiceStep(th, i18n.T("onb_lang_title"), i18n.T("onb_lang_desc"), languageChoices, m.onboard.choiceIdx)
	case stepTheme:
		body = renderChoiceStep(th, i18n.T("onb_theme_title"), i18n.T("onb_theme_desc"), themeChoices, m.onboard.choiceIdx)
	case stepBrowser:
		body = renderChoiceStep(th, i18n.T("onb_browser_title"), i18n.T("onb_browser_desc"), browserChoices, m.onboard.choiceIdx)
	case stepAction:
		body = renderChoiceStep(th, i18n.T("onb_action_title"), i18n.T("onb_action_desc"), actionChoices, m.onboard.choiceIdx)
	case stepFirstAlias:
		body = m.renderFirstAlias()
	case stepDone:
		body = renderDone(th)
	}

	footer := th.Help.Render(i18n.T("onb_footer"))
	return th.BoxFocused.Width(w).Height(h).Render(body + "\n\n" + footer)
}

func renderWelcome(th Theme) string {
	title := th.Title.Render(i18n.T("onb_welcome_title"))
	body := th.Item.Render(i18n.T("onb_welcome_body"))
	tip := th.Subtitle.Render(i18n.T("onb_welcome_tip"))
	return title + "\n\n" + body + "\n\n" + tip
}

func renderChoiceStep(th Theme, title, desc string, opts []string, idx int) string {
	var b strings.Builder
	b.WriteString(th.Title.Render(title))
	b.WriteString("\n")
	b.WriteString(th.Desc.Render(desc))
	b.WriteString("\n\n")
	for i, o := range opts {
		if i == idx {
			b.WriteString(th.ItemSel.Render("▶ " + o))
		} else {
			b.WriteString(th.Item.Render("  " + o))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m *model) renderFirstAlias() string {
	var b strings.Builder
	th := m.theme
	b.WriteString(th.Title.Render(i18n.T("onb_firstalias_title")))
	b.WriteString("\n")
	b.WriteString(th.Desc.Render(i18n.T("onb_firstalias_desc")))
	b.WriteString("\n\n")
	b.WriteString(th.Subtitle.Render(i18n.T("tui_field_name") + ":"))
	b.WriteString("\n  ")
	b.WriteString(m.onboard.aliasInputs[0].View())
	b.WriteString("\n\n")
	b.WriteString(th.Subtitle.Render(i18n.T("tui_field_url") + ":"))
	b.WriteString("\n  ")
	b.WriteString(m.onboard.aliasInputs[1].View())
	b.WriteString("\n\n")
	b.WriteString(th.Help.Render(i18n.T("onb_firstalias_skip")))
	return b.String()
}

func renderDone(th Theme) string {
	return th.Title.Render(i18n.T("onb_done_title")) + "\n\n" +
		th.Item.Render(i18n.T("onb_done_body")) + "\n\n" +
		th.Subtitle.Render(i18n.T("onb_done_tip"))
}
