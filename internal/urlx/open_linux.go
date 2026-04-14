//go:build linux

package urlx

import "os/exec"

var linuxBrowserBins = map[string]string{
	"chrome":  "google-chrome",
	"firefox": "firefox",
	"brave":   "brave-browser",
	"edge":    "microsoft-edge",
	"opera":   "opera",
	"vivaldi": "vivaldi",
}

func platformOpen(url, browser string) error {
	if browser == "" || browser == "default" {
		return runOpen("xdg-open", url)
	}
	bin := linuxBrowserBins[browser]
	if bin == "" {
		bin = browser
	}
	if _, err := exec.LookPath(bin); err != nil {
		return runOpen("xdg-open", url)
	}
	return runOpen(bin, url)
}
