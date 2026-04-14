package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Open    key.Binding
	Add     key.Binding
	Edit    key.Binding
	Delete  key.Binding
	Filter  key.Binding
	Yank    key.Binding
	Tag     key.Binding
	Help    key.Binding
	Quit    key.Binding
	Escape  key.Binding
	Confirm key.Binding
}

func defaultKeys() keyMap {
	return keyMap{
		Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Open:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("↵", "open")),
		Add:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		Edit:    key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		Delete:  key.NewBinding(key.WithKeys("d", "x"), key.WithHelp("d", "delete")),
		Filter:  key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		Yank:    key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "copy url")),
		Tag:     key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "by tag")),
		Help:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Escape:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		Confirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("↵", "confirm")),
	}
}

// ShortHelp / FullHelp for bubbles/help integration
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Open, k.Add, k.Edit, k.Delete, k.Filter, k.Help, k.Quit}
}
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Open, k.Filter},
		{k.Add, k.Edit, k.Delete, k.Yank, k.Tag},
		{k.Help, k.Escape, k.Quit},
	}
}
