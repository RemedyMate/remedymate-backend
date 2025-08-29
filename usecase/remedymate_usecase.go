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

var validTopicKeys = []string{
	"indigestion", "headache", "sore_throat", "cough", "fever", "back_pain",
}

type RemedyMateUsecase struct {
	triageService    interfaces.TriageService
	contentService   interfaces.ContentService
	guidanceComposer interfaces.GuidanceComposerService
	mapService       interfaces.MapTopicService
}

// NewRemedyMateUsecase creates a new RemedyMate usecase
func NewRemedyMateUsecase(
	triageService interfaces.TriageService,
	contentService interfaces.ContentService,
	guidanceComposer interfaces.GuidanceComposerService,
	mapService interfaces.MapTopicService,
) interfaces.RemedyMateUsecase {
	return &RemedyMateUsecase{
		triageService:    triageService,
		contentService:   contentService,
		guidanceComposer: guidanceComposer,
		mapService:       mapService,
	}
}

// GetTriage performs only triage classification
func (rmu *RemedyMateUsecase) GetTriage(ctx context.Context, text, lang string) (*dto.TriageResponse, error) {
	result, err := rmu.triageService.ClassifySymptoms(ctx, text, lang)
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

// MapTopic maps user symptom input to a valid topic key
func (rmu *RemedyMateUsecase) MapTopic(ctx context.Context, input string) (string, error) {
	topicKey, err := rmu.mapService.MapSymptomToTopic(ctx, input, validTopicKeys)
	if err != nil {
		return "", fmt.Errorf("failed to map symptom to topic: %w", err)
	}

	// Validate the returned topic key
	isValid := false
	for _, validKey := range validTopicKeys {
		if topicKey == validKey {
			isValid = true
			break
		}
	}
	if !isValid {
		return "", fmt.Errorf("invalid topic key returned: %s", topicKey)
	}

	return topicKey, nil
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

// ComposeGuidance performs only guidance composition
func (rmu *RemedyMateUsecase) ComposeGuidance(ctx context.Context, req dto.ComposeRequest) (*dto.ComposeResponse, error) {
	guidanceCard, err := rmu.guidanceComposer.ComposeGuidance(ctx, req.TopicKey, req.Language)
	if err != nil {
		return nil, err
	}

	return &dto.ComposeResponse{
		GuidanceCard: *guidanceCard,
		SessionID:    generateSessionID(),
	}, nil
}
