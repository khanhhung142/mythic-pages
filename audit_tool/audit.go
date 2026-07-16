package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func runAudit(entryPath, outputPath string, verbose bool) error {
	// Read entry
	content, err := os.ReadFile(entryPath)
	if err != nil {
		return fmt.Errorf("read entry: %w", err)
	}

	entryTitle := extractTitle(string(content), entryPath)
	logf(verbose, "Auditing: %s\n", entryTitle)

	// Pass 1 — extract claims
	logf(verbose, "Pass 1: extracting claims...\n")
	claims, err := extractClaims(string(content))
	if err != nil {
		return fmt.Errorf("extract claims: %w", err)
	}
	logf(verbose, "  %d claims found\n", len(claims))

	// Pass 3 — verify claims (concurrent, max 5 at a time)
	logf(verbose, "Pass 3: verifying claims...\n")
	results := verifyClaims(claims, verbose)

	// Pass 5+6 — pattern scan
	logf(verbose, "Pass 5+6: scanning patterns...\n")
	patterns := scanPatterns(string(content))
	logf(verbose, "  %d pattern issues\n", len(patterns))

	// Derive verdict
	verdict, summary := deriveVerdict(results, patterns)

	report := AuditReport{
		EntryPath:  entryPath,
		EntryTitle: entryTitle,
		Claims:     claims,
		Results:    results,
		Patterns:   patterns,
		Verdict:    verdict,
		Summary:    summary,
	}

	// Output path
	if outputPath == "" {
		base := strings.TrimSuffix(entryPath, filepath.Ext(entryPath))
		outputPath = base + "-audit.md"
	}

	rendered := renderReport(report)
	if err := os.WriteFile(outputPath, []byte(rendered), 0644); err != nil {
		return fmt.Errorf("write report: %w", err)
	}

	fmt.Printf("\nVerdict: %s\n", verdict)
	fmt.Printf("Summary: %s\n", summary)
	fmt.Printf("Report:  %s\n", outputPath)

	return nil
}

func verifyClaims(claims []Claim, verbose bool) []VerificationResult {
	results := make([]VerificationResult, len(claims))
	sem := make(chan struct{}, 5) // max 5 concurrent
	var wg sync.WaitGroup

	for i, claim := range claims {
		wg.Add(1)
		go func(idx int, c Claim) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			logf(verbose, "  [%d/%d] %s\n", idx+1, len(claims), truncate(c.Text, 60))

			res, err := verifyClaim(c, verbose)
			if err != nil {
				res = VerificationResult{
					Claim:      c,
					Status:     "not_found",
					Evidence:   fmt.Sprintf("Error: %v", err),
					Confidence: "low",
				}
			}
			results[idx] = res

			// Rate limit: avoid hammering APIs
			time.Sleep(200 * time.Millisecond)
		}(i, claim)
	}

	wg.Wait()
	return results
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
