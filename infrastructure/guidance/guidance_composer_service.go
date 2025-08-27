package guidance

import (
	"context"
	"fmt"

	"github.com/RemedyMate/remedymate-backend/domain/entities"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
)

type GuidanceComposerService struct {
	contentService interfaces.ContentService
	llmClient      interfaces.LLMClient
}

// NewGuidanceComposerService creates a new guidance composer service
func NewGuidanceComposerService(contentService interfaces.ContentService, llmClient interfaces.LLMClient) interfaces.GuidanceComposerService {
	return &GuidanceComposerService{
		contentService: contentService,
		llmClient:      llmClient,
	}
}

// ComposeGuidance composes a guidance card for a given topic and language
func (gcs *GuidanceComposerService) ComposeGuidance(ctx context.Context, topicKey, language string) (*entities.GuidanceCard, error) {
	content, err := gcs.contentService.GetContentByTopic(topicKey, language)
	if err != nil {
		return nil, fmt.Errorf("failed to get content for topic %s: %w", topicKey, err)
	}

	return gcs.ComposeFromBlocks(ctx, topicKey, language, *content)
}

// ComposeFromBlocks composes a guidance card from approved content blocks using LLM
func (gcs *GuidanceComposerService) ComposeFromBlocks(ctx context.Context, topicKey, language string, blocks entities.ContentTranslation) (*entities.GuidanceCard, error) {
	if topicKey == "" {
		return nil, fmt.Errorf("topic key cannot be empty")
	}

	if language != "en" && language != "am" {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	guidanceCard := &entities.GuidanceCard{
		TopicKey:  topicKey,
		Language:  language,
		SelfCare:  blocks.SelfCare,
		OTCCategories: blocks.OTCCategories,
		SeekCareIf: blocks.SeekCareIf,
		Disclaimer: blocks.Disclaimer,
		IsOffline: false,
	}

	return guidanceCard, nil
}
