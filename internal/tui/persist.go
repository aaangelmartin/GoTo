package tui

import "github.com/aaangelmartin/goto/internal/config"

// persistLanguage writes the chosen language into config.toml so the choice
// survives future sessions. The in-memory cfg is updated in place.
func persistLanguage(cfg *config.Config, lang string) error {
	cfg.Language = lang
	path, err := config.ConfigPath()
	if err != nil {
		return err
	}
	return config.Save(path, *cfg)
}
