package alias

import (
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	emailRE = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	sshRE   = regexp.MustCompile(`^[a-zA-Z0-9._\-]+@[a-zA-Z0-9.\-]+(:\d+)?(/.+)?$`)
	protoRE = regexp.MustCompile(`^([a-zA-Z][a-zA-Z0-9+.\-]*):`)
)

// Detect returns the best guess for the type of target, performing a
// filesystem stat() if the target looks like a path. The result is never
// TypeAuto — Detect is the function that resolves "auto" to a concrete type.
func Detect(target string) Type {
	t := strings.TrimSpace(target)
	if t == "" {
		return TypeURL
	}

	// Windows drive letter (`C:\…`, `C:/…`) is not a URL scheme — route it
	// through the path detector so it can be stat()'d.
	if isWindowsDrivePath(t) {
		return detectPath(t)
	}

	// Explicit protocol prefixes.
	if strings.HasPrefix(t, "mailto:") {
		return TypeMailto
	}
	if strings.HasPrefix(t, "ssh://") {
		return TypeSSH
	}
	if strings.HasPrefix(t, "file://") {
		return TypeFile
	}
	if m := protoRE.FindStringSubmatch(t); m != nil {
		scheme := strings.ToLower(m[1])
		switch scheme {
		case "http", "https", "ftp", "ftps", "about", "chrome", "edge", "ws", "wss":
			return TypeURL
		case "mailto":
			return TypeMailto
		case "ssh", "sftp", "scp":
			return TypeSSH
		case "file":
			return TypeFile
		default:
			// Unknown protocol: treat as URL (opener will pass through).
			return TypeURL
		}
	}

	// Bare email address.
	if emailRE.MatchString(t) {
		return TypeMailto
	}

	// Path-like tokens: absolute, home, or explicit relative.
	if looksLikePath(t) {
		return detectPath(t)
	}

	// user@host[:port] short-form SSH (must contain @ and not be an email with dots-only).
	if strings.Contains(t, "@") && sshRE.MatchString(t) && !emailRE.MatchString(t) {
		return TypeSSH
	}

	// Bare domain or URL-looking — validate with net/url.
	if strings.Contains(t, ".") && !strings.ContainsAny(t, " \t") {
		if _, err := url.Parse("https://" + t); err == nil {
			return TypeURL
		}
	}

	// Fallback — treat as a command.
	if strings.ContainsAny(t, " \t") {
		return TypeCommand
	}
	return TypeURL
}

// Resolve returns the effective type for an alias, consulting Detect when the
// alias is marked TypeAuto.
func Resolve(a Alias) Type {
	if a.Type == "" || a.Type == TypeAuto {
		return Detect(a.Target)
	}
	return a.Type
}

func detectPath(t string) Type {
	expanded := expandHome(t)
	if st, err := os.Stat(expanded); err == nil {
		if st.IsDir() {
			return TypeDirectory
		}
		return TypeFile
	}
	// Target path doesn't exist; trailing separator ⇒ dir, else file.
	if strings.HasSuffix(t, string(os.PathSeparator)) || strings.HasSuffix(t, "/") || strings.HasSuffix(t, "\\") {
		return TypeDirectory
	}
	return TypeFile
}

func isWindowsDrivePath(t string) bool {
	if len(t) < 3 {
		return false
	}
	if !((t[0] >= 'a' && t[0] <= 'z') || (t[0] >= 'A' && t[0] <= 'Z')) {
		return false
	}
	return t[1] == ':' && (t[2] == '\\' || t[2] == '/')
}

func looksLikePath(t string) bool {
	if t == "" {
		return false
	}
	if strings.HasPrefix(t, "/") || strings.HasPrefix(t, "./") || strings.HasPrefix(t, "../") {
		return true
	}
	if strings.HasPrefix(t, "~") {
		return true
	}
	if len(t) >= 2 && t[1] == ':' && ((t[0] >= 'a' && t[0] <= 'z') || (t[0] >= 'A' && t[0] <= 'Z')) {
		// Windows-style drive letter (C:\foo)
		return true
	}
	return false
}

func expandHome(t string) string {
	if strings.HasPrefix(t, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			rest := strings.TrimPrefix(t, "~")
			return filepath.Join(home, rest)
		}
	}
	return t
}

// ExpandHome is the exported version for callers that need the absolute path.
func ExpandHome(t string) string { return expandHome(t) }
