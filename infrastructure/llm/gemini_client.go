package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/interfaces"
)

// GeminiClient implements LLMClient using Google's Gemini API
type GeminiClient struct {
	config     dto.LLMConfig
	httpClient *http.Client
	baseURL    string
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient(config dto.LLMConfig) interfaces.LLMClient {
	return &GeminiClient{
		config:     config,
		httpClient: &http.Client{Timeout: time.Duration(config.Timeout) * time.Second},
		baseURL:    "https://generativelanguage.googleapis.com/v1beta/models",
	}
}

// ClassifyTriage calls Gemini API for triage classification
func (g *GeminiClient) ClassifyTriage(ctx context.Context, prompt string) (string, error) {
	return g.callGemini(ctx, prompt)
}

// callGemini makes the actual API call to Gemini
func (g *GeminiClient) callGemini(ctx context.Context, prompt string) (string, error) {
	// 1. Construct the request for Gemini's API format
	geminiReq := dto.GeminiRequest{
		Contents: []dto.GeminiContent{
			{
				Parts: []dto.GeminiPart{{Text: prompt}},
				Role:  "user",
			},
		},
		GenerationConfig: dto.GeminiGenConfig{
			Temperature:     g.config.Temperature,
			MaxOutputTokens: g.config.MaxTokens,
			TopP:            0.95,
		},
		SafetySettings: []dto.GeminiSafetySetting{
			{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
			{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "BLOCK_MEDIUM_AND_ABOVE"},
		},
	}

	jsonData, err := json.Marshal(geminiReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 2. Build the HTTP request
	modelURL := fmt.Sprintf("%s/%s:generateContent?key=%s", g.baseURL, g.config.Model, g.config.APIKey)
	req, err := http.NewRequestWithContext(ctx, "POST", modelURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 3. Execute the request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// 4. Parse the Gemini response
	var geminiResp dto.GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if geminiResp.Error != nil {
		return "", fmt.Errorf("gemini API error: %s", geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response content returned from Gemini")
	}

	// 5. Return the generated text
	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
