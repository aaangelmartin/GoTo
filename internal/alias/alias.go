// Package alias defines the Alias model and lookup logic.
package alias

import (
	"encoding/json"
	"time"
)

// Type classifies what kind of target an alias points at.
//
// It determines which opener is invoked and how the TUI renders the alias.
type Type string

const (
	// TypeAuto means "detect at open time" based on Target content.
	TypeAuto Type = "auto"
	// TypeURL is a web URL (http/https or bare domain).
	TypeURL Type = "url"
	// TypeMailto is an email address (rendered as mailto:).
	TypeMailto Type = "mailto"
	// TypeSSH is an ssh:// URL or user@host[:port] short form.
	TypeSSH Type = "ssh"
	// TypeFile is a filesystem file path.
	TypeFile Type = "file"
	// TypeDirectory is a filesystem directory path (opens via shell wrapper if sourced).
	TypeDirectory Type = "directory"
	// TypeCommand is a raw shell command to execute.
	TypeCommand Type = "command"
)

// Alias represents a named shortcut to a target.
type Alias struct {
	Name        string    `json:"name"`
	Target      string    `json:"target"`
	Type        Type      `json:"type,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	LastOpened  time.Time `json:"last_opened,omitempty"`
	HitCount    int       `json:"hit_count"`
}

// UnmarshalJSON accepts both the legacy {"url": "..."} shape and the new
// {"target": "...", "type": "..."} shape, keeping existing aliases.json files
// working without a migration step.
func (a *Alias) UnmarshalJSON(b []byte) error {
	type wire struct {
		Name        string    `json:"name"`
		Target      string    `json:"target"`
		URL         string    `json:"url"` // legacy
		Type        Type      `json:"type,omitempty"`
		Tags        []string  `json:"tags,omitempty"`
		Description string    `json:"description,omitempty"`
		CreatedAt   time.Time `json:"created_at"`
		LastOpened  time.Time `json:"last_opened,omitempty"`
		HitCount    int       `json:"hit_count"`
	}
	var w wire
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}
	a.Name = w.Name
	a.Target = w.Target
	if a.Target == "" {
		a.Target = w.URL
	}
	a.Type = w.Type
	if a.Type == "" {
		a.Type = TypeAuto
	}
	a.Tags = w.Tags
	a.Description = w.Description
	a.CreatedAt = w.CreatedAt
	a.LastOpened = w.LastOpened
	a.HitCount = w.HitCount
	return nil
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
