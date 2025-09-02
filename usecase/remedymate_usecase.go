package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	derrors "remedymate-backend/domain/AppError"
	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/util"
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

	// If mapping returns an empty topic key, surface a clear error for the controller
	if err := util.ValidateTopicKey(topicKey); err != nil {
		return "", err
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
	if err := util.ValidateTopicKey(topicKey); err != nil {
		return nil, err
	}
	if err := util.ValidateLanguage(language); err != nil {
		return nil, err
	}

	// Get content from content service
	content, err := rmu.contentService.GetContentByTopic(topicKey, language)
	if err != nil {
		// Use sentinel errors from domain/errors for robust handling
		switch {
		case errors.Is(err, derrors.ErrTopicNotFound):
			return nil, fmt.Errorf("topic '%s' not found", topicKey)
		case errors.Is(err, derrors.ErrLanguageNotAvailable):
			return nil, fmt.Errorf("language '%s' not available for topic '%s'", language, topicKey)
		default:
			return nil, err
		}
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

// GetRemedy orchestrates triage, topic mapping, content retrieval and LLM composition.
func (rmu *RemedyMateUsecase) GetRemedy(ctx context.Context, req dto.RemedyRequest) (*dto.RemedyResponse, error) {
	// 1) Triage
	triageRes, err := rmu.triageService.ClassifySymptoms(ctx, req.Text, req.Language)
	if err != nil {
		return nil, err
	}
	// Build base response and embed triage result
	base := dto.RemedyResponse{
		SessionID: generateSessionID(),
	}
	base.Triage = dto.TriageResponse{
		Level:    triageRes.Level,
		RedFlags: triageRes.RedFlags,
		Message:  triageRes.Message,
	}

	// If RED, return early with triage only
	if triageRes.Level == entities.TriageLevelRed {
		return &base, nil
	}

	// 2) Map topic (use existing MapTopic which validates the key)
	topicKey, err := rmu.MapTopic(ctx, req.Text)
	if err != nil {
		return nil, err
	}

	// 3) Get content blocks
	content, err := rmu.contentService.GetContentByTopic(topicKey, req.Language)
	if err != nil {
		return nil, err
	}

	// 4) Compose guidance
	guidanceCard, err := rmu.guidanceComposer.ComposeFromBlocks(ctx, topicKey, req.Language, *content)
	if err != nil {
		return nil, fmt.Errorf("failed to compose guidance: %w", err)
	}

	base.Content = guidanceCard
	return &base, nil
}
