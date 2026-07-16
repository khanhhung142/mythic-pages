package main

import (
	"fmt"
	"strings"
	"time"
)

func renderReport(r AuditReport) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Audit Report: %s\n\n", r.EntryTitle))
	sb.WriteString(fmt.Sprintf("**File:** `%s`  \n", r.EntryPath))
	sb.WriteString(fmt.Sprintf("**Date:** %s  \n", time.Now().Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Verdict:** %s\n\n", verdictBadge(r.Verdict)))
	sb.WriteString("---\n\n")

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString(r.Summary + "\n\n")

	// Stats
	counts := countStatuses(r.Results)
	sb.WriteString("| Status | Count |\n|---|---|\n")
	sb.WriteString(fmt.Sprintf("| ✅ verified | %d |\n", counts["verified"]))
	sb.WriteString(fmt.Sprintf("| ❌ wrong | %d |\n", counts["wrong"]))
	sb.WriteString(fmt.Sprintf("| ⚠️ suspicious | %d |\n", counts["suspicious"]))
	sb.WriteString(fmt.Sprintf("| 🔍 not_found | %d |\n", counts["not_found"]))
	sb.WriteString(fmt.Sprintf("| 🔶 pattern issues | %d |\n", len(r.Patterns)))
	sb.WriteString("\n---\n\n")

	// Claims needing attention first
	urgent := filterByStatus(r.Results, "wrong")
	suspicious := filterByStatus(r.Results, "suspicious")

	if len(urgent) > 0 {
		sb.WriteString("## ❌ Wrong Claims — Fix These\n\n")
		for _, res := range urgent {
			sb.WriteString(renderResult(res))
		}
	}

	if len(suspicious) > 0 {
		sb.WriteString("## ⚠️ Suspicious — Needs Human Review\n\n")
		for _, res := range suspicious {
			sb.WriteString(renderResult(res))
		}
	}

	// All results
	sb.WriteString("## All Claims\n\n")
	for _, res := range r.Results {
		sb.WriteString(renderResult(res))
	}

	// Pattern issues
	sb.WriteString("---\n\n## Writing Pattern Issues\n\n")
	if len(r.Patterns) == 0 {
		sb.WriteString("None.\n\n")
	} else {
		for _, p := range r.Patterns {
			sb.WriteString(fmt.Sprintf("**Line %d** [Pass %d] `%s`  \n", p.Line, p.Pass, p.Pattern))
			sb.WriteString(fmt.Sprintf("Text: _%s_  \n", p.Text))
			sb.WriteString(fmt.Sprintf("Fix: %s\n\n", p.Fix))
		}
	}

	return sb.String()
}

func renderResult(res VerificationResult) string {
	icon := statusIcon(res.Status)
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("### %s [%d] %s\n\n", icon, res.Claim.ID, res.Claim.Type))
	sb.WriteString(fmt.Sprintf("**Claim:** %s  \n", res.Claim.Text))
	if res.Claim.Source != "" {
		sb.WriteString(fmt.Sprintf("**Entry cite:** %s  \n", res.Claim.Source))
	}
	sb.WriteString(fmt.Sprintf("**Status:** %s (%s confidence)  \n", res.Status, res.Confidence))
	sb.WriteString(fmt.Sprintf("**Evidence:** %s  \n", res.Evidence))
	if res.SourceURL != "" {
		sb.WriteString(fmt.Sprintf("**Source:** %s  \n", res.SourceURL))
	}
	sb.WriteString("\n")

	return sb.String()
}

func verdictBadge(v string) string {
	switch v {
	case "PASS":
		return "✅ PASS"
	case "REVISE":
		return "⚠️ REVISE"
	case "REJECT":
		return "❌ REJECT"
	default:
		return v
	}
}

func statusIcon(s string) string {
	switch s {
	case "verified":
		return "✅"
	case "wrong":
		return "❌"
	case "suspicious":
		return "⚠️"
	default:
		return "🔍"
	}
}

func countStatuses(results []VerificationResult) map[string]int {
	counts := map[string]int{"verified": 0, "wrong": 0, "suspicious": 0, "not_found": 0}
	for _, r := range results {
		counts[r.Status]++
	}
	return counts
}

func filterByStatus(results []VerificationResult, status string) []VerificationResult {
	var out []VerificationResult
	for _, r := range results {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}

func deriveVerdict(results []VerificationResult, patterns []PatternIssue) (string, string) {
	counts := countStatuses(results)
	total := len(results)
	if total == 0 {
		return "REVISE", "No claims extracted — entry may be empty or unparseable."
	}

	wrongPct := float64(counts["wrong"]) / float64(total) * 100
	suspPct := float64(counts["suspicious"]) / float64(total) * 100
	stanceIssues := 0
	for _, p := range patterns {
		if p.Pass == 5 {
			stanceIssues++
		}
	}

	verdict := "PASS"
	if counts["wrong"] >= 3 || wrongPct >= 15 {
		verdict = "REJECT"
	} else if counts["wrong"] >= 1 || suspPct >= 20 || stanceIssues > 0 || len(patterns) >= 5 {
		verdict = "REVISE"
	}

	summary := fmt.Sprintf(
		"%d claims checked. %d verified, %d wrong (%.0f%%), %d suspicious (%.0f%%), %d not found. %d pattern issues (%d stance).",
		total,
		counts["verified"],
		counts["wrong"], wrongPct,
		counts["suspicious"], suspPct,
		counts["not_found"],
		len(patterns), stanceIssues,
	)

	return verdict, summary
}
