package main

import (
	"os"
	"testing"
)

// Live smoke test: the exact claim the auditor previously flagged as
// "invented" (the JSEAS 2013 article by Nguyễn Thị Điểu, which exists).
// Run: go test -run Smoke -v   (needs both API keys; skipped otherwise)
func TestSmokeVerifyRealSource(t *testing.T) {
	if os.Getenv("PERPLEXITY_API_KEY") == "" {
		t.Skip("PERPLEXITY_API_KEY not set")
	}
	rt, err := NewRuntime(AIConfig{LLMProvider: envOr("AUDIT_LLM", "claude")})
	if err != nil {
		t.Skip(err)
	}
	c := Claim{
		ID:     1,
		Text:   "Nguyễn Thị Điểu published a study tracking rewritings of the Âu Cơ legend in 2013",
		Source: "Nguyễn Thị Điểu, JSEAS 44/2, 2013, tr. 315-337",
		Type:   "source-existence",
		Block:  "section:scholarship",
		Field:  "scholarship",
		Risk:   "high",
	}
	r, err := verifyClaim(c, true, rt)
	if err != nil {
		t.Fatalf("verifyClaim: %v", err)
	}
	t.Logf("status=%s conf=%s evidence=%s url=%s", r.Status, r.Confidence, r.Evidence, r.SourceURL)
	if r.Status == "invented" || r.Status == "wrong" {
		t.Errorf("real source judged %s (was the exact false positive this test guards against)", r.Status)
	}
}
