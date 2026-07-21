package main

// EntryBlock is one parsed entry section.
type EntryBlock struct {
	Kind    string
	Section string
	Content string
}

// Claim extracted from a specific block.
type Claim struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Source string `json:"source"`
	Type   string `json:"type"`
	Block  string `json:"block"`
	Field  string `json:"field"`
	Risk   string `json:"risk"`
}

// VerificationResult for one claim.
type VerificationResult struct {
	Claim      Claim  `json:"claim"`
	Status     string `json:"status"`
	Evidence   string `json:"evidence"`
	SourceURL  string `json:"source_url"`
	Confidence string `json:"confidence"`
}

// PatternIssue from regex scan.
type PatternIssue struct {
	Pass    int    `json:"pass"`
	Line    int    `json:"line"`
	Pattern string `json:"pattern"`
	Text    string `json:"text"`
	Fix     string `json:"fix"`
}

// BlockAudit stores the audit result for one block.
type BlockAudit struct {
	Block   EntryBlock
	Claims  []Claim
	Results []VerificationResult
}

// AuditReport final output.
type AuditReport struct {
	EntryPath    string             `json:"entry_path"`
	EntryTitle   string             `json:"entry_title"`
	EntryType    string             `json:"entry_type"`
	BlockAudits  []BlockAudit       `json:"block_audits"`
	SourceLinks  []SourceLinkResult `json:"source_links"`
	Patterns     []PatternIssue     `json:"patterns"`
	Verdict      string             `json:"verdict"`
	RejectReason string             `json:"reject_reason,omitempty"`
	Summary      string             `json:"summary"`
}

func (r *AuditReport) AllResults() []VerificationResult {
	var out []VerificationResult
	for _, ba := range r.BlockAudits {
		out = append(out, ba.Results...)
	}
	return out
}
