package alias

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect(t *testing.T) {
	dir := t.TempDir()
	fileName := filepath.Join(dir, "sample.txt")
	if err := os.WriteFile(fileName, []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		in   string
		want Type
	}{
		{"https://github.com", TypeURL},
		{"google.com", TypeURL},
		{"mailto:a@b.com", TypeMailto},
		{"a@b.com", TypeMailto},
		{"ssh://user@host", TypeSSH},
		// user@host without a TLD-looking domain leans SSH
		{"root@myserver", TypeSSH},
		// With a TLD it's ambiguous; we prefer mailto — users can prefix ssh://
		{"user@host.example.com", TypeMailto},
		{"file:///etc/hosts", TypeFile},
		{dir, TypeDirectory},
		{fileName, TypeFile},
		{"~", TypeDirectory},
		{"./nope/", TypeDirectory},
		{"./missing.txt", TypeFile},
		{"git commit -m 'x'", TypeCommand},
	}
	for _, c := range cases {
		if got := Detect(c.in); got != c.want {
			t.Errorf("Detect(%q) = %s, want %s", c.in, got, c.want)
		}
	}
}

func TestResolveAuto(t *testing.T) {
	a := Alias{Target: "https://x.com", Type: TypeAuto}
	if got := Resolve(a); got != TypeURL {
		t.Errorf("Resolve auto = %s, want %s", got, TypeURL)
	}
	b := Alias{Target: "/etc", Type: TypeURL} // explicit wins
	if got := Resolve(b); got != TypeURL {
		t.Errorf("explicit type should not be overridden")
	}
}
