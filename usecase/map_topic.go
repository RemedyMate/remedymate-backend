package usecase

import (
	"context"
	"fmt"
	// "encoding/json"
	// "fmt"
	// "net/http"
)

var validTopicKeys = []string{
	"indigestion", "headache", "sore_throat", "cough", "fever", "back_pain",
}

type RemedyAIRepository interface {
	MapSymptomToTopic(ctx context.Context, userInput string, availableTopics []string) (string, error)
}

type RemedyUsecase struct {
	remedyRepo RemedyAIRepository
}

func NewRemedyUsecase(repo RemedyAIRepository) *RemedyUsecase {
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
