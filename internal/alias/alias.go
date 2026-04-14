// Package alias defines the Alias model and lookup logic.
package alias

import "time"

// Alias represents a named shortcut to a URL.
type Alias struct {
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Tags        []string  `json:"tags,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	LastOpened  time.Time `json:"last_opened,omitempty"`
	HitCount    int       `json:"hit_count"`
}

// HasTag reports whether the alias has the given tag (case-insensitive).
func (a Alias) HasTag(tag string) bool {
	for _, t := range a.Tags {
		if equalFold(t, tag) {
			return true
		}
	}
	return false
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
