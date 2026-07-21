package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseSources(t *testing.T) {
	fm := `name_vi: Test
sources:
  - title: Lĩnh Nam chích quái
    author: Trần Thế Pháp
    edition: NXB Văn hoá, 1960
    url: https://example.org/lncq
  - title: Việt điện u linh
    author: Lý Tế Xuyên
summary: x`
	got := parseSources(fm)
	if len(got) != 2 {
		t.Fatalf("want 2 sources, got %d", len(got))
	}
	if got[0].URL != "https://example.org/lncq" {
		t.Fatalf("url=%q", got[0].URL)
	}
	if got[1].URL != "" {
		t.Fatalf("second url should be empty, got %q", got[1].URL)
	}
}

func TestAuditOneSourceURL_missing(t *testing.T) {
	r := auditOneSourceURL(SourceRef{Title: "X"}, &http.Client{})
	if r.Status != "missing_url" {
		t.Fatalf("status=%s", r.Status)
	}
}

func TestAuditOneSourceURL_bannedDomain(t *testing.T) {
	r := auditOneSourceURL(SourceRef{Title: "X", URL: "https://en.wikipedia.org/wiki/Foo"}, &http.Client{})
	if r.Status != "bad_domain" {
		t.Fatalf("status=%s", r.Status)
	}
}

func TestAuditOneSourceURL_ok(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	r := auditOneSourceURL(SourceRef{Title: "X", URL: srv.URL}, srv.Client())
	if r.Status != "ok" {
		t.Fatalf("status=%s evidence=%s", r.Status, r.Evidence)
	}
}

func TestScanUnlinkedCitations(t *testing.T) {
	blocks := []EntryBlock{{
		Kind:    "section",
		Section: "Nghiên cứu",
		Content: "Chu Xuân Giao (*Tạp chí Nghiên cứu Tôn giáo*, số 1/2014) lưu ý thêm.\n\n[Cao Huy Đỉnh (1969)](https://example.org/a) đã nêu.",
	}}
	issues := scanUnlinkedCitations(blocks)
	if len(issues) != 1 {
		t.Fatalf("want 1 unlinked cite, got %d", len(issues))
	}
}
