package main

// Claim extracted from entry by Claude
type Claim struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`   // the claim as stated in the entry
	Source string `json:"source"` // citation the entry gave, empty if none
	Type   string `json:"type"`   // person | date | event | quote | motif | place
}

// VerificationResult for one claim
type VerificationResult struct {
	Claim      Claim  `json:"claim"`
	Status     string `json:"status"` // verified | wrong | not_found | suspicious
	Evidence   string `json:"evidence"`
	SourceURL  string `json:"source_url"`
	Confidence string `json:"confidence"` // high | medium | low
}

// PatternIssue from AI writing pattern scan
type PatternIssue struct {
	Pass    int    `json:"pass"`
	Line    int    `json:"line"`
	Pattern string `json:"pattern"`
	Text    string `json:"text"`
	Fix     string `json:"fix"`
}

// AuditReport final output
type AuditReport struct {
	EntryPath   string               `json:"entry_path"`
	EntryTitle  string               `json:"entry_title"`
	Claims      []Claim              `json:"claims"`
	Results     []VerificationResult `json:"results"`
	Patterns    []PatternIssue       `json:"patterns"`
	Verdict     string               `json:"verdict"` // PASS | REVISE | REJECT
	Summary     string               `json:"summary"`
}
