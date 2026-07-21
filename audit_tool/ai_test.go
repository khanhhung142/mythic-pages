package main

import (
	"os"
	"testing"
)

func TestNewLLM_unknownProvider(t *testing.T) {
	_, err := NewLLM("unknown", "")
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestNewRuntime_splitDefaults(t *testing.T) {
	if os.Getenv("DEEPSEEK_API_KEY") == "" {
		t.Skip("DEEPSEEK_API_KEY not set")
	}
	if os.Getenv("ANTHROPIC_API_KEY_API_PLATFORM") == "" {
		t.Skip("ANTHROPIC_API_KEY_API_PLATFORM not set")
	}
	rt, err := NewRuntime(AIConfig{})
	if err != nil {
		t.Fatal(err)
	}
	if rt.Extract == nil || rt.Judge == nil || rt.Search == nil {
		t.Fatal("expected non-nil Extract, Judge, Search")
	}
	if rt.ExtractProvider != "deepseek" || rt.JudgeProvider != "claude" {
		t.Fatalf("extract=%s judge=%s", rt.ExtractProvider, rt.JudgeProvider)
	}
}

func TestNewRuntime_legacyLLMOverride(t *testing.T) {
	if os.Getenv("DEEPSEEK_API_KEY") == "" {
		t.Skip("DEEPSEEK_API_KEY not set")
	}
	rt, err := NewRuntime(AIConfig{LLMProvider: "deepseek"})
	if err != nil {
		t.Fatal(err)
	}
	if rt.ExtractProvider != "deepseek" || rt.JudgeProvider != "deepseek" {
		t.Fatalf("legacy --llm should set both: extract=%s judge=%s", rt.ExtractProvider, rt.JudgeProvider)
	}
}

func TestCleanLLMJSON(t *testing.T) {
	raw := "```json\n[{\"id\":1}]\n```"
	got := cleanLLMJSON(raw)
	if got != `[{"id":1}]` {
		t.Fatalf("got %q", got)
	}
}
