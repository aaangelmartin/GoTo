package tui

import (
	"os"
	"os/exec"

	"github.com/aaangelmartin/goto/internal/config"
)

// openConfigInEditor opens the user's config.toml in $EDITOR if set, else
// returns the path without starting anything. The file is created with
// defaults if it doesn't exist so the editor never sees an empty buffer.
func openConfigInEditor() (string, error) {
	path, err := config.ConfigPath()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := config.Save(path, config.Default()); err != nil {
			return path, err
		}
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return path, nil
	}
	c := exec.Command(editor, path)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return path, c.Run()
}
