package tui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/aaangelmartin/goto/internal/alias"
)

// Theme holds all styles for the TUI. Switch themes via config.
type Theme struct {
	Name    string
	Accent  lipgloss.Color
	Accent2 lipgloss.Color
	FG      lipgloss.Color
	FGDim   lipgloss.Color
	BG      lipgloss.Color
	BGDim   lipgloss.Color
	Success lipgloss.Color
	Danger  lipgloss.Color
	Warning lipgloss.Color
	Border  lipgloss.Color

	Title      lipgloss.Style
	Subtitle   lipgloss.Style
	Item       lipgloss.Style
	ItemSel    lipgloss.Style
	Tag        lipgloss.Style
	Key        lipgloss.Style
	Desc       lipgloss.Style
	Box        lipgloss.Style
	BoxFocused lipgloss.Style
	Status     lipgloss.Style
	Danger_    lipgloss.Style
	URL        lipgloss.Style
	Help       lipgloss.Style
	Empty      lipgloss.Style

	// TypeStyles maps each alias type to a color-coded badge style.
	TypeStyles map[alias.Type]lipgloss.Style
}

// TypeBadge returns a short label wrapped in the type's badge style.
func (t Theme) TypeBadge(typ alias.Type, label string) string {
	st, ok := t.TypeStyles[typ]
	if !ok {
		st = lipgloss.NewStyle().Foreground(t.FGDim)
	}
	return st.Render(" " + label + " ")
}

func themeByName(name string) Theme {
	// Each theme uses its own signature accent so they read as distinct
	// palettes. The project brand cyan (#00B5E2) is kept as the "default"
	// theme and appears as a secondary accent in a few others.
	switch name {
	case "dracula":
		return build("dracula",
			"#FF79C6", "#8BE9FD", "#F8F8F2", "#6272A4",
			"#282A36", "#44475A", "#50FA7B", "#FF5555", "#F1FA8C", "#44475A")
	case "catppuccin", "catppuccin-mocha":
		return build("catppuccin",
			"#CBA6F7", "#89DCEB", "#CDD6F4", "#6C7086",
			"#1E1E2E", "#313244", "#A6E3A1", "#F38BA8", "#F9E2AF", "#45475A")
	case "nord":
		return build("nord",
			"#88C0D0", "#81A1C1", "#ECEFF4", "#4C566A",
			"#2E3440", "#3B4252", "#A3BE8C", "#BF616A", "#EBCB8B", "#434C5E")
	case "tokyonight":
		return build("tokyonight",
			"#BB9AF7", "#7AA2F7", "#C0CAF5", "#565F89",
			"#1A1B26", "#24283B", "#9ECE6A", "#F7768E", "#E0AF68", "#3B4261")
	default:
		return build("default",
			"#00B5E2", "#7DD3FC", "#E5E7EB", "#6B7280",
			"#0B1220", "#1F2937", "#22D3A6", "#F87171", "#FBBF24", "#273043")
	}
}

func build(name, accent, accent2, fg, fgDim, bg, bgDim, ok, bad, warn, border string) Theme {
	t := Theme{
		Name:    name,
		Accent:  lipgloss.Color(accent),
		Accent2: lipgloss.Color(accent2),
		FG:      lipgloss.Color(fg),
		FGDim:   lipgloss.Color(fgDim),
		BG:      lipgloss.Color(bg),
		BGDim:   lipgloss.Color(bgDim),
		Success: lipgloss.Color(ok),
		Danger:  lipgloss.Color(bad),
		Warning: lipgloss.Color(warn),
		Border:  lipgloss.Color(border),
	}
	t.Title = lipgloss.NewStyle().Bold(true).Foreground(t.Accent)
	t.Subtitle = lipgloss.NewStyle().Foreground(t.Accent2)
	t.Item = lipgloss.NewStyle().Foreground(t.FG)
	t.ItemSel = lipgloss.NewStyle().Bold(true).Foreground(t.Accent).Background(t.BGDim)
	t.Tag = lipgloss.NewStyle().Foreground(t.Success)
	t.Key = lipgloss.NewStyle().Bold(true).Foreground(t.Accent2)
	t.Desc = lipgloss.NewStyle().Foreground(t.FGDim).Italic(true)
	t.Box = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(t.Border).Padding(0, 1)
	t.BoxFocused = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(t.Accent).Padding(0, 1)
	t.Status = lipgloss.NewStyle().Foreground(t.FGDim)
	t.Danger_ = lipgloss.NewStyle().Foreground(t.Danger).Bold(true)
	t.URL = lipgloss.NewStyle().Foreground(t.Accent2).Underline(true)
	t.Help = lipgloss.NewStyle().Foreground(t.FGDim)
	t.Empty = lipgloss.NewStyle().Foreground(t.FGDim).Italic(true).Align(lipgloss.Center)

	mk := func(fg string) lipgloss.Style {
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(fg))
	}
	t.TypeStyles = map[alias.Type]lipgloss.Style{
		alias.TypeURL:       mk("#00B5E2"),
		alias.TypeMailto:    mk("#F472B6"),
		alias.TypeSSH:       mk("#FBBF24"),
		alias.TypeFile:      mk("#A78BFA"),
		alias.TypeDirectory: mk("#22D3A6"),
		alias.TypeCommand:   mk("#FB7185"),
		alias.TypeAuto:      mk("#94A3B8"),
	}
	return t
}
