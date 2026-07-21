package main

import (
	"context"
	"fmt"
	"os"
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
// Extract: cheap/fast JSON parsing. Judge: strict instruction-following (Claude).
type Runtime struct {
	Extract         LLM
	Judge           LLM
	Search          SearchProvider
	ExtractProvider string
	JudgeProvider   string
}

// AIConfig selects LLM and search providers from CLI / env.
type AIConfig struct {
	ExtractProvider string // default deepseek
	ExtractModel    string
	JudgeProvider   string // default claude
	JudgeModel      string
	SearchProvider  string // default perplexity

	// Legacy: --llm / AUDIT_LLM overrides both extract and judge (batch experiment).
	LLMProvider string
	LLMModel    string
}

func NewRuntime(cfg AIConfig) (Runtime, error) {
	extractProv := cfg.ExtractProvider
	extractModel := cfg.ExtractModel
	judgeProv := cfg.JudgeProvider
	judgeModel := cfg.JudgeModel

	if cfg.LLMProvider != "" {
		extractProv = cfg.LLMProvider
		judgeProv = cfg.LLMProvider
		if cfg.LLMModel != "" {
			extractModel = cfg.LLMModel
			judgeModel = cfg.LLMModel
		}
	}
	if extractProv == "" {
		extractProv = "deepseek"
	}
	if judgeProv == "" {
		judgeProv = "claude"
	}

	extract, err := NewLLM(extractProv, extractModel)
	if err != nil {
		return Runtime{}, fmt.Errorf("extract LLM: %w", err)
	}
	judge, err := NewLLM(judgeProv, judgeModel)
	if err != nil {
		return Runtime{}, fmt.Errorf("judge LLM: %w", err)
	}
	search, err := NewSearch(cfg.SearchProvider)
	if err != nil {
		return Runtime{}, err
	}

	warnJudgeProvider(judgeProv)

	return Runtime{
		Extract:         extract,
		Judge:           judge,
		Search:          search,
		ExtractProvider: extractProv,
		JudgeProvider:   judgeProv,
	}, nil
}

// warnJudgeProvider flags risky judge choices for Vietnamese/Hán Nôm encyclopedia work.
func warnJudgeProvider(provider string) {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "gemini", "deepseek":
		fmt.Fprintf(os.Stderr,
			"warning: judge=%s is not recommended for Vietnamese/Hán Nôm fact-checking; use --judge-llm claude\n",
			provider,
		)
	}
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
