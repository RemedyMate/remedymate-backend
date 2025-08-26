package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/RemedyMate/remedymate-backend/domain/dto"
	"github.com/RemedyMate/remedymate-backend/domain/entities"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
)

type RemedyMateUsecase struct {
	triageService  interfaces.TriageService
	contentService interfaces.ContentService
}

// NewRemedyMateUsecase creates a new RemedyMate usecase
func NewRemedyMateUsecase(
	triageService interfaces.TriageService,
	contentService interfaces.ContentService,

) interfaces.RemedyMateUsecase {
	return &RemedyMateUsecase{
		triageService:  triageService,
		contentService: contentService,
	}
}

// GetTriage performs only triage classification
func (rmu *RemedyMateUsecase) GetTriage(ctx context.Context, req dto.TriageRequest) (*dto.TriageResponse, error) {
	input := entities.SymptomInput{
		Text:     req.Text,
		Language: req.Language,
	}

	result, err := rmu.triageService.ClassifySymptoms(ctx, input)
	if err != nil {
		return nil, err
	}

	return &dto.TriageResponse{
		Level:     result.Level,
		RedFlags:  result.RedFlags,
		Message:   result.Message,
		SessionID: generateSessionID(),
	}, nil
}

// GetContent retrieves approved content for a given topic and language
func (rmu *RemedyMateUsecase) GetContent(ctx context.Context, topicKey, language string) (*entities.ContentTranslation, error) {
	// Validate topic key and language
	if topicKey == "" {
		return nil, fmt.Errorf("topic key is required")
	}
	if language == "" {
		return nil, fmt.Errorf("language is required")
	}
	
	// Validate language is supported
	if language != "en" && language != "am" {
		return nil, fmt.Errorf("unsupported language: %s. Supported languages: en, am", language)
	}

	// Get content from content service
	content, err := rmu.contentService.GetContentByTopic(topicKey, language)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	return content, nil
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// fallback to timestamp if random fails
		return fmt.Sprintf("session_%d", time.Now().UnixNano())
	}
	return "session_" + hex.EncodeToString(b)
}
