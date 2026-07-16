package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type pattern struct {
	re   *regexp.Regexp
	name string
	fix  string
	pass int
}

var bannedPatterns = []pattern{
	// Pass 6 — EN banned words
	{re: regexp.MustCompile(`(?i)\bdelve\b`), name: "AI word: delve", fix: "rephrase", pass: 6},
	{re: regexp.MustCompile(`(?i)\btapestry\b`), name: "AI word: tapestry", fix: "rephrase", pass: 6},
	{re: regexp.MustCompile(`(?i)\binterplay\b`), name: "AI word: interplay", fix: "state the relationship directly", pass: 6},
	{re: regexp.MustCompile(`(?i)\bintricate\b`), name: "AI word: intricate", fix: "rephrase or cut", pass: 6},
	{re: regexp.MustCompile(`(?i)\bpivotal\b`), name: "AI word: pivotal", fix: "state what changed", pass: 6},
	{re: regexp.MustCompile(`(?i)\bshowcase[sd]?\b`), name: "AI word: showcase", fix: "use 'shows' or 'demonstrates'", pass: 6},
	{re: regexp.MustCompile(`(?i)\btestament\b`), name: "AI word: testament", fix: "rephrase", pass: 6},
	{re: regexp.MustCompile(`(?i)\bunderscore[sd]?\b`), name: "AI word: underscore", fix: "state the point directly", pass: 6},
	{re: regexp.MustCompile(`(?i)\bvibrant\b`), name: "AI word: vibrant", fix: "cut or be specific", pass: 6},
	{re: regexp.MustCompile(`(?i)\bfoster[sd]?\b`), name: "AI word: foster", fix: "rephrase", pass: 6},
	{re: regexp.MustCompile(`(?i)\bencompass(es|ed)?\b`), name: "AI word: encompass", fix: "use 'includes' or list directly", pass: 6},
	{re: regexp.MustCompile(`(?i)\bembod(y|ies|ied)\b`), name: "AI word: embody", fix: "state what it represents directly", pass: 6},
	{re: regexp.MustCompile(`(?i)\bmultifaceted\b`), name: "AI word: multifaceted", fix: "list the facets", pass: 6},
	{re: regexp.MustCompile(`(?i)\bgroundbreaking\b`), name: "AI word: groundbreaking", fix: "state what was new specifically", pass: 6},
	{re: regexp.MustCompile(`(?i)\brenowned\b`), name: "AI word: renowned", fix: "use 'known for X' or cut", pass: 6},
	{re: regexp.MustCompile(`(?i)\bfascinating\b`), name: "AI word: fascinating", fix: "cut — let the fact speak", pass: 6},
	{re: regexp.MustCompile(`(?i)\bremarkable\b`), name: "AI word: remarkable", fix: "cut or be specific", pass: 6},

	// Pass 6 — VN banned phrases
	{re: regexp.MustCompile(`tỏa sáng`), name: "AI phrase VN: tỏa sáng", fix: "cụ thể hóa", pass: 6},
	{re: regexp.MustCompile(`nổi bật`), name: "AI phrase VN: nổi bật", fix: "nói cụ thể nổi bật cái gì", pass: 6},
	{re: regexp.MustCompile(`đặc sắc`), name: "AI phrase VN: đặc sắc (filler)", fix: "cụ thể hóa hoặc bỏ", pass: 6},
	{re: regexp.MustCompile(`góp phần quan trọng`), name: "AI phrase VN: góp phần quan trọng", fix: "nêu cụ thể đóng góp gì", pass: 6},
	{re: regexp.MustCompile(`mang ý nghĩa sâu sắc`), name: "AI phrase VN: mang ý nghĩa sâu sắc", fix: "nêu ý nghĩa cụ thể", pass: 6},
	{re: regexp.MustCompile(`thể hiện rõ nét`), name: "AI phrase VN: thể hiện rõ nét", fix: "nói thẳng thể hiện cái gì", pass: 6},
	{re: regexp.MustCompile(`di sản vô giá`), name: "AI ending VN: di sản vô giá", fix: "kết bằng fact cụ thể", pass: 6},
	{re: regexp.MustCompile(`sống mãi trong lòng dân tộc`), name: "AI ending VN: sống mãi...", fix: "kết bằng fact cụ thể", pass: 6},

	// Pass 5 — Stance issues
	{re: regexp.MustCompile(`(?i)south china sea`), name: "Stance: South China Sea", fix: "use 'East Sea' / 'Biển Đông'", pass: 5},
	{re: regexp.MustCompile(`(?i)disputed (waters|islands|territory)`), name: "Stance: neutralizing sovereignty", fix: "state VN position with sources", pass: 5},
	{re: regexp.MustCompile(`(?i)both sides claim`), name: "Stance: both-sides framing", fix: "present VN position as entry voice, note PRC claim separately", pass: 5},

	// Pass 6 — structural patterns
	{re: regexp.MustCompile(` — | —|— `), name: "Em dash in prose", fix: "restructure with comma, colon, or period", pass: 6},
	{re: regexp.MustCompile(`(?i)(let's explore|hãy cùng tìm hiểu)`), name: "Signposting opener", fix: "start directly", pass: 6},
	{re: regexp.MustCompile(`(?i)(scholars argue|experts believe|nhiều nhà nghiên cứu cho rằng)`), name: "Vague attribution", fix: "name the scholar + work + year", pass: 6},
	{re: regexp.MustCompile(`(?i)not (just|only) .{1,60}(;|,) (it'?s?|but also)`), name: "Negative parallelism: not just X but Y", fix: "state what it IS", pass: 6},
	{re: regexp.MustCompile(`(?i)(serves as|stands as|functions as|marks a) (a |an )?`), name: "Copula avoidance", fix: "use 'is'", pass: 6},
}

func scanPatterns(content string) []PatternIssue {
	var issues []PatternIssue

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip YAML frontmatter lines and code blocks
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			continue
		}

		for _, p := range bannedPatterns {
			if p.re.MatchString(line) {
				snippet := line
				if len(snippet) > 100 {
					// Find match position and show context
					loc := p.re.FindStringIndex(line)
					start := loc[0] - 30
					if start < 0 {
						start = 0
					}
					end := loc[1] + 30
					if end > len(line) {
						end = len(line)
					}
					snippet = "..." + line[start:end] + "..."
				}
				issues = append(issues, PatternIssue{
					Pass:    p.pass,
					Line:    lineNum,
					Pattern: p.name,
					Text:    strings.TrimSpace(snippet),
					Fix:     p.fix,
				})
			}
		}
	}

	return issues
}

func formatPatternIssues(issues []PatternIssue) string {
	if len(issues) == 0 {
		return "No pattern issues found."
	}

	var sb strings.Builder
	for _, issue := range issues {
		sb.WriteString(fmt.Sprintf("Line %d [Pass %d] %s\n  Text: %s\n  Fix: %s\n\n",
			issue.Line, issue.Pass, issue.Pattern, issue.Text, issue.Fix))
	}
	return sb.String()
}
