// Package ai talks to the Google Gemini REST API (generateContent).
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultEndpoint = "https://generativelanguage.googleapis.com/v1beta/models"

// Client calls Gemini with an API key.
type Client struct {
	apiKey   string
	model    string
	endpoint string
	http     *http.Client
}

// NewClient returns nil when apiKey is empty (AI features disabled).
func NewClient(apiKey, model string) *Client {
	if apiKey == "" {
		return nil
	}
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &Client{
		apiKey:   apiKey,
		model:    model,
		endpoint: defaultEndpoint,
		http:     &http.Client{Timeout: 30 * time.Second},
	}
}

type part struct {
	Text string `json:"text"`
}

type content struct {
	Parts []part `json:"parts"`
}

type generateRequest struct {
	Contents         []content      `json:"contents"`
	GenerationConfig map[string]any `json:"generationConfig"`
}

type generateResponse struct {
	Candidates []struct {
		Content content `json:"content"`
	} `json:"candidates"`
}

const promptTemplate = `Please provide brief and concise suggestions for the following task description.
Format your response as plain text (no markdown) with:
1. 2-3 suggestions for improving the task description
2. 3-4 specific subtasks

Task description: `

// TaskSuggestions returns cleaned-up plain-text suggestions for the description.
func (c *Client) TaskSuggestions(ctx context.Context, description string) (string, error) {
	body, err := json.Marshal(generateRequest{
		Contents: []content{{Parts: []part{{Text: promptTemplate + description}}}},
		GenerationConfig: map[string]any{
			"temperature":     0.7,
			"topK":            40,
			"topP":            0.95,
			"maxOutputTokens": 1024,
		},
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s:generateContent", c.endpoint, c.model)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	// Key in a header, not the URL, so it never ends up in access logs.
	req.Header.Set("x-goog-api-key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("gemini request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini returned HTTP %d", resp.StatusCode)
	}

	var parsed generateResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("decode gemini response: %w", err)
	}
	if len(parsed.Candidates) == 0 {
		return "", errors.New("no response received from Gemini")
	}

	var sb strings.Builder
	for _, p := range parsed.Candidates[0].Content.Parts {
		sb.WriteString(p.Text)
	}
	text := Clean(sb.String())
	if text == "" {
		return "", errors.New("no text content in Gemini response")
	}
	return text, nil
}

// Clean strips markdown emphasis and blank lines from model output.
func Clean(text string) string {
	text = strings.ReplaceAll(text, "**", "")
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "- ") {
			line = "• " + strings.TrimSpace(line[2:])
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
