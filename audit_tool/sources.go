package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// SourceRef is one item from frontmatter sources[].
type SourceRef struct {
	Title  string
	Author string
	URL    string
}

// SourceLinkResult is the audit outcome for one source URL.
type SourceLinkResult struct {
	Source   SourceRef
	Status   string // ok | missing_url | bad_domain | unreachable | invalid_url
	Evidence string
}

var bannedSourceHosts = []string{
	"medium.com",
	"substack.com",
	"wordpress.com",
	"blogspot.com",
	"tumblr.com",
	"facebook.com",
	"zalo.me",
	"wikipedia.org",
	"wikisource.org",
	"wikimedia.org",
	"dotchuoinon.com",
	"blog.",
	"tia-sang.com",
}

var trustedSourceHints = []string{
	"jstor.org",
	"doi.org",
	"vjol.info.vn",
	"han-nom.org",
	"thuvienquocgia.gov.vn",
	"gallica.bnf.fr",
	"worldcat.org",
	"archive.org",
	"nxb",
	".edu",
	".ac.vn",
	".gov.vn",
	"befeo",
	"persee.fr",
	"cnrs.fr",
	"unesco.org",
}

var (
	inlineLinkRe     = regexp.MustCompile(`\[[^\]]+\]\((https?://[^)]+)\)`)
	unlinkedYearRe   = regexp.MustCompile(`\((19|20)\d{2}[^)]*\)|,\s*(19|20)\d{2}\b|số\s+\d+/\d{4}|tr\.\s*\d+`)
	italicTitleRe    = regexp.MustCompile(`\*[^*]{3,}\*`)
	sourcesItemStart = regexp.MustCompile(`^\s*-\s+title:\s*`)
)

func extractFrontmatter(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return ""
	}
	var fm []string
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			return strings.Join(fm, "\n")
		}
		fm = append(fm, lines[i])
	}
	return ""
}

func parseSources(frontmatter string) []SourceRef {
	var out []SourceRef
	var cur *SourceRef
	inSources := false

	for _, line := range strings.Split(frontmatter, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "sources:" {
			inSources = true
			continue
		}
		if !inSources {
			continue
		}
		if trimmed != "" && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") && !sourcesItemStart.MatchString(line) {
			break
		}
		if sourcesItemStart.MatchString(line) {
			if cur != nil && (cur.Title != "" || cur.Author != "" || cur.URL != "") {
				out = append(out, *cur)
			}
			cur = &SourceRef{Title: unquoteYAML(strings.TrimPrefix(trimmed, "- title:"))}
			continue
		}
		if cur == nil {
			continue
		}
		switch {
		case strings.HasPrefix(trimmed, "author:"):
			cur.Author = unquoteYAML(strings.TrimPrefix(trimmed, "author:"))
		case strings.HasPrefix(trimmed, "url:"):
			cur.URL = strings.TrimSpace(unquoteYAML(strings.TrimPrefix(trimmed, "url:")))
		}
	}
	if cur != nil && (cur.Title != "" || cur.Author != "" || cur.URL != "") {
		out = append(out, *cur)
	}
	return out
}

func unquoteYAML(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `"'`)
	return s
}

func auditSourceLinks(sources []SourceRef, client *http.Client) []SourceLinkResult {
	if client == nil {
		client = &http.Client{Timeout: 12 * time.Second}
	}
	out := make([]SourceLinkResult, 0, len(sources))
	for _, src := range sources {
		out = append(out, auditOneSourceURL(src, client))
	}
	return out
}

func auditOneSourceURL(src SourceRef, client *http.Client) SourceLinkResult {
	res := SourceLinkResult{Source: src}
	if strings.TrimSpace(src.URL) == "" {
		res.Status = "missing_url"
		res.Evidence = "sources[] item has no url field (skill §2.7)"
		return res
	}
	u, err := url.Parse(src.URL)
	if err != nil || u.Scheme != "http" && u.Scheme != "https" || u.Host == "" {
		res.Status = "invalid_url"
		res.Evidence = fmt.Sprintf("malformed url: %q", src.URL)
		return res
	}
	host := strings.ToLower(u.Host)
	if isBannedSourceHost(host) {
		res.Status = "bad_domain"
		res.Evidence = fmt.Sprintf("host %q is not an allowed institutional/scholarly source", host)
		return res
	}
	code, fetchErr := probeURL(client, src.URL)
	if fetchErr != nil {
		res.Status = "unreachable"
		res.Evidence = fmt.Sprintf("request failed: %v", fetchErr)
		return res
	}
	// Bot-blocking (403/401/429) means the host refuses automated probes, NOT
	// that the page is dead. Flagging these as unreachable produced false
	// REVISE verdicts on gov/scholarly hosts. Surface for human eyes instead.
	if code == 401 || code == 403 || code == 429 {
		res.Status = "blocked"
		res.Evidence = fmt.Sprintf("HTTP %d — host blocks automated probes; verify link by hand", code)
		return res
	}
	if code >= 400 {
		res.Status = "unreachable"
		res.Evidence = fmt.Sprintf("HTTP %d", code)
		return res
	}
	res.Status = "ok"
	if isTrustedSourceHost(host) {
		res.Evidence = fmt.Sprintf("HTTP %d, trusted host", code)
	} else {
		res.Evidence = fmt.Sprintf("HTTP %d, host not on trusted list — human review recommended", code)
	}
	return res
}

func auditInlineLinks(body string) []SourceLinkResult {
	client := &http.Client{Timeout: 12 * time.Second}
	seen := map[string]bool{}
	var out []SourceLinkResult
	for _, m := range inlineLinkRe.FindAllStringSubmatch(body, -1) {
		link := m[1]
		if seen[link] {
			continue
		}
		seen[link] = true
		src := SourceRef{Title: "inline citation", URL: link}
		out = append(out, auditOneSourceURL(src, client))
	}
	return out
}

func scanUnlinkedCitations(blocks []EntryBlock) []PatternIssue {
	var issues []PatternIssue
	for _, b := range blocks {
		if b.Kind == "chuyen-ke" || b.Kind == "frontmatter" {
			continue
		}
		lineNo := 0
		for _, line := range strings.Split(b.Content, "\n") {
			lineNo++
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				continue
			}
			if strings.Contains(line, "](http") {
				continue
			}
			hasYear := unlinkedYearRe.MatchString(line)
			hasTitle := italicTitleRe.MatchString(line)
			if !hasYear && !hasTitle {
				continue
			}
			section := b.Section
			if section == "" {
				section = b.Kind
			}
			issues = append(issues, PatternIssue{
				Pass:    7,
				Line:    lineNo,
				Pattern: "Unlinked scholarly citation",
				Text:    truncate(trimmed, 120),
				Fix:     fmt.Sprintf("wrap citation in [%s](url) pointing to a reputable source (skill §2.7)", truncate(trimmed, 40)),
			})
		}
	}
	return issues
}

func isBannedSourceHost(host string) bool {
	for _, bad := range bannedSourceHosts {
		if strings.Contains(host, bad) {
			return true
		}
	}
	return false
}

func isTrustedSourceHost(host string) bool {
	for _, ok := range trustedSourceHints {
		if strings.Contains(host, ok) {
			return true
		}
	}
	return false
}

func probeURL(client *http.Client, rawURL string) (int, error) {
	code, err := doProbe(client, http.MethodHead, rawURL)
	if err == nil && code != http.StatusMethodNotAllowed && code != http.StatusNotFound {
		return code, nil
	}
	return doProbe(client, http.MethodGet, rawURL)
}

func doProbe(client *http.Client, method, rawURL string) (int, error) {
	req, err := http.NewRequest(method, rawURL, nil)
	if err != nil {
		return 0, err
	}
	// Real browser UA: many gov/scholarly hosts 403 a bot-looking UA, which
	// otherwise shows up as a false dead-link.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func countSourceLinkIssues(results []SourceLinkResult) (missing, bad, unreachable int) {
	for _, r := range results {
		switch r.Status {
		case "missing_url":
			missing++
		case "bad_domain", "invalid_url":
			bad++
		case "unreachable":
			unreachable++
		}
	}
	return
}
