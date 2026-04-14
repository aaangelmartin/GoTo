// Package urlx handles URL normalization and opening in the default browser.
package urlx

import (
	"regexp"
	"strings"
)

var protocolRE = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*:`)

// Normalize returns a URL with a protocol prefix suitable for the OS opener.
//
//	"google.com"          -> "https://google.com"
//	"http://example.com"  -> unchanged
//	"mailto:foo@bar.com"  -> unchanged
//	"file:///etc/hosts"   -> unchanged
//
// If forceHTTP is true, missing protocols become "http://" instead of "https://".
// An empty input returns an empty string.
func Normalize(target string, forceHTTP bool) string {
	t := strings.TrimSpace(target)
	if t == "" {
		return ""
	}
	if protocolRE.MatchString(t) {
		return t
	}
	if forceHTTP {
		return "http://" + t
	}
	return "https://" + t
}

// LooksLikeURL is a permissive heuristic: no whitespace and contains a dot
// or a known protocol. Used to decide whether a bare token should be opened
// as a URL or interpreted as a search query.
func LooksLikeURL(target string) bool {
	t := strings.TrimSpace(target)
	if t == "" {
		return false
	}
	if strings.ContainsAny(t, " \t\n") {
		return false
	}
	if protocolRE.MatchString(t) {
		return true
	}
	return strings.Contains(t, ".")
}
