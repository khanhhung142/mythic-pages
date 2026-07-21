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

func TestNewRuntime_defaults(t *testing.T) {
	if os.Getenv("ANTHROPIC_API_KEY_API_PLATFORM") == "" {
		t.Skip("ANTHROPIC_API_KEY_API_PLATFORM not set")
	}
	rt, err := NewRuntime(AIConfig{})
	if err != nil {
		t.Fatal(err)
	}
	if rt.LLM == nil || rt.Search == nil {
		t.Fatal("expected non-nil LLM and Search")
	}
}

func TestCleanLLMJSON(t *testing.T) {
	raw := "```json\n[{\"id\":1}]\n```"
	got := cleanLLMJSON(raw)
	if got != `[{"id":1}]` {
		t.Fatalf("got %q", got)
	}
}
