package urlx

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct {
		in        string
		forceHTTP bool
		want      string
	}{
		{"google.com", false, "https://google.com"},
		{"google.com", true, "http://google.com"},
		{"https://github.com", false, "https://github.com"},
		{"http://example.com", false, "http://example.com"},
		{"mailto:foo@bar.com", false, "mailto:foo@bar.com"},
		{"file:///etc/hosts", false, "file:///etc/hosts"},
		{"ssh://user@host", false, "ssh://user@host"},
		{"  google.com  ", false, "https://google.com"},
		{"", false, ""},
	}
	for _, c := range cases {
		got := Normalize(c.in, c.forceHTTP)
		if got != c.want {
			t.Errorf("Normalize(%q, %v) = %q, want %q", c.in, c.forceHTTP, got, c.want)
		}
	}
}

func TestLooksLikeURL(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"google.com", true},
		{"https://x.com", true},
		{"claude code", false},
		{"hello", false},
		{"", false},
		{"mailto:a@b.c", true},
	}
	for _, c := range cases {
		if got := LooksLikeURL(c.in); got != c.want {
			t.Errorf("LooksLikeURL(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
