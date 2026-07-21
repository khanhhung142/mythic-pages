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
	sb.WriteString(fmt.Sprintf("**Type:** %s  \n", r.EntryType))
	sb.WriteString(fmt.Sprintf("**Date:** %s  \n", time.Now().Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Verdict:** %s\n", verdictBadge(r.Verdict)))
	if r.RejectReason != "" {
		sb.WriteString(fmt.Sprintf("**Reject reason:** %s\n", r.RejectReason))
	}
	sb.WriteString("\n---\n\n")

	sb.WriteString("## Summary\n\n")
	sb.WriteString(r.Summary + "\n\n")

	allResults := flattenResults(r.BlockAudits)
	counts := countStatuses(allResults)
	sb.WriteString("| Status | Count |\n|---|---|\n")
	sb.WriteString(fmt.Sprintf("| ✅ verified | %d |\n", counts["verified"]))
	sb.WriteString(fmt.Sprintf("| ❌ wrong | %d |\n", counts["wrong"]))
	sb.WriteString(fmt.Sprintf("| 💀 invented | %d |\n", counts["invented"]))
	sb.WriteString(fmt.Sprintf("| ⚠️ suspicious | %d |\n", counts["suspicious"]))
	sb.WriteString(fmt.Sprintf("| 🔍 not_found | %d |\n", counts["not_found"]))
	sb.WriteString(fmt.Sprintf("| 🔶 pattern issues | %d |\n\n", len(r.Patterns)))

	urgent := filterByStatuses(allResults, "invented", "wrong")
	if len(urgent) > 0 {
		sb.WriteString("---\n\n## 🚨 Fix Immediately (wrong + invented)\n\n")
		for _, res := range urgent {
			sb.WriteString(renderResult(res))
		}
	}

	suspicious := filterByStatuses(allResults, "suspicious")
	if len(suspicious) > 0 {
		sb.WriteString("---\n\n## ⚠️ Suspicious — Human Review\n\n")
		for _, res := range suspicious {
			sb.WriteString(renderResult(res))
		}
	}

	sb.WriteString("---\n\n## Per-Block Results\n\n")
	for _, ba := range r.BlockAudits {
		if len(ba.Claims) == 0 {
			continue
		}

		blockLabel := ba.Block.Kind
		if ba.Block.Section != "" {
			blockLabel = fmt.Sprintf("%s / %s", ba.Block.Kind, ba.Block.Section)
		}

		bc := countStatuses(ba.Results)
		sb.WriteString(fmt.Sprintf("### %s\n\n", blockLabel))
		sb.WriteString(fmt.Sprintf("Claims: %d | ✅%d ❌%d 💀%d ⚠️%d 🔍%d\n\n",
			len(ba.Claims),
			bc["verified"], bc["wrong"], bc["invented"], bc["suspicious"], bc["not_found"]))

		for _, res := range ba.Results {
			if res.Claim.Text == "" {
				continue
			}
			sb.WriteString(renderResult(res))
		}
	}

	sb.WriteString("---\n\n## Source URL Audit\n\n")
	if len(r.SourceLinks) == 0 {
		sb.WriteString("No `sources[]` entries or inline citation links found.\n\n")
	} else {
		sb.WriteString("| # | Title | URL | Status | Note |\n|---|---|---|---|---|\n")
		for i, sl := range r.SourceLinks {
			title := sl.Source.Title
			if sl.Source.Author != "" && title != "inline citation" {
				title = sl.Source.Author + " — " + title
			}
			url := sl.Source.URL
			if url == "" {
				url = "—"
			}
			sb.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s |\n",
				i+1, title, url, sl.Status, sl.Evidence))
		}
		sb.WriteString("\n")
	}

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

	riskBadge := ""
	if res.Claim.Risk == "high" {
		riskBadge = " `HIGH`"
	}

	sb.WriteString(fmt.Sprintf("#### %s [%d]%s — %s\n\n", icon, res.Claim.ID, riskBadge, res.Claim.Type))
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
	case "invented":
		return "💀"
	case "suspicious":
		return "⚠️"
	default:
		return "🔍"
	}
}

func countStatuses(results []VerificationResult) map[string]int {
	counts := map[string]int{"verified": 0, "wrong": 0, "invented": 0, "suspicious": 0, "not_found": 0}
	for _, r := range results {
		counts[r.Status]++
	}
	return counts
}

func filterByStatuses(results []VerificationResult, statuses ...string) []VerificationResult {
	set := map[string]bool{}
	for _, status := range statuses {
		set[status] = true
	}

	var out []VerificationResult
	for _, r := range results {
		if set[r.Status] {
			out = append(out, r)
		}
	}
	return out
}
