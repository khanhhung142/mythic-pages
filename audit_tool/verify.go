package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const judgeSystem = `You are a Vietnamese history/mythology fact-checker. Given a claim and search results, output a JSON verdict.

Output ONLY JSON, no markdown fences:
{
  "status": "verified|wrong|not_found|suspicious",
  "evidence": "<1-2 sentences: what the evidence says>",
  "source_url": "<most relevant URL from citations, empty if none>",
  "confidence": "high|medium|low"
}

Status rules:
- verified: evidence clearly supports the claim
- wrong: evidence contradicts the claim
- not_found: no evidence either way
- suspicious: evidence is partial or contradictory; needs human review`

func verifyClaim(claim Claim, verbose bool) (VerificationResult, error) {
	result := VerificationResult{Claim: claim}

	// Search
	sr, err := searchClaim(claim.Text)
	if err != nil {
		result.Status = "not_found"
		result.Evidence = fmt.Sprintf("Search failed: %v", err)
		result.Confidence = "low"
		return result, nil
	}

	if verbose {
		fmt.Printf("  search done, %d citations\n", len(sr.Citations))
	}

	// Judge
	citationList := strings.Join(sr.Citations, "\n")
	prompt := fmt.Sprintf(`Claim: %s
Entry's cited source: %s

Search results:
%s

Citations:
%s`, claim.Text, claim.Source, sr.Answer, citationList)

	raw, err := callClaude(judgeSystem, prompt, 512)
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
		result.Evidence = fmt.Sprintf("Parse error: %v\nRaw: %s", err, raw)
		result.Confidence = "low"
		return result, nil
	}

	result.Status = verdict.Status
	result.Evidence = verdict.Evidence
	result.SourceURL = verdict.SourceURL
	result.Confidence = verdict.Confidence

	return result, nil
}
