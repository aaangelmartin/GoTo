package urlx

import (
	"fmt"
	"os/exec"
)

// Opener opens a URL. Platform implementations are in open_{darwin,linux,windows}.go.
type Opener interface {
	Open(url string, browser string) error
}

// Open opens the URL via the platform-default opener.
// If browser is non-empty and recognized, it's used; otherwise the OS default is used.
func Open(url, browser string) error {
	if url == "" {
		return fmt.Errorf("empty url")
	}
	return platformOpen(url, browser)
}

func runOpen(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return cmd.Process.Release()
}
