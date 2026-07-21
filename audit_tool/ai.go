package main

import (
	"context"
	"fmt"
	"strings"
)

// LLM completes a system + user prompt. Implementations: Claude, OpenAI, DeepSeek, Gemini.
type LLM interface {
	Complete(ctx context.Context, system, prompt string, maxTokens int) (string, error)
}

// SearchProvider runs a grounded web search for claim verification.
type SearchProvider interface {
	Search(ctx context.Context, query string) (*SearchResult, error)
}

// SearchResult is normalized output from any search backend.
type SearchResult struct {
	Answer    string
	Citations []string
}

// Runtime holds injectable AI backends for the audit pipeline.
type Runtime struct {
	LLM    LLM
	Search SearchProvider
}

// AIConfig selects LLM and search providers from CLI / env.
type AIConfig struct {
	LLMProvider    string // claude | openai | deepseek | gemini
	LLMModel       string // optional override
	SearchProvider string // perplexity (default)
}

func NewRuntime(cfg AIConfig) (Runtime, error) {
	llm, err := NewLLM(cfg.LLMProvider, cfg.LLMModel)
	if err != nil {
		return Runtime{}, err
	}
	search, err := NewSearch(cfg.SearchProvider)
	if err != nil {
		return Runtime{}, err
	}
	return Runtime{LLM: llm, Search: search}, nil
}

func NewLLM(provider, model string) (LLM, error) {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "", "claude":
		return newClaudeLLM(model)
	case "openai":
		return newOpenAILLM(model, openAIEndpoint, "OPENAI_API_KEY")
	case "deepseek":
		return newOpenAILLM(model, deepSeekEndpoint, "DEEPSEEK_API_KEY")
	case "gemini":
		return newGeminiLLM(model)
	default:
		return nil, fmt.Errorf("unknown LLM provider %q (use claude, openai, deepseek, gemini)", provider)
	}
}

func NewSearch(provider string) (SearchProvider, error) {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "", "perplexity":
		return newPerplexitySearch(), nil
	default:
		return nil, fmt.Errorf("unknown search provider %q (use perplexity)", provider)
	}
}

func callLLMJSON(llm LLM, systemPrompt, userPrompt string, maxTokens int) (string, error) {
	// ponytail: retry once for truncated JSON; switch to streaming/tool-use only if this proves flaky.
	var lastRaw string
	var lastErr error
	for range 2 {
		raw, err := llm.Complete(context.Background(), systemPrompt, userPrompt, maxTokens)
		if err != nil {
			return "", err
		}
		raw = cleanLLMJSON(raw)
		lastRaw = raw

		if isJSONArray(raw) || isJSONObject(raw) {
			return raw, nil
		}
		lastErr = fmt.Errorf("parse LLM JSON")
	}
	return "", fmt.Errorf("%w\nraw output: %s", lastErr, lastRaw)
}

func cleanLLMJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	if start := strings.Index(raw, "["); start >= 0 {
		if end := strings.LastIndex(raw, "]"); end > start {
			return strings.TrimSpace(raw[start : end+1])
		}
	}
	if start := strings.Index(raw, "{"); start >= 0 {
		if end := strings.LastIndex(raw, "}"); end > start {
			return strings.TrimSpace(raw[start : end+1])
		}
	}
	return raw
}
