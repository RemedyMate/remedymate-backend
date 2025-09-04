package usecase

import (
	"context"
	"time"

	"remedymate-backend/domain/AppError"
	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicUsecase struct {
	topicRepository interfaces.TopicRepository
}

func NewTopicUsecase(topicRepo interfaces.TopicRepository) *TopicUsecase {
	return &TopicUsecase{
		topicRepository: topicRepo,
	}
}

func (tu *TopicUsecase) CreateTopic(ctx context.Context, request dto.TopicCreateRequest) (*entities.Topic, error) {
	// Validate required fields
	if request.TopicKey == "" || request.NameEN == "" || request.NameAM == "" {
		return nil, AppError.ErrInvalidInput
	}

	createdByUserID, err := extractUserID(ctx)
	if err != nil {
		return nil, err
	}
	// ensure createdBy is a valid ObjectID
	createdByOID, err := primitive.ObjectIDFromHex(createdByUserID)
	if err != nil {
		return nil, AppError.ErrInvalidInput
	}

	// Prevent duplicate topic_key
	if existing, _ := tu.topicRepository.GetTopicByKey(ctx, request.TopicKey); existing != nil {
		return nil, AppError.ErrTopicAlreadyExists
	}

	now := time.Now()
	topic := &entities.Topic{
		TopicKey:      request.TopicKey,
		NameEN:        request.NameEN,
		NameAM:        request.NameAM,
		DescriptionEN: request.DescriptionEN,
		DescriptionAM: request.DescriptionAM,
		Status:        entities.TopicStatusActive,
		Translations:  request.Translations,
		Version:       1,
		CreatedAt:     now,
		UpdatedAt:     now,
		CreatedBy:     createdByOID,
		UpdatedBy:     createdByOID,
	}

	if err := tu.topicRepository.CreateTopic(ctx, topic); err != nil {
		// map repository duplicate-key or other domain errors if repository returns them
		return nil, err
	}

	// Return the created topic (fresh from DB)
	created, err := tu.topicRepository.GetTopicByKey(ctx, request.TopicKey)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (tu *TopicUsecase) GetTopicByKey(ctx context.Context, key string) (*entities.Topic, error) {
	if key == "" {
		return nil, AppError.ErrInvalidInput
	}
	topic, err := tu.topicRepository.GetTopicByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (tu *TopicUsecase) ListAllTopics(ctx context.Context, params dto.TopicListQueryParams) (*dto.PaginatedTopicsResult, error) {
	// Validate pagination parameters
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	topics, total, err := tu.topicRepository.ListAllTopics(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert []*entities.Topic to []entities.Topic
	topicsVal := make([]entities.Topic, len(topics))
	for i, t := range topics {
		if t != nil {
			topicsVal[i] = *t
		}
	}

	return &dto.PaginatedTopicsResult{
		TotalCount: total,
		Topics:     topicsVal,
	}, nil
}

// extractUserID extracts the user ID from context, expecting "userID" key.
func extractUserID(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", AppError.ErrUserNotAuthenticated
	}
	if v := ctx.Value("userID"); v != nil {
		if s, ok := v.(string); ok {
			if s != "" {
				return s, nil
			}
		}
	}
	return "", AppError.ErrUserNotAuthenticated
}

func (tu *TopicUsecase) UpdateTopic(ctx context.Context, topicKey string, request dto.TopicUpdateRequest) (*entities.Topic, error) {
	// Validate required fields
	if topicKey == "" {
		return nil, AppError.ErrInvalidInput
	}
	updatedByUserID, err := extractUserID(ctx)
	if err != nil {
		return nil, err
	}
	updatedByOID, err := primitive.ObjectIDFromHex(updatedByUserID)
	if err != nil {
		return nil, AppError.ErrInvalidInput
	}

	// Ensure topic exists
	existing, err := tu.topicRepository.GetTopicByKey(ctx, topicKey)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, AppError.ErrTopicNotFound
	}

	// Apply allowed updates and preserve CreatedAt/CreatedBy
	existing.NameEN = request.NameEN
	existing.NameAM = request.NameAM
	existing.DescriptionEN = request.DescriptionEN
	existing.DescriptionAM = request.DescriptionAM
	if request.Translations != nil {
		existing.Translations = request.Translations
	}
	existing.UpdatedAt = time.Now()
	existing.UpdatedBy = updatedByOID
	existing.Version = existing.Version + 1

	if err := tu.topicRepository.UpdateTopic(ctx, topicKey, existing); err != nil {
		return nil, err
	}

	updated, err := tu.topicRepository.GetTopicByKey(ctx, topicKey)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (tu *TopicUsecase) SoftDeleteTopic(ctx context.Context, topicKey string) error {
	deletedByUserID, err := extractUserID(ctx)
	if err != nil {
		return err
	}
	// verify existence
	_, err = tu.topicRepository.GetTopicByKey(ctx, topicKey)
	if err != nil {
		return err
	}
	// repository should set status to deleted and record deletedBy if supported
	if err := tu.topicRepository.DeleteTopic(ctx, topicKey, deletedByUserID); err != nil {
		return err
	}
	return nil
}
