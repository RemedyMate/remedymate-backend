package dto

import "net/http"

// LLMConfig holds configuration for LLM client
type LLMConfig struct {
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float32
	Timeout     int // seconds
}

// GeminiClient implements LLMClient using Gemini API
type GeminiClient struct {
	config     LLMConfig
	httpClient *http.Client
	baseURL    string
}

// Gemini-specific request/response structures
type GeminiRequest struct {
	Contents         []GeminiContent       `json:"contents"`
	GenerationConfig GeminiGenConfig       `json:"generationConfig"`
	SafetySettings   []GeminiSafetySetting `json:"safetySettings,omitempty"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiGenConfig struct {
	Temperature     float32 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
	TopP            float32 `json:"topP"`
}

type GeminiSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
	Error      *APIError         `json:"error,omitempty"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}
