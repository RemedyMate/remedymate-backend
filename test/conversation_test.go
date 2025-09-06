package test

import (
	"context"
	"strings"
	"testing"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/infrastructure/conversation"

	"github.com/stretchr/testify/mock"
)

// MockLLMClient is a mock implementation of LLMClient for testing
type MockLLMClient struct {
	mock.Mock
}

func (m *MockLLMClient) ClassifyTriage(ctx context.Context, prompt string) (string, error) {
	args := m.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

// TestConversationFlow tests the complete conversation flow
func TestConversationFlow(t *testing.T) {
	// Create mock LLM client
	mockLLM := &MockLLMClient{}

	// Mock responses for question generation - match by content instead of length
	mockLLM.On("ClassifyTriage", mock.Anything, mock.MatchedBy(func(prompt string) bool {
		return strings.Contains(prompt, "Generate exactly 5 targeted follow-up questions")
	})).Return(`[
		{"id": 1, "text": "How long have you had this headache?", "type": "duration", "required": true},
		{"id": 2, "text": "Where is the pain located?", "type": "location", "required": true},
		{"id": 3, "text": "How severe is the pain?", "type": "severity", "required": true},
		{"id": 4, "text": "Any medical history?", "type": "history", "required": false},
		{"id": 5, "text": "What triggers the pain?", "type": "triggers", "required": false}
	]`, nil)

	// Mock responses for answer validation
	mockLLM.On("ClassifyTriage", mock.Anything, mock.MatchedBy(func(prompt string) bool {
		return strings.Contains(prompt, "Validate this answer to a medical question")
	})).Return(`{"valid": true, "feedback": ""}`, nil)

	// Mock response for health report generation
	mockLLM.On("ClassifyTriage", mock.Anything, mock.MatchedBy(func(prompt string) bool {
		return strings.Contains(prompt, "Create a structured health report")
	})).Return(`{
		"symptom": "Headache",
		"duration": "3 days",
		"location": "Front of head",
		"severity": "Moderate",
		"associated_symptoms": ["Nausea"],
		"medical_history": "None",
		"triggers": "Stress",
		"possible_conditions": ["Tension headache"],
		"recommendations": ["Rest", "Pain medication"],
		"urgency_level": "GREEN"
	}`, nil)

	// Create conversation service
	conversationService := conversation.NewConversationService(mockLLM)

	// Test question generation
	questions, err := conversationService.GenerateQuestions(context.Background(), "headache", "en")
	if err != nil {
		t.Fatalf("Failed to generate questions: %v", err)
	}

	if len(questions) != 5 {
		t.Errorf("Expected 5 questions, got %d", len(questions))
	}

	// Test answer validation
	question := entities.Question{ID: 1, Text: "How long have you had this headache?", Type: "duration", Required: true}
	isValid, feedback, err := conversationService.ValidateAnswer(context.Background(), question, "3 days")
	if err != nil {
		t.Fatalf("Failed to validate answer: %v", err)
	}

	if !isValid {
		t.Errorf("Expected valid answer, got invalid with feedback: %s", feedback)
	}

	// Test health report generation
	conv := &entities.Conversation{
		Symptom:   "headache",
		Language:  "en",
		Questions: questions,
		Answers: []entities.Answer{
			{QuestionID: 1, Text: "3 days", IsValid: true},
			{QuestionID: 2, Text: "Front of head", IsValid: true},
			{QuestionID: 3, Text: "Moderate", IsValid: true},
			{QuestionID: 4, Text: "None", IsValid: true},
			{QuestionID: 5, Text: "Stress", IsValid: true},
		},
	}

	report, err := conversationService.GenerateHealthReport(context.Background(), conv)
	if err != nil {
		t.Fatalf("Failed to generate health report: %v", err)
	}

	if report.Symptom != "Headache" {
		t.Errorf("Expected symptom 'Headache', got '%s'", report.Symptom)
	}

	mockLLM.AssertExpectations(t)
}

// TestConversationDTOs tests the DTO structures
func TestConversationDTOs(t *testing.T) {
	// Test StartConversationRequest (no authentication required)
	startReq := dto.StartConversationRequest{
		Symptom:  "headache",
		Language: "en",
		UserID:   "", // Optional for unauthenticated users
	}

	if startReq.Symptom != "headache" {
		t.Errorf("Expected symptom 'headache', got '%s'", startReq.Symptom)
	}

	// Test SubmitAnswerRequest
	answerReq := dto.SubmitAnswerRequest{
		ConversationID: "conv123",
		Answer:         "3 days",
	}

	if answerReq.ConversationID != "conv123" {
		t.Errorf("Expected conversation ID 'conv123', got '%s'", answerReq.ConversationID)
	}
}

// TestConversationEntities tests the entity structures
func TestConversationEntities(t *testing.T) {
	// Test Conversation entity
	conversation := entities.Conversation{
		ID:          "conv123",
		UserID:      "user123",
		Symptom:     "headache",
		Language:    "en",
		Status:      entities.ConversationStatusActive,
		CurrentStep: 1,
		TotalSteps:  5,
	}

	if conversation.Status != entities.ConversationStatusActive {
		t.Errorf("Expected status 'ACTIVE', got '%s'", conversation.Status)
	}

	// Test Question entity
	question := entities.Question{
		ID:       1,
		Text:     "How long have you had this headache?",
		Type:     "duration",
		Required: true,
	}

	if question.Type != "duration" {
		t.Errorf("Expected type 'duration', got '%s'", question.Type)
	}

	// Test Answer entity
	answer := entities.Answer{
		QuestionID: 1,
		Text:       "3 days",
		IsValid:    true,
		Feedback:   "",
	}

	if !answer.IsValid {
		t.Error("Expected valid answer")
	}
}
