//go:build windows

package urlx

import "os/exec"

var winBrowserBins = map[string]string{
	"chrome":  "chrome",
	"firefox": "firefox",
	"edge":    "msedge",
	"brave":   "brave",
	"opera":   "opera",
	"vivaldi": "vivaldi",
}

func platformOpen(url, browser string) error {
	if browser == "" || browser == "default" {
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	}
	bin := winBrowserBins[browser]
	if bin == "" {
		bin = browser
	}
	return exec.Command("cmd", "/c", "start", "", bin, url).Start()
}
