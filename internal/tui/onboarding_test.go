package tui

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aaangelmartin/goto/internal/config"
	"github.com/aaangelmartin/goto/internal/store"
)

// TestOnboardingThemeLivePreview walks the wizard to the Theme step and
// asserts that pressing ↓ actually swaps the model's active theme (the bug
// the user reported: "los temas no funcionan").
func TestOnboardingThemeLivePreview(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GOTO_CONFIG", filepath.Join(dir, "config.toml"))
	t.Setenv("GOTO_ALIASES", filepath.Join(dir, "aliases.json"))
	st := store.New(filepath.Join(dir, "aliases.json"))
	if err := st.Load(); err != nil {
		t.Fatal(err)
	}
	cfg := config.Default()
	m := newModel(st, cfg)
	if m.screen != screenOnboard {
		t.Fatalf("expected onboarding to start, got screen %d", m.screen)
	}

	// Step 0 (welcome) → step 1 (language): press enter.
	mAny, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mAny.(*model)
	// language → theme
	mAny, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mAny.(*model)
	if m.onboard.step != stepTheme {
		t.Fatalf("expected step=stepTheme, got %d", m.onboard.step)
	}

	initialTheme := m.theme.Name
	if initialTheme != "default" {
		t.Fatalf("expected default theme at start, got %s", initialTheme)
	}

	// Simulate ↓ on the Theme step — first move should switch theme.
	mAny, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = mAny.(*model)
	if m.theme.Name == initialTheme {
		t.Fatalf("theme did not change after ↓; still %s", m.theme.Name)
	}
	if m.theme.Name != themeChoices[1] {
		t.Errorf("expected theme=%s, got %s", themeChoices[1], m.theme.Name)
	}
}

// TestOnboardingThemeAccentDiffer asserts that distinct themes render with
// distinct accent colors — the regression the user reported was that all
// themes looked the same because they shared the #00B5E2 accent.
func TestOnboardingThemeAccentDiffer(t *testing.T) {
	seen := map[string]string{}
	for _, name := range themeChoices {
		th := themeByName(name)
		accent := string(th.Accent)
		if other, ok := seen[accent]; ok {
			t.Errorf("themes %s and %s share the same accent %s", other, name, accent)
		}
		seen[accent] = name
	}
}
