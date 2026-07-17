package main

import "strings"

func parseEntry(content string) []EntryBlock {
	var blocks []EntryBlock
	lines := strings.Split(content, "\n")

	const (
		stateStart = iota
		stateFrontmatter
		stateBody
	)

	state := stateStart
	var fmLines []string
	var currentSection string
	var currentLines []string
	inFrontmatter := false

	flush := func(kind, section string, ls []string) {
		text := strings.TrimSpace(strings.Join(ls, "\n"))
		if text == "" {
			return
		}
		blocks = append(blocks, EntryBlock{
			Kind:    kind,
			Section: section,
			Content: text,
		})
	}

	for i, line := range lines {
		switch state {
		case stateStart:
			if strings.TrimSpace(line) == "---" {
				inFrontmatter = true
				state = stateFrontmatter
				fmLines = []string{}
			}
		case stateFrontmatter:
			if strings.TrimSpace(line) == "---" && i > 0 && inFrontmatter {
				flush("frontmatter", "", fmLines)
				inFrontmatter = false
				state = stateBody
				currentSection = ""
				currentLines = []string{}
				continue
			}
			fmLines = append(fmLines, line)
		case stateBody:
			trimmed := strings.TrimSpace(line)
			if isChuyenKeHeading(trimmed) {
				if len(currentLines) > 0 {
					kind, sec := classifySection(currentSection)
					flush(kind, sec, currentLines)
				}
				currentSection = "chuyen-ke"
				currentLines = []string{}
				continue
			}
			if strings.HasPrefix(trimmed, "## ") {
				if len(currentLines) > 0 {
					kind, sec := classifySection(currentSection)
					flush(kind, sec, currentLines)
				}
				currentSection = strings.TrimPrefix(trimmed, "## ")
				currentLines = []string{}
				continue
			}
			currentLines = append(currentLines, line)
		}
	}

	if len(currentLines) > 0 && state == stateBody {
		kind, sec := classifySection(currentSection)
		flush(kind, sec, currentLines)
	}

	return blocks
}

func isChuyenKeHeading(s string) bool {
	lower := strings.ToLower(s)
	return strings.Contains(lower, "chuyện kể") ||
		strings.Contains(lower, "chuyen ke") ||
		strings.Contains(lower, "the story") ||
		s == "## Chuyện kể" ||
		s == "## The Story"
}

func classifySection(name string) (kind, section string) {
	if name == "" {
		return "unknown", ""
	}
	if name == "chuyen-ke" {
		return "chuyen-ke", "Chuyện kể"
	}
	return "section", name
}

func detectEntryType(content string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "category:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "category:"))
			return strings.Trim(val, `"' `)
		}
	}
	return "unknown"
}
