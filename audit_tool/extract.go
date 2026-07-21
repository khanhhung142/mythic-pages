package main

import (
	"encoding/json"
	"fmt"
)

const frontmatterExtractSystem = `You extract verifiable claims from the YAML frontmatter of a Vietnamese mythology encyclopedia entry.

Output ONLY a JSON array. No preamble, no markdown fences.

Each claim:
{
  "id": <number>,
  "text": "<self-contained claim in English>",
  "source": "<citation given in the field, or empty>",
  "type": "<person|date|event|quote|motif|place|han-tu|atu|source-existence|source-url>",
  "block": "frontmatter",
  "field": "<yaml field name, e.g. name_han, sources, birth_death, events, motifs>",
  "risk": "<high|medium|low>"
}

Extract from EVERY field that can be verified:
- name_han: Hán tự characters -> type=han-tu, risk=high
- aliases: alternate names -> type=person or place
- era: era claim -> type=date, risk=medium
- sources[]: EACH source is THREE claims: (1) "Source exists: [title] by [author]" type=source-existence risk=high; (2) chapter/edition/page if present type=event risk=high; (3) if url present, "Source URL: [url] for [title]" type=source-url risk=high — if url missing, emit "Source URL missing for [title]" type=source-url risk=high
- characters[]: each named character -> type=person
- motifs[]: each ATU/Thompson code -> type=atu, risk=high
- birth_death: -> type=date, risk=high
- dynasty: -> type=event, risk=medium
- events[]: each event name+year+role -> type=event, risk=high
- temples[]: each temple name+location -> type=place, risk=medium
- sovereignty_sources[]: title+author+year+relevance -> type=source-existence, risk=high
- coordinates: -> type=place, risk=medium
- summary: any factual claim in the summary -> risk=medium`

const chuyenKeExtractSystem = `You check narrative faithfulness in the "Chuyện kể" section of a Vietnamese mythology encyclopedia entry.

Output ONLY a JSON array. No preamble, no markdown fences.

Each item:
{
  "id": <number>,
  "text": "<the specific narrative claim or event>",
  "source": "",
  "type": "event",
  "block": "chuyen-ke",
  "field": "narrative",
  "risk": "<high|medium|low>"
}

Extract ONLY:
1. Specific plot events
2. Named characters and their actions
3. Specific locations mentioned in the narrative
4. Any specific numbers
5. Any dialogue or quoted speech
6. Any objects or artifacts with specific properties
7. Sensory details that go beyond what a source typically implies

Mark risk=high for: invented-feeling specific details, exact numbers, named dialogue, precise locations not in the title.
Mark risk=medium for: character actions central to the plot.
Mark risk=low for: general narrative flow that matches the known tale.`

const sectionExtractSystem = `You extract verifiable claims from an analysis section of a Vietnamese mythology encyclopedia entry.

Section name: %s
Entry type: %s

Output ONLY a JSON array. No preamble, no markdown fences.

Each claim:
{
  "id": <number>,
  "text": "<self-contained verifiable claim>",
  "source": "<attribution given in this section, or empty>",
  "type": "<person|date|event|quote|motif|place|han-tu|atu|source-existence|source-url>",
  "block": "section:%s",
  "field": "%s",
  "risk": "<high|medium|low>"
}

Rules by section type:
- named scholar + claim -> risk=high
- source references (LNCQ tale, ĐVSKTT quyển, etc.) -> risk=high
- inline markdown links [text](url) in analysis -> note url for verification
- dates, battles, dynasties, title names with years -> risk=high
- Hán tự -> risk=high
- ATU/Thompson codes -> risk=high
- ancient place = modern place -> risk=medium
- temple names, locations, festival dates -> risk=medium

Max 40 claims per section. Skip pure interpretive prose with no verifiable anchor.`

func extractFromBlock(block EntryBlock, entryType string, baseID int, llm LLM) ([]Claim, error) {
	var systemPrompt string
	var userPrompt string

	switch block.Kind {
	case "frontmatter":
		systemPrompt = frontmatterExtractSystem
		userPrompt = fmt.Sprintf("Extract all verifiable claims from this YAML frontmatter:\n\n---\n%s\n---", block.Content)
	case "chuyen-ke":
		systemPrompt = chuyenKeExtractSystem
		userPrompt = fmt.Sprintf("Check this narrative section for verifiable/suspicious claims:\n\n%s", block.Content)
	case "section":
		sectionName := block.Section
		systemPrompt = fmt.Sprintf(sectionExtractSystem, sectionName, entryType, sectionName, sectionName)
		userPrompt = fmt.Sprintf("Extract verifiable claims from this section:\n\nSection: %s\n\n%s", sectionName, block.Content)
	default:
		return nil, nil
	}

	raw, err := callLLMJSON(llm, systemPrompt, userPrompt, 8192)
	if err != nil {
		return nil, fmt.Errorf("extract [%s/%s]: %w", block.Kind, block.Section, err)
	}
	if raw == "[]" || raw == "" {
		return nil, nil
	}

	var claims []Claim
	if err := json.Unmarshal([]byte(raw), &claims); err != nil {
		return nil, fmt.Errorf("parse claims [%s/%s]: %w\nraw: %s", block.Kind, block.Section, err, raw)
	}

	for i := range claims {
		claims[i].ID = baseID + i
		if claims[i].Block == "" {
			claims[i].Block = block.Kind
		}
	}

	return claims, nil
}

func isJSONArray(raw string) bool {
	var claims []Claim
	return json.Unmarshal([]byte(raw), &claims) == nil
}

func isJSONObject(raw string) bool {
	var obj map[string]any
	return json.Unmarshal([]byte(raw), &obj) == nil
}
