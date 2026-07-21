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
	openAIEndpoint       = "https://api.openai.com/v1/chat/completions"
	deepSeekEndpoint     = "https://api.deepseek.com/chat/completions"
	defaultOpenAIModel   = "gpt-4o"
	defaultDeepSeekModel = "deepseek-chat"
)

type openAILLM struct {
	model    string
	apiKey   string
	endpoint string
	client   *http.Client
}

type openAIRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func newOpenAILLM(model, endpoint, apiKeyEnv string) (*openAILLM, error) {
	apiKey := os.Getenv(apiKeyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf("%s not set", apiKeyEnv)
	}
	if model == "" {
		if endpoint == deepSeekEndpoint {
			model = defaultDeepSeekModel
		} else {
			model = defaultOpenAIModel
		}
	}
	return &openAILLM{
		model:    model,
		apiKey:   apiKey,
		endpoint: endpoint,
		client:   http.DefaultClient,
	}, nil
}

func (o *openAILLM) Complete(_ context.Context, system, prompt string, maxTokens int) (string, error) {
	req := openAIRequest{
		Model:     o.model,
		MaxTokens: maxTokens,
		Messages: []openAIMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest(http.MethodPost, o.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var or openAIResponse
	if err := json.Unmarshal(raw, &or); err != nil {
		return "", fmt.Errorf("openai-compatible parse error: %w\nraw: %s", err, raw)
	}
	if or.Error != nil {
		return "", fmt.Errorf("openai-compatible error: %s", or.Error.Message)
	}
	if len(or.Choices) == 0 {
		return "", fmt.Errorf("openai-compatible empty response")
	}
	return or.Choices[0].Message.Content, nil
}
