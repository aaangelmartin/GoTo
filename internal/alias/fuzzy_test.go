package alias

import "testing"

func TestRankOrdering(t *testing.T) {
	aliases := []Alias{
		{Name: "github"},
		{Name: "gitlab"},
		{Name: "google", Tags: []string{"search"}},
		{Name: "bitbucket", Description: "git hosting"},
	}
	matches := Rank("git", aliases, 0.2)
	if len(matches) == 0 {
		t.Fatalf("expected matches, got none")
	}
	if matches[0].Alias.Name != "github" {
		t.Errorf("expected github first, got %s", matches[0].Alias.Name)
	}
	for i := 1; i < len(matches); i++ {
		if matches[i-1].Score < matches[i].Score {
			t.Errorf("not sorted desc: %v", matches)
		}
	}
}

func TestRankExact(t *testing.T) {
	aliases := []Alias{
		{Name: "gh"},
		{Name: "github"},
	}
	matches := Rank("gh", aliases, 0.4)
	if matches[0].Score != 1.0 {
		t.Errorf("exact should score 1.0, got %f", matches[0].Score)
	}
	if matches[0].Alias.Name != "gh" {
		t.Errorf("expected gh first, got %s", matches[0].Alias.Name)
	}
}

func TestRankEmpty(t *testing.T) {
	if got := Rank("", []Alias{{Name: "x"}}, 0.4); got != nil {
		t.Errorf("empty query should return nil, got %v", got)
	}
}
