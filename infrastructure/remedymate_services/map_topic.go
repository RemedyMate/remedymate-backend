package remedymate_services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/util"
	"strings"
)

type MapTopicService struct {
	apiKey  string
	baseURL string
	model   string
}

// NewGeminiRemedyRepo creates a new repository for Gemini interactions.
func NewMapTopicService(apiKey, model string) interfaces.MapTopicService {
	return &MapTopicService{
		apiKey:  apiKey,
		baseURL: "https://generativelanguage.googleapis.com/v1beta/models/",
		model:   model,
	}
}

// MapSymptomToTopic implements the Usecase interface method.
func (r *MapTopicService) MapSymptomToTopic(ctx context.Context, userInput string, availableTopics []string) (string, error) {
	prompt := r.BuildMapTopicPrompt(userInput, availableTopics)

	payload := r.CreatePayload(prompt)
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("payload marshaling failed: %w", err)
	}

	respBody, err := r.ExecuteAPIRequest(ctx, body)
	if err != nil {
		return "", err
	}

	topicKey, err := r.ExtractTopicKeyResponse(respBody)
	if err != nil {
		return "", fmt.Errorf("failed to extract topic key: %w", err)
	}

	return topicKey, nil
}

// buildMapTopicPrompt creates the specific prompt for the classification task.
func (r *MapTopicService) BuildMapTopicPrompt(userInput string, availableTopics []string) string {
	// Convert the slice of topics into a formatted string for the prompt
	topicListString := "[\n"
	for _, topic := range availableTopics {
		topicListString += fmt.Sprintf("  \"%s\",\n", topic)
	}
	topicListString += "]"

	return fmt.Sprintf(`
You are an expert AI assistant for a health advisory app. Your task is to analyze the user's symptoms and map them to the single most relevant topic from the provided list.

**Instructions:**
1. Read the user's symptom description carefully.
2. You MUST choose exactly one topic key from the list.
3. If the user's query is vague or does not fit any topic well, you MUST return 'DOES NOT FIT IN ANY TOPIC'.
4. Your response MUST be a single, valid JSON object in the format: {"topic_key": "your_chosen_key"}
5. Do not add any other text, explanations, or markdown formatting around the JSON object.

**Available Topic List:**
%s

**User's Symptom:**
"%s"

**Your JSON Response:**
`, topicListString, userInput)
}

// createPayload builds the JSON payload for the Gemini API call.
func (r *MapTopicService) CreatePayload(prompt string) map[string]any {
	return map[string]any{
		"contents": []any{
			map[string]any{
				"parts": []any{
					map[string]any{"text": prompt},
				},
			},
		},
		// Safety settings are important for a health app
		"safetySettings": []map[string]string{
			{"category": "HARM_CATEGORY_DANGEROUS_CONTENT", "threshold": "BLOCK_ONLY_HIGH"},
			{"category": "HARM_CATEGORY_HATE_SPEECH", "threshold": "BLOCK_MEDIUM_AND_ABOVE"},
			{"category": "HARM_CATEGORY_HARASSMENT", "threshold": "BLOCK_MEDIUM_AND_ABOVE"},
			{"category": "HARM_CATEGORY_SEXUALLY_EXPLICIT", "threshold": "BLOCK_MEDIUM_AND_ABOVE"},
		},
		"generationConfig": map[string]any{
			"maxOutputTokens": 100, // Small, as we only expect a short JSON response
			"temperature":     0.1, // Low temperature for deterministic classification
			"topP":            0.8,
			"topK":            10,
		},
	}
}

// executeAPIRequest sends the request to the Gemini API.
func (r *MapTopicService) ExecuteAPIRequest(ctx context.Context, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s:generateContent?key=%s", r.baseURL, r.model, r.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(errorBody))
	}

	return io.ReadAll(resp.Body)
}

// extractTopicKeyResponse parses the Gemini API response to get the topic key.
func (r *MapTopicService) ExtractTopicKeyResponse(respBody []byte) (string, error) {
	// Define structs to match the Gemini API response structure
	type apiPart struct {
		Text string `json:"text"`
	}
	type apiContent struct {
		Parts []apiPart `json:"parts"`
	}
	type apiCandidate struct {
		Content apiContent `json:"content"`
	}
	type apiResponse struct {
		Candidates []apiCandidate `json:"candidates"`
	}

	var geminiResp apiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return "", fmt.Errorf("gemini response parsing failed: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content returned from Gemini API, possibly due to safety filters")
	}

	// The text part from Gemini should be our target JSON object
	jsonText := geminiResp.Candidates[0].Content.Parts[0].Text

	// Clean the response: LLMs sometimes wrap JSON in markdown code blocks
	jsonText = strings.TrimSpace(jsonText)
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)

	// Define a struct to unmarshal the final JSON object
	type topicResponse struct {
		TopicKey string `json:"topic_key"`
	}

	var finalResp topicResponse
	if err := json.Unmarshal([]byte(jsonText), &finalResp); err != nil {
		return "", fmt.Errorf("failed to parse topic key JSON from LLM response: %w. Response text: %s", err, jsonText)
	}

	if err := util.ValidateTopicKey(finalResp.TopicKey); err != nil {
		return "", err
	}


	return finalResp.TopicKey, nil
}
