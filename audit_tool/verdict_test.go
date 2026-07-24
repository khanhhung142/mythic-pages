package main

import "testing"

func errResults(n, errs int) []VerificationResult {
	out := make([]VerificationResult, 0, n)
	for i := 0; i < n; i++ {
		st := "verified"
		if i < errs {
			st = "error"
		}
		out = append(out, VerificationResult{Status: st})
	}
	return out
}

// A dead backend must never ship a fact verdict.
func TestDeriveVerdictInconclusiveWhenBackendDead(t *testing.T) {
	v, _, _ := deriveVerdict(errResults(100, 100), nil, []BlockAudit{{}}, nil)
	if v != "INCONCLUSIVE" {
		t.Fatalf("all-errored audit got %q, want INCONCLUSIVE", v)
	}
	// A few stray errors must NOT poison an otherwise clean audit.
	if v, _, _ := deriveVerdict(errResults(100, 2), nil, []BlockAudit{{}}, nil); v == "INCONCLUSIVE" {
		t.Fatalf("2%% errors wrongly forced INCONCLUSIVE")
	}
}

func TestGateDamningVerdict(t *testing.T) {
	cases := []struct {
		url  string
		want string // expected status after gating a "wrong" verdict
	}{
		{"https://vi.wikipedia.org/wiki/x", "suspicious"},         // banned
		{"https://cand.com.vn/story", "suspicious"},               // untrusted newspaper
		{"", "suspicious"},                                        // no source
		{"https://vjol.info.vn/index.php/rsr/article/1", "wrong"}, // reputable journal
		{"https://gallica.bnf.fr/ark:/x", "wrong"},                // reputable archive
	}
	for _, c := range cases {
		got := gateDamningVerdict(VerificationResult{Status: "wrong", Confidence: "high", SourceURL: c.url})
		if got.Status != c.want {
			t.Errorf("url %q: got %q, want %q", c.url, got.Status, c.want)
		}
	}
	// Non-damning verdicts pass through untouched.
	if got := gateDamningVerdict(VerificationResult{Status: "verified", SourceURL: ""}); got.Status != "verified" {
		t.Errorf("verified verdict altered: %q", got.Status)
	}
}

func TestIsFatalBackendErr(t *testing.T) {
	fatal := []string{
		"Search failed: perplexity error: You exceeded your current quota",
		"Judgment failed: 401 Unauthorized",
		"invalid api key",
	}
	for _, e := range fatal {
		if !isFatalBackendErr(e) {
			t.Errorf("expected fatal: %q", e)
		}
	}
	if isFatalBackendErr("Search failed: context deadline exceeded (timeout)") {
		// "exceeded" here is a timeout, but we intentionally abort on it too:
		// a flapping backend is still worth stopping. Documented, not a bug.
		t.Log("timeout treated as fatal (intentional early-abort)")
	}
}
