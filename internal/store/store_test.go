package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/aaangelmartin/goto/internal/alias"
)

func TestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "aliases.json")
	s := New(path)
	if err := s.Load(); err != nil {
		t.Fatalf("load empty: %v", err)
	}
	a := alias.Alias{Name: "gh", Target: "https://github.com", Type: alias.TypeURL, Tags: []string{"dev"}, CreatedAt: time.Now()}
	if err := s.Put(a); err != nil {
		t.Fatalf("put: %v", err)
	}
	if err := s.Put(a); err == nil {
		t.Errorf("put duplicate should fail")
	}
	if err := s.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	s2 := New(path)
	if err := s2.Load(); err != nil {
		t.Fatalf("reload: %v", err)
	}
	got, err := s2.Get("GH") // case-insensitive
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Target != a.Target {
		t.Errorf("target mismatch: %s", got.Target)
	}
	if err := s2.Delete("gh"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s2.Get("gh"); err == nil {
		t.Errorf("get after delete should fail")
	}
}

func TestLoadMissing(t *testing.T) {
	s := New(filepath.Join(t.TempDir(), "missing.json"))
	if err := s.Load(); err != nil {
		t.Fatalf("load missing should not error: %v", err)
	}
	if s.Len() != 0 {
		t.Errorf("expected empty, got %d", s.Len())
	}
}
