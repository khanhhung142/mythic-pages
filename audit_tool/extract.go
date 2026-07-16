package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const extractSystem = `You extract factual claims from Vietnamese mythology encyclopedia entries.

Output ONLY a JSON array of claims. No preamble, no markdown fences.

Each claim:
{
  "id": <number>,
  "text": "<the claim, self-contained, in the language it appears>",
  "source": "<citation the entry gave, empty string if none>",
  "type": "<person|date|event|quote|motif|place>"
}

Rules:
- One atomic claim per object. Split compound claims.
- Skip interpretive/editorial statements ("this reflects...", "scholars have noted...") UNLESS they cite a named scholar — then include as type "quote".
- Include ALL: names with dates, chapter/tale numbers, Hán tự claims, ATU/Thompson codes, specific years, direct quotes from classical texts.
- Vague attributions ("nhiều nhà nghiên cứu") = include as type "quote" with source empty — these need flagging.
- Max 60 claims per entry. Prioritize verifiable specifics.`

func extractClaims(entryContent string) ([]Claim, error) {
	prompt := fmt.Sprintf("Extract all factual claims from this entry:\n\n%s", entryContent)

	raw, err := callClaude(extractSystem, prompt, 4096)
	if err != nil {
		return nil, fmt.Errorf("claim extraction: %w", err)
	}

	// Strip markdown fences if Claude added them anyway
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var claims []Claim
	if err := json.Unmarshal([]byte(raw), &claims); err != nil {
		return nil, fmt.Errorf("parse claims JSON: %w\nraw output: %s", err, raw)
	}

	return claims, nil
}
