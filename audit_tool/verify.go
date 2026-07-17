package main

import (
	"encoding/json"
	"fmt"
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

Special rules by claim type:
- source-existence: mark invented ONLY if the search results positively state the work does not exist or positively attribute it to a different author (quote that statement). Vietnamese author names are often romanized differently (given-name order, dropped diacritics); a name mismatch alone is not contradiction.
- han-tu: if characters don't match the meaning claimed, mark wrong
- atu: if code doesn't exist in ATU index, mark invented; if code exists but description doesn't match, mark wrong
- chuyen-ke narrative: be skeptical of very specific details absent from known versions`

func verifyClaim(claim Claim, verbose bool) (VerificationResult, error) {
	result, err := judgeOnce(claim, buildSearchQuery(claim), verbose)
	if err != nil {
		return result, err
	}

	// ponytail: second opinion only for damning verdicts; the judge only sees a
	// search summary, so a single pass hallucinates false wrong/invented.
	// Upgrade path: fetch top citation page content if this still misfires.
	if result.Status == "wrong" || result.Status == "invented" {
		second, err := judgeOnce(claim, altSearchQuery(claim), verbose)
		if err == nil && second.Status == "verified" {
			result.Status = "suspicious"
			result.Evidence = fmt.Sprintf("Conflicting verdicts across two searches. First (%s): %s | Second (verified): %s",
				result.Status, result.Evidence, second.Evidence)
			result.SourceURL = second.SourceURL
			result.Confidence = "low"
		}
	}

	return result, nil
}

func judgeOnce(claim Claim, searchQuery string, verbose bool) (VerificationResult, error) {
	result := VerificationResult{Claim: claim}

	sr, err := searchClaim(searchQuery)
	if err != nil {
		result.Status = "not_found"
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

	raw, err := callClaudeJSON(judgeSystem, prompt, 512)
	if err != nil {
		result.Status = "not_found"
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
