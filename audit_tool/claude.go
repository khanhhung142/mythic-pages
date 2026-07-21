package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const claudeModel = "claude-sonnet-4-6"
const claudeEndpoint = "https://api.anthropic.com/v1/messages"

// Explicit cache on the system block (static across calls). Do NOT use
// top-level automatic cache_control — that breakpoints on the last user
// message, which changes every claim and never reuses the system prefix.
// https://platform.claude.com/docs/en/build-with-claude/prompt-caching
type cacheControl struct {
	Type string `json:"type"`
	TTL  string `json:"ttl,omitempty"` // "5m" (default) or "1h"
}

type systemBlock struct {
	Type         string        `json:"type"`
	Text         string        `json:"text"`
	CacheControl *cacheControl `json:"cache_control,omitempty"`
}

type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	System    []systemBlock   `json:"system,omitempty"`
	Messages  []claudeMessage `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func callClaude(system, prompt string, maxTokens int) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY_API_PLATFORM")
	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY_API_PLATFORM not set")
	}

	req := claudeRequest{
		Model:     claudeModel,
		MaxTokens: maxTokens,
		System: []systemBlock{{
			Type: "text",
			Text: system,
			// 1h: audit runs can exceed the default 5m TTL across many claims.
			CacheControl: &cacheControl{Type: "ephemeral", TTL: "1h"},
		}},
		Messages: []claudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", claudeEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var cr claudeResponse
	if err := json.Unmarshal(raw, &cr); err != nil {
		return "", fmt.Errorf("parse error: %w\nraw: %s", err, raw)
	}
	if cr.Error != nil {
		return "", fmt.Errorf("claude error: %s", cr.Error.Message)
	}
	if len(cr.Content) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return cr.Content[0].Text, nil
}
