package usecase

import (
	"context"
	"fmt"
	"remedymate-backend/domain/interfaces"
)

var validTopicKeys = []string{
	"indigestion", "headache", "sore_throat", "cough", "fever", "back_pain",
}

type RemedyUsecase struct {
	remedyRepo interfaces.RemedyAIRepository
}

func NewRemedyUsecase(repo interfaces.RemedyAIRepository) *RemedyUsecase {
	return &RemedyUsecase{remedyRepo: repo}
}

func (uc *RemedyUsecase) MapTopic(ctx context.Context, input string) (string, error) {
	topicKey, err := uc.remedyRepo.MapSymptomToTopic(ctx, input, validTopicKeys)
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
