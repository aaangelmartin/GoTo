// Package store persists aliases to a JSON file atomically.
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aaangelmartin/goto/internal/alias"
)

// ErrNotFound is returned when a requested alias does not exist.
var ErrNotFound = errors.New("alias not found")

// ErrExists is returned when adding an alias that already exists.
var ErrExists = errors.New("alias already exists")

// Store manages a persistent collection of aliases.
type Store struct {
	Path    string
	aliases map[string]alias.Alias
}

// New returns a Store backed by path (not yet loaded).
func New(path string) *Store {
	return &Store{Path: path, aliases: map[string]alias.Alias{}}
}

// Load reads the file into memory. Missing file is treated as empty.
func (s *Store) Load() error {
	b, err := os.ReadFile(s.Path)
	if errors.Is(err, fs.ErrNotExist) {
		s.aliases = map[string]alias.Alias{}
		return nil
	}
	if err != nil {
		return err
	}
	if len(b) == 0 {
		s.aliases = map[string]alias.Alias{}
		return nil
	}
	var list []alias.Alias
	if err := json.Unmarshal(b, &list); err != nil {
		return fmt.Errorf("parse %s: %w", s.Path, err)
	}
	s.aliases = make(map[string]alias.Alias, len(list))
	for _, a := range list {
		s.aliases[strings.ToLower(a.Name)] = a
	}
	return nil
}

// Save writes all aliases atomically (temp file + rename).
func (s *Store) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return err
	}
	list := s.List()
	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(s.Path), ".goto-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := tmp.Write(b); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, s.Path)
}

// Get returns the alias by (case-insensitive) name.
func (s *Store) Get(name string) (alias.Alias, error) {
	a, ok := s.aliases[strings.ToLower(name)]
	if !ok {
		return alias.Alias{}, ErrNotFound
	}
	return a, nil
}

// Put inserts a new alias. Returns ErrExists if name is taken.
func (s *Store) Put(a alias.Alias) error {
	key := strings.ToLower(a.Name)
	if _, ok := s.aliases[key]; ok {
		return ErrExists
	}
	s.aliases[key] = a
	return nil
}

// Set inserts or replaces an alias (upsert).
func (s *Store) Set(a alias.Alias) {
	s.aliases[strings.ToLower(a.Name)] = a
}

// Delete removes an alias by name.
func (s *Store) Delete(name string) error {
	key := strings.ToLower(name)
	if _, ok := s.aliases[key]; !ok {
		return ErrNotFound
	}
	delete(s.aliases, key)
	return nil
}

// List returns all aliases sorted by name.
func (s *Store) List() []alias.Alias {
	out := make([]alias.Alias, 0, len(s.aliases))
	for _, a := range s.aliases {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Len returns the number of aliases.
func (s *Store) Len() int { return len(s.aliases) }
