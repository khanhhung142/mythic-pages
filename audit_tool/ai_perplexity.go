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
	perplexityEndpoint   = "https://api.perplexity.ai/chat/completions"
	defaultPerplexityModel = "sonar"
)

type perplexitySearch struct {
	model  string
	apiKey string
	client *http.Client
}

type perplexityRequest struct {
	Model    string              `json:"model"`
	Messages []perplexityMessage `json:"messages"`
}

type perplexityMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type perplexityResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Citations []string `json:"citations"`
	Error     *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func newPerplexitySearch() *perplexitySearch {
	return &perplexitySearch{
		model:  defaultPerplexityModel,
		apiKey: os.Getenv("PERPLEXITY_API_KEY"),
		client: http.DefaultClient,
	}
}

func (p *perplexitySearch) Search(_ context.Context, query string) (*SearchResult, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("PERPLEXITY_API_KEY not set")
	}

	prompt := fmt.Sprintf(`Verify this claim about Vietnamese history/mythology. Be specific: does evidence support, contradict, or not address it?

Claim: %s

Respond concisely:
1. Verdict: supported / contradicted / not_found / uncertain
2. Evidence: what sources say
3. Key source titles if found`, query)

	req := perplexityRequest{
		Model: p.model,
		Messages: []perplexityMessage{
			{
				Role:    "system",
				Content: "You are a Vietnamese history and mythology fact-checker. Cite specific primary sources (LNCQ, ĐVSKTT, Việt Điện U Linh, etc.) when available. Be direct and brief.",
			},
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, perplexityEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pr perplexityResponse
	if err := json.Unmarshal(raw, &pr); err != nil {
		return nil, fmt.Errorf("perplexity parse error: %w\nraw: %s", err, raw)
	}
	if pr.Error != nil {
		return nil, fmt.Errorf("perplexity error: %s", pr.Error.Message)
	}
	if len(pr.Choices) == 0 {
		return nil, fmt.Errorf("perplexity empty response")
	}

	return &SearchResult{
		Answer:    pr.Choices[0].Message.Content,
		Citations: pr.Citations,
	}, nil
}
