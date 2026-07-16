package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const perplexityEndpoint = "https://api.perplexity.ai/chat/completions"
const perplexityModel = "sonar"

type perplexityRequest struct {
	Model    string               `json:"model"`
	Messages []perplexityMessage  `json:"messages"`
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

type searchResult struct {
	Answer    string
	Citations []string
}

func searchClaim(claim string) (*searchResult, error) {
	apiKey := os.Getenv("PERPLEXITY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("PERPLEXITY_API_KEY not set")
	}

	prompt := fmt.Sprintf(`Verify this claim about Vietnamese history/mythology. Be specific: does evidence support, contradict, or not address it?

Claim: %s

Respond concisely:
1. Verdict: supported / contradicted / not_found / uncertain
2. Evidence: what sources say
3. Key source titles if found`, claim)

	req := perplexityRequest{
		Model: perplexityModel,
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

	httpReq, err := http.NewRequest("POST", perplexityEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
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
		return nil, fmt.Errorf("parse error: %w\nraw: %s", err, raw)
	}
	if pr.Error != nil {
		return nil, fmt.Errorf("perplexity error: %s", pr.Error.Message)
	}
	if len(pr.Choices) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	return &searchResult{
		Answer:    pr.Choices[0].Message.Content,
		Citations: pr.Citations,
	}, nil
}
