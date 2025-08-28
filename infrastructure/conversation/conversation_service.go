package conversation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RemedyMate/remedymate-backend/domain/entities"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
)

type ConversationServiceImpl struct {
	llmClient interfaces.LLMClient
}

// NewConversationService creates a new conversation service
func NewConversationService(llmClient interfaces.LLMClient) interfaces.ConversationService {
	return &ConversationServiceImpl{
		llmClient: llmClient,
	}
}

// GenerateQuestions generates 5 follow-up questions based on the initial symptom
func (cs *ConversationServiceImpl) GenerateQuestions(ctx context.Context, symptom, language string) ([]entities.Question, error) {
	prompt := cs.buildQuestionGenerationPrompt(symptom, language)

	response, err := cs.llmClient.ClassifyTriage(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate questions: %w", err)
	}

	questions, err := cs.parseQuestionsFromResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse questions: %w", err)
	}

	// Ensure we have exactly 5 questions
	if len(questions) < 5 {
		// Generate additional questions if needed
		additionalQuestions := cs.generateDefaultQuestions(symptom, language)
		for i := len(questions); i < 5; i++ {
			questions = append(questions, additionalQuestions[i-len(questions)])
		}
	} else if len(questions) > 5 {
		questions = questions[:5]
	}

	return questions, nil
}

// ValidateAnswer validates a user's answer to a question
func (cs *ConversationServiceImpl) ValidateAnswer(ctx context.Context, question entities.Question, answer string) (bool, string, error) {
	prompt := cs.buildValidationPrompt(question, answer)

	response, err := cs.llmClient.ClassifyTriage(ctx, prompt)
	if err != nil {
		return false, "", fmt.Errorf("failed to validate answer: %w", err)
	}

	isValid, feedback := cs.parseValidationResponse(response)
	return isValid, feedback, nil
}

// GenerateHealthReport creates a structured health report from conversation data
func (cs *ConversationServiceImpl) GenerateHealthReport(ctx context.Context, conversation *entities.Conversation) (*entities.HealthReport, error) {
	prompt := cs.buildReportGenerationPrompt(conversation)

	response, err := cs.llmClient.ClassifyTriage(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate health report: %w", err)
	}

	report, err := cs.parseHealthReportFromResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse health report: %w", err)
	}

	return report, nil
}

// buildQuestionGenerationPrompt creates the prompt for generating follow-up questions
func (cs *ConversationServiceImpl) buildQuestionGenerationPrompt(symptom, language string) string {
	langText := "English"
	if language == "am" {
		langText = "Amharic"
	}

	return fmt.Sprintf(`Generate 5 follow-up questions for a patient with the symptom: "%s". 
	
Requirements:
- Generate exactly 5 questions in %s
- Questions should cover: duration, location, severity, medical history, and triggers
- Each question should be specific and actionable
- Format as JSON array with objects containing: id, text, type, required

Example format:
[
  {"id": 1, "text": "How long have you had this symptom?", "type": "duration", "required": true},
  {"id": 2, "text": "Where exactly is the pain located?", "type": "location", "required": true}
]

Symptom: %s
Language: %s

Generate the questions:`, symptom, langText, symptom, langText)
}

// buildValidationPrompt creates the prompt for validating answers
func (cs *ConversationServiceImpl) buildValidationPrompt(question entities.Question, answer string) string {
	return fmt.Sprintf(`Validate this answer to a medical question.

Question: %s
Question Type: %s
Answer: %s

Requirements:
- Check if the answer is relevant and informative
- For duration: should include time period (days, hours, etc.)
- For location: should specify body part or area
- For severity: should indicate pain level or intensity
- For history: should mention relevant medical background
- For triggers: should describe what causes or worsens the symptom

Respond with JSON format:
{"valid": true/false, "feedback": "explanation if invalid"}

Validation result:`, question.Text, question.Type, answer)
}

// buildReportGenerationPrompt creates the prompt for generating health reports
func (cs *ConversationServiceImpl) buildReportGenerationPrompt(conversation *entities.Conversation) string {
	// Build context from conversation
	context := fmt.Sprintf("Symptom: %s\nLanguage: %s\n", conversation.Symptom, conversation.Language)

	for i, answer := range conversation.Answers {
		if i < len(conversation.Questions) {
			context += fmt.Sprintf("Q%d: %s\nA%d: %s\n",
				answer.QuestionID,
				conversation.Questions[answer.QuestionID-1].Text,
				answer.QuestionID,
				answer.Text)
		}
	}

	return fmt.Sprintf(`Create a structured health report based on this conversation:

%s

Generate a comprehensive health report in JSON format with the following fields:
- symptom: the main symptom
- duration: how long the symptom has been present
- location: where the symptom is located
- severity: how severe the symptom is
- associated_symptoms: any other symptoms mentioned
- medical_history: relevant medical background
- triggers: what causes or worsens the symptom
- possible_conditions: potential diagnoses
- recommendations: suggested next steps
- urgency_level: GREEN/YELLOW/RED based on severity

Format as JSON object.`, context)
}

// parseQuestionsFromResponse parses questions from LLM response
func (cs *ConversationServiceImpl) parseQuestionsFromResponse(response string) ([]entities.Question, error) {
	// Try to extract JSON from the response
	jsonStart := strings.Index(response, "[")
	jsonEnd := strings.LastIndex(response, "]")

	if jsonStart == -1 || jsonEnd == -1 {
		return cs.generateDefaultQuestions("", "en"), nil
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var questions []entities.Question
	err := json.Unmarshal([]byte(jsonStr), &questions)
	if err != nil {
		return cs.generateDefaultQuestions("", "en"), nil
	}

	return questions, nil
}

// parseValidationResponse parses validation result from LLM response
func (cs *ConversationServiceImpl) parseValidationResponse(response string) (bool, string) {
	// Try to extract JSON from the response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return true, "" // Default to valid if can't parse
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var result struct {
		Valid    bool   `json:"valid"`
		Feedback string `json:"feedback"`
	}

	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return true, "" // Default to valid if can't parse
	}

	return result.Valid, result.Feedback
}

// parseHealthReportFromResponse parses health report from LLM response
func (cs *ConversationServiceImpl) parseHealthReportFromResponse(response string) (*entities.HealthReport, error) {
	// Try to extract JSON from the response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("invalid response format")
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var report entities.HealthReport
	err := json.Unmarshal([]byte(jsonStr), &report)
	if err != nil {
		return nil, fmt.Errorf("failed to parse health report: %w", err)
	}

	return &report, nil
}

// generateDefaultQuestions generates default questions if LLM fails
func (cs *ConversationServiceImpl) generateDefaultQuestions(symptom, language string) []entities.Question {
	questions := []entities.Question{
		{ID: 1, Text: "How long have you had this symptom?", Type: "duration", Required: true},
		{ID: 2, Text: "Where exactly is the symptom located?", Type: "location", Required: true},
		{ID: 3, Text: "How severe is the symptom on a scale of 1-10?", Type: "severity", Required: true},
		{ID: 4, Text: "Do you have any relevant medical history?", Type: "history", Required: false},
		{ID: 5, Text: "What makes the symptom better or worse?", Type: "triggers", Required: false},
	}

	// Translate to Amharic if needed
	if language == "am" {
		questions = []entities.Question{
			{ID: 1, Text: "ይህ ምልክት ለምን ያህል ጊዜ አለው?", Type: "duration", Required: true},
			{ID: 2, Text: "ይህ ምልክት የት አለ?", Type: "location", Required: true},
			{ID: 3, Text: "ይህ ምልክት ምን ያህል ከባድ ነው? (1-10)", Type: "severity", Required: true},
			{ID: 4, Text: "ምንም የሕክምና ታሪክ አለዎት?", Type: "history", Required: false},
			{ID: 5, Text: "ምን ያደርገዋል ይህ ምልክት የተሻለ ወይስ የተከሰተ?", Type: "triggers", Required: false},
		}
	}

	return questions
}
