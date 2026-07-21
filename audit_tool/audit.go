package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func runAudit(entryPath, outputPath string, verbose bool, rt Runtime) error {
	entryPath, err := resolveEntryPath(entryPath)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(entryPath)
	if err != nil {
		return fmt.Errorf("read entry: %w", err)
	}
	raw := string(content)

	entryTitle := extractTitle(raw, entryPath)
	entryType := detectEntryType(raw)
	logf(verbose, "Entry: %s (type: %s)\n", entryTitle, entryType)
	logf(verbose, "AI: extract=%s judge=%s search=perplexity\n", rt.ExtractProvider, rt.JudgeProvider)

	logf(verbose, "Parsing blocks...\n")
	blocks := parseEntry(raw)
	logf(verbose, "  %d blocks found\n", len(blocks))
	for _, b := range blocks {
		logf(verbose, "    [%s] %s\n", b.Kind, truncate(b.Section, 40))
	}

	logf(verbose, "Extracting claims per block...\n")
	var blockAudits []BlockAudit
	claimID := 1
	for _, block := range blocks {
		logf(verbose, "  extracting [%s/%s]...\n", block.Kind, truncate(block.Section, 30))

		claims, err := extractFromBlock(block, entryType, claimID, rt.Extract)
		if err != nil {
			logf(verbose, "  WARN: %v\n", err)
			blockAudits = append(blockAudits, BlockAudit{Block: block})
			continue
		}
		claimID += len(claims)
		logf(verbose, "    %d claims\n", len(claims))
		blockAudits = append(blockAudits, BlockAudit{
			Block:  block,
			Claims: claims,
		})

		time.Sleep(300 * time.Millisecond)
	}

	logf(verbose, "Verifying claims...\n")
	totalClaims := 0
	for _, ba := range blockAudits {
		totalClaims += len(ba.Claims)
	}
	logf(verbose, "  total claims to verify: %d\n", totalClaims)
	blockAudits = verifyAllBlocks(blockAudits, verbose, rt)

	logf(verbose, "Scanning patterns...\n")
	patterns := scanPatterns(raw)
	logf(verbose, "  %d pattern issues\n", len(patterns))

	logf(verbose, "Auditing source URLs...\n")
	fm := extractFrontmatter(raw)
	sources := parseSources(fm)
	sourceLinks := auditSourceLinks(sources, nil)
	sourceLinks = append(sourceLinks, auditInlineLinks(raw)...)
	unlinked := scanUnlinkedCitations(blocks)
	patterns = append(patterns, unlinked...)
	missingURL, badURL, unreachable := countSourceLinkIssues(sourceLinks)
	logf(verbose, "  sources=%d missing_url=%d bad=%d unreachable=%d unlinked_cites=%d\n",
		len(sources), missingURL, badURL, unreachable, len(unlinked))

	allResults := flattenResults(blockAudits)
	verdict, rejectReason, summary := deriveVerdict(allResults, patterns, blockAudits, sourceLinks)

	report := AuditReport{
		EntryPath:    entryPath,
		EntryTitle:   entryTitle,
		EntryType:    entryType,
		BlockAudits:  blockAudits,
		SourceLinks:  sourceLinks,
		Patterns:     patterns,
		Verdict:      verdict,
		RejectReason: rejectReason,
		Summary:      summary,
	}

	outputPath, err = resolveAuditOutput(entryPath, outputPath)
	if err != nil {
		return err
	}

	rendered := renderReport(report)
	if err := os.WriteFile(outputPath, []byte(rendered), 0644); err != nil {
		return fmt.Errorf("write report: %w", err)
	}

	fmt.Printf("\nVerdict: %s\n", verdict)
	if rejectReason != "" {
		fmt.Printf("Reason:  %s\n", rejectReason)
	}
	fmt.Printf("Summary: %s\n", summary)
	fmt.Printf("Report:  %s\n", outputPath)

	return nil
}

func resolveEntryPath(arg string) (string, error) {
	if _, err := os.Stat(arg); err == nil {
		return arg, nil
	}

	if filepath.Ext(arg) != "" || strings.ContainsRune(arg, os.PathSeparator) {
		return arg, nil
	}

	// ponytail: support the two common launch dirs now; switch to executable-relative
	// lookup later if the tool needs to be run from arbitrary directories.
	candidates := []string{
		filepath.Join("src", "content", "vi", "entries", arg+".md"),
		filepath.Join("..", "src", "content", "vi", "entries", arg+".md"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return arg, nil
}

// resolveAuditOutput always writes under audit/<slug>-audit.md (or -o basename).
func resolveAuditOutput(entryPath, custom string) (string, error) {
	slug := strings.TrimSuffix(filepath.Base(entryPath), filepath.Ext(entryPath))
	name := slug + "-audit.md"
	if custom != "" {
		name = filepath.Base(custom)
	}

	dir := resolveAuditDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create audit dir: %w", err)
	}
	return filepath.Join(dir, name), nil
}

func resolveAuditDir() string {
	for _, d := range []string{"audit", filepath.Join("..", "audit")} {
		if st, err := os.Stat(d); err == nil && st.IsDir() {
			return d
		}
	}
	if _, err := os.Stat(filepath.Join("..", "src", "content")); err == nil {
		return filepath.Join("..", "audit")
	}
	return "audit"
}

func verifyAllBlocks(blockAudits []BlockAudit, verbose bool, rt Runtime) []BlockAudit {
	type indexedClaim struct {
		blockIdx int
		claimIdx int
		claim    Claim
	}

	var all []indexedClaim
	for bi, ba := range blockAudits {
		for ci, c := range ba.Claims {
			all = append(all, indexedClaim{bi, ci, c})
		}
	}

	results := make([]VerificationResult, len(all))
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for i, ic := range all {
		wg.Add(1)
		go func(idx int, ic indexedClaim) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			logf(verbose, "  [%d/%d] [%s] %s\n",
				idx+1, len(all),
				ic.claim.Block,
				truncate(ic.claim.Text, 55))

			res, err := verifyClaim(ic.claim, verbose, rt)
			if err != nil {
				res = VerificationResult{
					Claim:      ic.claim,
					Status:     "not_found",
					Evidence:   fmt.Sprintf("Error: %v", err),
					Confidence: "low",
				}
			}
			results[idx] = res

			time.Sleep(200 * time.Millisecond)
		}(i, ic)
	}

	wg.Wait()

	for i, ic := range all {
		blockAudits[ic.blockIdx].Results = ensureSize(blockAudits[ic.blockIdx].Results, ic.claimIdx+1)
		blockAudits[ic.blockIdx].Results[ic.claimIdx] = results[i]
	}

	return blockAudits
}

func ensureSize(s []VerificationResult, n int) []VerificationResult {
	for len(s) < n {
		s = append(s, VerificationResult{})
	}
	return s
}

func flattenResults(blockAudits []BlockAudit) []VerificationResult {
	var out []VerificationResult
	for _, ba := range blockAudits {
		out = append(out, ba.Results...)
	}
	return out
}

func deriveVerdict(results []VerificationResult, patterns []PatternIssue, blockAudits []BlockAudit, sourceLinks []SourceLinkResult) (verdict, rejectReason, summary string) {
	total := len(results)
	if total == 0 {
		return "REVISE", "No claims extracted", "Entry may be empty or unparseable."
	}

	counts := map[string]int{}
	highRiskWrong := 0
	inventedCount := 0
	chuyenKeInvented := 0

	for _, r := range results {
		counts[r.Status]++
		if r.Status == "invented" {
			inventedCount++
			if r.Claim.Block == "chuyen-ke" {
				chuyenKeInvented++
			}
		}
		if (r.Status == "wrong" || r.Status == "invented") && r.Claim.Risk == "high" {
			highRiskWrong++
		}
	}

	wrongPct := float64(counts["wrong"]+counts["invented"]) / float64(total) * 100
	suspPct := float64(counts["suspicious"]) / float64(total) * 100

	stanceIssues := 0
	unlinkedCites := 0
	for _, p := range patterns {
		if p.Pass == 5 {
			stanceIssues++
		}
		if p.Pass == 7 && p.Pattern == "Unlinked scholarly citation" {
			unlinkedCites++
		}
	}
	missingURL, badURL, unreachable := countSourceLinkIssues(sourceLinks)

	switch {
	case chuyenKeInvented >= 1:
		verdict = "REJECT"
		rejectReason = fmt.Sprintf("Invented plot event(s) in Chuyện kể (%d detected)", chuyenKeInvented)
	case inventedCount >= 2:
		verdict = "REJECT"
		rejectReason = fmt.Sprintf("Multiple invented claims (%d): fabricated sources or ATU codes", inventedCount)
	case highRiskWrong >= 3:
		verdict = "REJECT"
		rejectReason = fmt.Sprintf("%d high-risk claims wrong (dates, sources, Hán tự)", highRiskWrong)
	case wrongPct >= 15:
		verdict = "REJECT"
		rejectReason = fmt.Sprintf("%.0f%% of claims wrong or invented", wrongPct)
	case counts["wrong"]+counts["invented"] >= 1 || suspPct >= 20 || stanceIssues > 0 || len(patterns) >= 5 || missingURL > 0 || badURL > 0 || unreachable > 0 || unlinkedCites > 0:
		verdict = "REVISE"
		if missingURL > 0 && rejectReason == "" {
			rejectReason = fmt.Sprintf("%d source(s) missing url in sources[]", missingURL)
		}
		if badURL > 0 && rejectReason == "" {
			rejectReason = fmt.Sprintf("%d source url(s) on banned/invalid domain", badURL)
		}
		if unreachable > 0 && rejectReason == "" {
			rejectReason = fmt.Sprintf("%d source url(s) unreachable (4xx/5xx or network error)", unreachable)
		}
		if unlinkedCites > 0 && rejectReason == "" {
			rejectReason = fmt.Sprintf("%d inline citation(s) without markdown link", unlinkedCites)
		}
	default:
		verdict = "PASS"
	}

	summary = fmt.Sprintf(
		"%d claims across %d blocks. verified=%d wrong=%d invented=%d suspicious=%d not_found=%d (wrong+invented=%.0f%%). Patterns=%d (stance=%d). Sources=%d (missing_url=%d bad=%d unreachable=%d unlinked_cites=%d).",
		total, len(blockAudits),
		counts["verified"], counts["wrong"], counts["invented"],
		counts["suspicious"], counts["not_found"],
		wrongPct, len(patterns), stanceIssues,
		len(sourceLinks), missingURL, badURL, unreachable, unlinkedCites,
	)

	return verdict, rejectReason, summary
}

func extractTitle(content, path string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name_vi:") {
			return strings.Trim(strings.TrimPrefix(line, "name_vi:"), ` "`)
		}
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return filepath.Base(path)
}

func logf(verbose bool, format string, args ...any) {
	if verbose {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
