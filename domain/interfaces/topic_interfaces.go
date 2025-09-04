package interfaces

import (
	"context"
	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
)

type TopicRepository interface {
	// CreateTopic inserts a new topic into the database.
	CreateTopic(ctx context.Context, topic *entities.Topic) error

	// GetTopicByKey retrieves a single topic by its unique topic_key.
	GetTopicByKey(ctx context.Context, topicKey string) (*entities.Topic, error)

	// ListAllTopics retrieves all topics with pagination and filtering.
	ListAllTopics(ctx context.Context, params dto.TopicListQueryParams) ([]*entities.Topic, int64, error)

	// UpdateTopic updates an existing topic by its topic_key.
	UpdateTopic(ctx context.Context, topicKey string, update *entities.Topic) error

	// DeleteTopic performs a soft delete on a topic by changing its status.
	DeleteTopic(ctx context.Context, topicKey string, deletedByUserID string) error

	// CheckTopicExists checks if a topic with the given topic_key exists.
	CheckTopicExists(ctx context.Context, topicKey string) (bool, error)
}

type TopicUsecase interface {
	// CreateTopic handles the creation of a new topic, including validation and setting audit fields.
	CreateTopic(ctx context.Context, request dto.TopicCreateRequest) (*entities.Topic, error)

	// GetTopic retrieves a single topic, potentially including its content.
	GetTopicByKey(ctx context.Context, topicKey string) (*entities.Topic, error)

	// ListAllTopics retrieves a paginated, filtered, and sorted list of topics.
	ListAllTopics(ctx context.Context, params dto.TopicListQueryParams) (*dto.PaginatedTopicsResult, error)

	// UpdateTopic handles updating an existing topic, including validation and setting audit fields.
	UpdateTopic(ctx context.Context, topicKey string, request dto.TopicUpdateRequest) (*entities.Topic, error)

	// SoftDeleteTopic performs a soft delete, marking a topic as inactive but retaining its data.
	SoftDeleteTopic(ctx context.Context, topicKey string) error
}
