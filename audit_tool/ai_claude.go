package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	defaultClaudeModel    = "claude-sonnet-4-6"
	claudeMessagesURL     = "https://api.anthropic.com/v1/messages"
	claudeAPIKeyEnv       = "ANTHROPIC_API_KEY_API_PLATFORM"
)

type claudeLLM struct {
	model  string
	apiKey string
	client *http.Client
}

// Explicit cache on the system block (static across calls). Do NOT use
// top-level automatic cache_control — that breakpoints on the last user
// message, which changes every claim and never reuses the system prefix.
// https://platform.claude.com/docs/en/build-with-claude/prompt-caching
type cacheControl struct {
	Type string `json:"type"`
	TTL  string `json:"ttl,omitempty"`
}

type claudeSystemBlock struct {
	Type         string        `json:"type"`
	Text         string        `json:"text"`
	CacheControl *cacheControl `json:"cache_control,omitempty"`
}

type claudeRequest struct {
	Model     string              `json:"model"`
	MaxTokens int                 `json:"max_tokens"`
	System    []claudeSystemBlock `json:"system,omitempty"`
	Messages  []claudeMessage     `json:"messages"`
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

func newClaudeLLM(model string) (*claudeLLM, error) {
	apiKey := os.Getenv(claudeAPIKeyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf("%s not set", claudeAPIKeyEnv)
	}
	if model == "" {
		model = defaultClaudeModel
	}
	return &claudeLLM{model: model, apiKey: apiKey, client: http.DefaultClient}, nil
}

func (c *claudeLLM) Complete(_ context.Context, system, prompt string, maxTokens int) (string, error) {
	req := claudeRequest{
		Model:     c.model,
		MaxTokens: maxTokens,
		System: []claudeSystemBlock{{
			Type:         "text",
			Text:         system,
			CacheControl: &cacheControl{Type: "ephemeral", TTL: "1h"},
		}},
		Messages: []claudeMessage{{Role: "user", Content: prompt}},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest(http.MethodPost, claudeMessagesURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(httpReq)
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
		return "", fmt.Errorf("claude parse error: %w\nraw: %s", err, raw)
	}
	if cr.Error != nil {
		return "", fmt.Errorf("claude error: %s", cr.Error.Message)
	}
	if len(cr.Content) == 0 {
		return "", fmt.Errorf("claude empty response")
	}
	return cr.Content[0].Text, nil
}
