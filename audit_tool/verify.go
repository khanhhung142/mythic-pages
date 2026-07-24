package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

const judgeSystem = `You are a Vietnamese history and mythology fact-checker. Given a claim, its context, and search results, output a JSON verdict.

Output ONLY JSON, no markdown fences:
{
  "status": "verified|wrong|not_found|suspicious|invented",
  "evidence": "<1-2 sentences: what the evidence says>",
  "source_url": "<most relevant URL from citations, empty if none>",
  "confidence": "high|medium|low"
}

Status definitions:
- verified: evidence clearly supports the claim
- wrong: evidence contradicts the claim
- not_found: no evidence either way
- suspicious: evidence is partial or contradictory; human review needed
- invented: strong signal this detail doesn't exist

HARD RULES (violating any of these is a critical failure):
- For status "wrong" or "invented", the evidence field MUST contain a verbatim quote from the search results that contradicts the claim. If you cannot quote the contradicting passage, use not_found or suspicious instead.
- NEVER mention a name, title, date, or fact that does not appear in the search results above. Do not fill gaps from your own knowledge.
- confidence "high" is only allowed when the evidence field quotes the search results directly.
- Absence from a single search is weak evidence. A person or work missing from these results is at most not_found with confidence medium, NEVER invented with high confidence.

SOURCE HIERARCHY (this is a scholarly encyclopedia; not all contradictions count):
- Wikipedia, newspapers (báo/tin tức), blogs, forums, and content farms CANNOT overturn a claim. If the ONLY contradicting evidence comes from such sources, the status is at most "suspicious", never "wrong" or "invented". Overturning a cited academic source requires an equally reputable source (journal/DOI, university/library, NXB, archive).
- When the entry cites a specific scholarly source and the search surfaces only popular sources, prefer "not_found" or "suspicious".

FOLKLORE VARIANT RULE (most false "wrong" verdicts come from here):
- Vietnamese legends exist in many regional and textual variants. Place names, market names, character names, and plot details LEGITIMATELY differ between versions.
- If the search shows a DIFFERENT variant (e.g. the market is called "Hà Thám" while the claim says "Hà Thị"), that is NOT a contradiction. Mark "suspicious" (variant), never "wrong". Only mark "wrong" if a source explicitly states the claim's specific version is false or erroneous.

INTERPRETIVE-CLAIM RULE:
- Analytical or interpretive statements (about origins, dating of a concept, scholarly framing, "product of modern studies", symbolism) are arguments, not checkable facts. NEVER mark these "wrong". At most "suspicious" if reputable scholarship visibly disagrees. Reserve "wrong"/"invented" for hard, checkable facts: dates, coordinates, Hán tự, source existence, ATU codes.

Special rules by claim type:
- source-existence: mark invented ONLY if the search results positively state the work does not exist or positively attribute it to a different author (quote that statement). Vietnamese author names are often romanized differently (given-name order, dropped diacritics); a name mismatch alone is not contradiction.
- han-tu: if characters don't match the meaning claimed, mark wrong
- atu: if code doesn't exist in ATU index, mark invented; if code exists but description doesn't match, mark wrong
- chuyen-ke narrative: be skeptical of very specific details absent from known versions`

func verifyClaim(claim Claim, verbose bool, rt Runtime) (VerificationResult, error) {
	result, err := judgeOnce(claim, buildSearchQuery(claim), verbose, rt)
	if err != nil {
		return result, err
	}

	// ponytail: second opinion only for damning verdicts; the judge only sees a
	// search summary, so a single pass hallucinates false wrong/invented.
	// Upgrade path: fetch top citation page content if this still misfires.
	if result.Status == "wrong" || result.Status == "invented" {
		second, err := judgeOnce(claim, altSearchQuery(claim), verbose, rt)
		if err == nil && second.Status == "verified" {
			result.Status = "suspicious"
			result.Evidence = fmt.Sprintf("Conflicting verdicts across two searches. First (%s): %s | Second (verified): %s",
				result.Status, result.Evidence, second.Evidence)
			result.SourceURL = second.SourceURL
			result.Confidence = "low"
		}
	}

	return gateDamningVerdict(result), nil
}

// gateDamningVerdict is the mechanical backstop for false "wrong"/"invented":
// such a verdict is only allowed to stand when the judge's own cited source is a
// reputable scholarly host. Wikipedia, newspapers, and blogs cannot overturn a
// cited academic source (skill §2.7), and observed false positives all rested on
// exactly those. Downgrade the rest to "suspicious" so a human reviews instead of
// the tool auto-rejecting on weak evidence. Independent of prompt compliance.
func gateDamningVerdict(r VerificationResult) VerificationResult {
	if r.Status != "wrong" && r.Status != "invented" {
		return r
	}
	if reputableVerdictSource(r.SourceURL) {
		return r
	}
	r.Evidence = "[auto-downgraded: contradicting source is not a reputable scholarly host — verify by hand] " + r.Evidence
	r.Status = "suspicious"
	r.Confidence = "low"
	return r
}

func reputableVerdictSource(rawURL string) bool {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return false
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return false
	}
	host := strings.ToLower(u.Host)
	return isTrustedSourceHost(host) && !isBannedSourceHost(host)
}

// isFatalBackendErr reports whether an "error"-status evidence string looks
// like a non-recoverable backend failure (quota/auth), so the pipeline can
// stop hammering a dead API instead of burning ~250 paid calls per bad run.
func isFatalBackendErr(evidence string) bool {
	e := strings.ToLower(evidence)
	for _, sig := range []string{"quota", "billing", "unauthorized", "invalid api key", "invalid_api_key", " 401", " 403", "exceeded"} {
		if strings.Contains(e, sig) {
			return true
		}
	}
	return false
}

func judgeOnce(claim Claim, searchQuery string, verbose bool, rt Runtime) (VerificationResult, error) {
	result := VerificationResult{Claim: claim}

	sr, err := rt.Search.Search(context.Background(), searchQuery)
	if err != nil {
		// Infra failure, NOT evidence of absence. Distinct status so a dead
		// backend can't masquerade as a benign not_found and produce a
		// confident verdict on an audit that checked nothing.
		result.Status = "error"
		result.Evidence = fmt.Sprintf("Search failed: %v", err)
		result.Confidence = "low"
		return result, nil
	}

	if verbose {
		fmt.Printf("    search done, %d citations\n", len(sr.Citations))
	}

	citationList := strings.Join(sr.Citations, "\n")
	prompt := fmt.Sprintf(`Claim: %s

Context: This claim is from the %s block, field/section "%s" of a Vietnamese mythology encyclopedia entry.
Claim type: %s
Entry's cited source: %s

Search results:
%s

Citations:
%s`, claim.Text, claim.Block, claim.Field, claim.Type, claim.Source, sr.Answer, citationList)

	raw, err := callLLMJSON(rt.Judge, judgeSystem, prompt, 512)
	if err != nil {
		result.Status = "error"
		result.Evidence = fmt.Sprintf("Judgment failed: %v", err)
		result.Confidence = "low"
		return result, nil
	}

	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var verdict struct {
		Status     string `json:"status"`
		Evidence   string `json:"evidence"`
		SourceURL  string `json:"source_url"`
		Confidence string `json:"confidence"`
	}
	if err := json.Unmarshal([]byte(raw), &verdict); err != nil {
		result.Status = "suspicious"
		result.Evidence = fmt.Sprintf("Parse error: %v | Raw: %s", err, truncate(raw, 200))
		result.Confidence = "low"
		return result, nil
	}

	result.Status = verdict.Status
	result.Evidence = verdict.Evidence
	result.SourceURL = verdict.SourceURL
	result.Confidence = verdict.Confidence

	return result, nil
}

func buildSearchQuery(c Claim) string {
	switch c.Type {
	case "source-existence":
		// Include the entry's own citation: it carries the exact title/journal,
		// which survives author-name romanization differences.
		if c.Source != "" {
			return fmt.Sprintf("%s %s", c.Text, c.Source)
		}
		return c.Text
	case "han-tu":
		return fmt.Sprintf("Vietnamese %s Han tu characters meaning", c.Text)
	case "atu":
		return fmt.Sprintf("ATU folktale type %s Thompson motif index", c.Text)
	case "person", "date":
		return fmt.Sprintf("Vietnamese history %s", c.Text)
	case "motif":
		return fmt.Sprintf("Vietnamese mythology %s motif", c.Text)
	case "place":
		return fmt.Sprintf("Vietnam %s location history", c.Text)
	default:
		base := c.Text
		if c.Source != "" {
			base = fmt.Sprintf("%s (source: %s)", base, c.Source)
		}
		return fmt.Sprintf("Vietnamese mythology history %s", base)
	}
}

// altSearchQuery gives a differently-phrased query for the double-check pass.
func altSearchQuery(c Claim) string {
	if c.Source != "" {
		return fmt.Sprintf("%s %s", c.Source, c.Text)
	}
	return c.Text
}
