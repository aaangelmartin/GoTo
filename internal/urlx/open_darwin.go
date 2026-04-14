//go:build darwin

package urlx

var macBrowserAliases = map[string]string{
	"chrome":  "Google Chrome",
	"firefox": "Firefox",
	"safari":  "Safari",
	"arc":     "Arc",
	"brave":   "Brave Browser",
	"edge":    "Microsoft Edge",
	"opera":   "Opera",
	"vivaldi": "Vivaldi",
}

func platformOpen(url, browser string) error {
	if browser == "" || browser == "default" {
		return runOpen("open", url)
	}
	app := macBrowserAliases[browser]
	if app == "" {
		app = browser
	}
	return runOpen("open", "-a", app, url)
}
