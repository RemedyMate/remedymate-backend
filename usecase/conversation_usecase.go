package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
)

type ConversationUsecaseImpl struct {
	conversationService interfaces.ConversationService
	conversationRepo    interfaces.ConversationRepository
	remedyMateUsecase   interfaces.RemedyMateUsecase
}

// NewConversationUsecase creates a new conversation usecase
func NewConversationUsecase(
	conversationService interfaces.ConversationService,
	conversationRepo interfaces.ConversationRepository,
	remedyMateUsecase interfaces.RemedyMateUsecase,
) interfaces.ConversationUsecase {
	return &ConversationUsecaseImpl{
		conversationService: conversationService,
		conversationRepo:    conversationRepo,
		remedyMateUsecase:   remedyMateUsecase,
	}
}

// ValidateSymptom validates if the provided symptom is medical and appropriate
func (cu *ConversationUsecaseImpl) ValidateSymptom(ctx context.Context, symptom, language string) (bool, string, error) {
	return cu.conversationService.ValidateSymptom(ctx, symptom, language)
}

// StartConversation starts a new conversation with the initial symptom
func (cu *ConversationUsecaseImpl) StartConversation(ctx context.Context, req dto.StartConversationRequest) (*dto.StartConversationResponse, error) {
	// Generate questions using AI
	questions, err := cu.conversationService.GenerateQuestions(ctx, req.Symptom, req.Language)
	if err != nil {
		return nil, fmt.Errorf("failed to generate questions: %w", err)
	}

	// Create conversation entity
	conversation := &entities.Conversation{
		ID:          generateConversationID(),
		UserID:      req.UserID,
		Symptom:     req.Symptom,
		Language:    req.Language,
		Questions:   questions,
		Answers:     []entities.Answer{},
		TotalSteps:  len(questions),
		CurrentStep: 1,
	}

	// Save conversation to database
	err = cu.conversationRepo.CreateConversation(ctx, conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Return first question
	firstQuestion := questions[0]
	return &dto.StartConversationResponse{
		ConversationID: conversation.ID,
		Question:       firstQuestion,
		TotalSteps:     conversation.TotalSteps,
		CurrentStep:    conversation.CurrentStep,
	}, nil
}

// SubmitAnswer submits an answer to the current question
func (cu *ConversationUsecaseImpl) SubmitAnswer(ctx context.Context, req dto.SubmitAnswerRequest) (*dto.SubmitAnswerResponse, error) {
	// Get conversation from database
	conversation, err := cu.conversationRepo.GetConversation(ctx, req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}

	// Check if conversation is still active
	if conversation.Status != entities.ConversationStatusActive {
		return nil, fmt.Errorf("conversation is not active")
	}

	// Get current question
	if conversation.CurrentStep > len(conversation.Questions) {
		return nil, fmt.Errorf("invalid conversation state")
	}

	currentQuestion := conversation.Questions[conversation.CurrentStep-1]

	// Validate answer using AI
	isValid, feedback, err := cu.conversationService.ValidateAnswer(ctx, currentQuestion, req.Answer)
	if err != nil {
		return nil, fmt.Errorf("failed to validate answer: %w", err)
	}

	// Create answer entity
	answer := entities.Answer{
		QuestionID: currentQuestion.ID,
		Text:       req.Answer,
		IsValid:    isValid,
		Feedback:   feedback,
		AnsweredAt: time.Now(),
	}

	// Add answer to conversation
	err = cu.conversationRepo.AddAnswer(ctx, req.ConversationID, answer)
	if err != nil {
		return nil, fmt.Errorf("failed to save answer: %w", err)
	}

	// If answer is invalid, return same question with feedback
	if !isValid {
		return &dto.SubmitAnswerResponse{
			ConversationID: req.ConversationID,
			Question:       &currentQuestion,
			Message:        feedback,
			IsComplete:     false,
			CurrentStep:    conversation.CurrentStep,
			TotalSteps:     conversation.TotalSteps,
		}, nil
	}

	// Move to next question
	conversation.CurrentStep++

	// Check if all questions are answered
	if conversation.CurrentStep > len(conversation.Questions) {
		// Generate final health report
		report, err := cu.conversationService.GenerateHealthReport(ctx, conversation)
		if err != nil {
			return nil, fmt.Errorf("failed to generate health report: %w", err)
		}

		// Get remedy using the original symptom and language
		remedyReq := dto.RemedyRequest{
			Text:     conversation.Symptom,
			Language: conversation.Language,
		}
		remedyResponse, err := cu.remedyMateUsecase.GetRemedy(ctx, remedyReq)

		if err != nil {
			// Log the error but don't fail the conversation completion
			// The conversation can still complete without remedy
			fmt.Printf("Warning: Failed to get remedy for conversation %s: %v\n", req.ConversationID, err)
			remedyResponse = nil
		}

		// Save final report with remedy information
		if remedyResponse != nil && remedyResponse.Content != nil {
			// Add remedy information to the report
			report.Remedy = &entities.Remedy{
				Triage: entities.TriageResult{
					Level:    remedyResponse.Triage.Level,
					RedFlags: remedyResponse.Triage.RedFlags,
					Message:  remedyResponse.Triage.Message,
				},
				SelfCare:      remedyResponse.Content.SelfCare,
				OTCCategories: remedyResponse.Content.OTCCategories,
				SeekCareIf:    remedyResponse.Content.SeekCareIf,
				Disclaimer:    remedyResponse.Content.Disclaimer,
				TopicKey:      remedyResponse.Content.TopicKey,
				Language:      remedyResponse.Content.Language,
			}
		}

		// Save final report
		err = cu.conversationRepo.SetFinalReport(ctx, req.ConversationID, report)
		if err != nil {
			return nil, fmt.Errorf("failed to save final report: %w", err)
		}

		// Mark conversation as complete
		err = cu.conversationRepo.UpdateConversationStatus(ctx, req.ConversationID, entities.ConversationStatusComplete)
		if err != nil {
			return nil, fmt.Errorf("failed to update conversation status: %w", err)
		}

		return &dto.SubmitAnswerResponse{
			ConversationID: req.ConversationID,
			Question:       nil,
			Message:        "All questions completed. You can now view your health report and remedy.",
			IsComplete:     true,
			CurrentStep:    conversation.TotalSteps,
			TotalSteps:     conversation.TotalSteps,
		}, nil
	}

	// Update conversation progress
	err = cu.conversationRepo.UpdateConversation(ctx, conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to update conversation: %w", err)
	}

	// Return next question
	nextQuestion := conversation.Questions[conversation.CurrentStep-1]
	return &dto.SubmitAnswerResponse{
		ConversationID: req.ConversationID,
		Question:       &nextQuestion,
		Message:        "",
		IsComplete:     false,
		CurrentStep:    conversation.CurrentStep,
		TotalSteps:     conversation.TotalSteps,
	}, nil
}

// GetReport retrieves the final health report for a completed conversation
func (cu *ConversationUsecaseImpl) GetReport(ctx context.Context, conversationID string) (*dto.GetReportResponse, error) {
	// Get conversation from database
	conversation, err := cu.conversationRepo.GetConversation(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}

	// Check if conversation is complete
	if conversation.Status != entities.ConversationStatusComplete {
		return nil, fmt.Errorf("conversation is not complete")
	}

	// Check if final report exists
	if conversation.FinalReport == nil {
		return nil, fmt.Errorf("final report not found")
	}

	return &dto.GetReportResponse{
		ConversationID: conversationID,
		Report:         conversation.FinalReport,
		Symptom:        conversation.Symptom,
		Status:         string(conversation.Status),
	}, nil
}

// generateConversationID generates a unique conversation ID
func generateConversationID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// fallback to timestamp if random fails
		return fmt.Sprintf("conv_%d", time.Now().UnixNano())
	}
	return "conv_" + hex.EncodeToString(b)
}

func (cu *ConversationUsecaseImpl) GetOfflineHealthTopics(ctx context.Context) ([]entities.HealthTopic, error) {
	return cu.conversationRepo.GetOfflineHealthTopics(ctx)
}
