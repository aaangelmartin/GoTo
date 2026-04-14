package alias

// Result is the outcome of resolving a user-provided target.
type Result struct {
	// Kind indicates how the target was resolved.
	Kind ResultKind
	// URL is populated when Kind is KindURL, KindAlias, or KindFuzzy.
	URL string
	// Alias is populated when Kind is KindAlias or KindFuzzy (exact candidate).
	Alias Alias
	// Candidates is populated when Kind is KindAmbiguous.
	Candidates []Match
}

// ResultKind classifies a resolution result.
type ResultKind int

const (
	// KindURL means the target was treated as a raw URL (normalized).
	KindURL ResultKind = iota
	// KindAlias means the target matched an alias exactly.
	KindAlias
	// KindFuzzy means the target matched a single fuzzy candidate.
	KindFuzzy
	// KindAmbiguous means multiple fuzzy candidates matched.
	KindAmbiguous
	// KindNotFound means the target could not be resolved.
	KindNotFound
)
