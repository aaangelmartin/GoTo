package i18n

import "testing"

func TestSetLang(t *testing.T) {
	SetLang("es")
	if T("short") == "" || T("short") == "short" {
		t.Fatalf("spanish translation missing")
	}
	if T("rm_aborted") != "cancelado" {
		t.Errorf("expected Spanish 'cancelado', got %q", T("rm_aborted"))
	}
	SetLang("en")
	if T("rm_aborted") != "aborted" {
		t.Errorf("expected English 'aborted', got %q", T("rm_aborted"))
	}
}

func TestSpanishResolution(t *testing.T) {
	SetLang("es")
	if got := T("added"); got == "" || got[0] == 'a' && got[1] == 'd' {
		// English is "added: ..."; Spanish should not start with "added"
		t.Errorf("expected Spanish translation, got %q", got)
	}
	SetLang("en")
}

func TestUnknownKeyReturnsKey(t *testing.T) {
	SetLang("en")
	if T("__totally_unknown__") != "__totally_unknown__" {
		t.Errorf("unknown key should round-trip")
	}
}

func TestTf(t *testing.T) {
	SetLang("en")
	got := Tf("added", "gh", "https://github.com")
	want := "added: gh -> https://github.com\n"
	if got != want {
		t.Errorf("Tf = %q, want %q", got, want)
	}
}

func TestCatalogParity(t *testing.T) {
	for k := range catalog[EN] {
		if _, ok := catalog[ES][k]; !ok {
			t.Errorf("missing ES translation for key %q", k)
		}
	}
	for k := range catalog[ES] {
		if _, ok := catalog[EN][k]; !ok {
			t.Errorf("missing EN translation for key %q", k)
		}
	}
}
