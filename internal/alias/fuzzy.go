package alias

import "strings"

// Match represents a candidate alias with a relevance score in [0,1].
type Match struct {
	Alias Alias
	Score float64
}

// Rank returns aliases ranked by fuzzy score against query, filtered by the
// minimum threshold. The best match is first. Exact name match scores 1.0,
// prefix 0.9, substring 0.75, subsequence scales by density.
func Rank(query string, aliases []Alias, threshold float64) []Match {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}
	var out []Match
	for _, a := range aliases {
		s := score(q, a)
		if s >= threshold {
			out = append(out, Match{Alias: a, Score: s})
		}
	}
	// simple insertion sort by score desc (lists are small)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].Score > out[j-1].Score; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

func score(q string, a Alias) float64 {
	name := strings.ToLower(a.Name)
	switch {
	case name == q:
		return 1.0
	case strings.HasPrefix(name, q):
		return 0.9
	case strings.Contains(name, q):
		return 0.75
	}
	// subsequence scoring on name, then tags
	if s := subseqScore(q, name); s > 0 {
		return 0.4 + 0.3*s // 0.4–0.7 band
	}
	for _, t := range a.Tags {
		if strings.Contains(strings.ToLower(t), q) {
			return 0.5
		}
	}
	if strings.Contains(strings.ToLower(a.Description), q) {
		return 0.3
	}
	return 0
}

// subseqScore returns 0 if q is not a subsequence of s, else a density score
// in (0, 1] where 1 means all query chars were contiguous.
func subseqScore(q, s string) float64 {
	if q == "" {
		return 0
	}
	qi := 0
	firstMatch, lastMatch := -1, -1
	for si := 0; si < len(s) && qi < len(q); si++ {
		if s[si] == q[qi] {
			if firstMatch < 0 {
				firstMatch = si
			}
			lastMatch = si
			qi++
		}
	}
	if qi < len(q) {
		return 0
	}
	span := lastMatch - firstMatch + 1
	if span <= 0 {
		return 1
	}
	return float64(len(q)) / float64(span)
}
