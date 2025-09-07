package repository

import (
	"context"
	"fmt"
	"time"

	"remedymate-backend/domain/AppError"
	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TopicRepository handles topic persistence.
type TopicRepository struct {
	TopicCollection *mongo.Collection
}

// NewTopicRepository creates a repository using the provided mongo.Database.
func NewTopicRepository() (*TopicRepository, error) {
	topicColl := database.Client.Database("remedymate").Collection("topics")

	// ensure text index on relevant fields
	idxModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name_en", Value: "text"},
			{Key: "name_am", Value: "text"},
			{Key: "topic_key", Value: "text"},
			{Key: "description_en", Value: "text"},
			{Key: "description_am", Value: "text"},
		},
	}
	if _, err := topicColl.Indexes().CreateOne(context.Background(), idxModel); err != nil {
		return nil, fmt.Errorf("failed to create topic_key index: %w", err)
	}

	return &TopicRepository{
		TopicCollection: topicColl,
	}, nil
}

// CreateTopic inserts a new topic, setting defaults for timestamps/status and mapping duplicate-key errors.
func (tr *TopicRepository) CreateTopic(ctx context.Context, topic *entities.Topic) error {
	now := time.Now()
	if topic.CreatedAt.IsZero() {
		topic.CreatedAt = now
	}
	topic.UpdatedAt = now
	if topic.Status == "" {
		topic.Status = entities.TopicStatusActive
	}

	_, err := tr.TopicCollection.InsertOne(ctx, topic)
	if err != nil {
		// map duplicate key to domain error
		if mongo.IsDuplicateKeyError(err) {
			return AppError.ErrTopicAlreadyExists
		}
		return fmt.Errorf("failed to insert topic: %w", err)
	}
	return nil
}

// GetTopicByKey returns the topic for the given key. By default excludes soft-deleted topics.
func (tr *TopicRepository) GetTopicByKey(ctx context.Context, key string) (*entities.Topic, error) {
	var topic entities.Topic
	filter := bson.M{
		"topic_key": key,
		"status":    bson.M{"$ne": entities.TopicStatusDeleted},
	}
	err := tr.TopicCollection.FindOne(ctx, filter).Decode(&topic)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, AppError.ErrTopicNotFound
		}
		return nil, fmt.Errorf("failed to query topic by key: %w", err)
	}
	return &topic, nil
}

// UpdateTopic applies a partial update to the topic identified by topicKey.
func (tr *TopicRepository) UpdateTopic(ctx context.Context, topicKey string, update *entities.Topic) error {
	updateFields := bson.M{
		"updated_at": time.Now(),
		"updated_by": update.UpdatedBy,
	}

	updateFields["name_en"] = update.NameEN
	updateFields["name_am"] = update.NameAM
	updateFields["description_en"] = update.DescriptionEN
	updateFields["description_am"] = update.DescriptionAM
	updateFields["status"] = update.Status

	if update.Translations != nil {
		updateFields["translations"] = update.Translations
	}
	// increment version if provided/expected
	if update.Version > 0 {
		updateFields["version"] = update.Version
	}

	updateDoc := bson.M{"$set": updateFields}

	result, err := tr.TopicCollection.UpdateOne(ctx, bson.M{"topic_key": topicKey}, updateDoc)
	if err != nil {
		return fmt.Errorf("failed to update topic %s: %w", topicKey, err)
	}
	// MatchedCount == 0 means no matching document found
	if result.MatchedCount == 0 {
		return AppError.ErrTopicNotFound
	}
	return nil
}

// DeleteTopic performs a soft-delete (sets status=deleted). It's idempotent.
func (tr *TopicRepository) DeleteTopic(ctx context.Context, topicKey string, deletedByUserID string) error {
	filter := bson.M{
		"topic_key": topicKey,
		"status":    bson.M{"$ne": entities.TopicStatusDeleted},
	}
	update := bson.M{
		"$set": bson.M{
			"status":     entities.TopicStatusDeleted,
			"updated_at": time.Now(),
			"updated_by": deletedByUserID,
		},
	}
	result, err := tr.TopicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to soft-delete topic %s: %w", topicKey, err)
	}
	if result.MatchedCount == 0 {
		// Either not found or already deleted
		return AppError.ErrTopicNotFound
	}
	return nil
}

// ListAllTopics returns topics matching params and the total count for pagination.
func (tr *TopicRepository) ListAllTopics(ctx context.Context, params dto.TopicListQueryParams) ([]*entities.Topic, int64, error) {
	var topics []*entities.Topic

	filter := bson.M{
		"status": bson.M{"$ne": entities.TopicStatusDeleted},
	}
	if params.Search != "" {
		regex := primitive.Regex{Pattern: params.Search, Options: "i"}
		filter["$or"] = []bson.M{
			{"name_en": regex},
			{"name_am": regex},
			{"topic_key": regex},
			{"description_en": regex},
			{"description_am": regex},
		}
	}
	findOptions := options.Find()
	// Sorting: default by name_en ascending
	if params.SortBy != "" {
		order := 1
		if params.Order == "desc" {
			order = -1
		}
		findOptions.SetSort(bson.D{{Key: params.SortBy, Value: order}})
	} else {
		findOptions.SetSort(bson.D{{Key: "name_en", Value: 1}})
	}

	// Pagination
	limit := int64(params.Limit)
	if limit <= 0 {
		limit = 20
	}
	page := int64(params.Page)
	if page <= 0 {
		page = 1
	}
	findOptions.SetLimit(limit)
	findOptions.SetSkip((page - 1) * limit)

	cursor, err := tr.TopicCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list topics: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var topic entities.Topic
		if err := cursor.Decode(&topic); err != nil {
			return nil, 0, fmt.Errorf("failed to decode topic: %w", err)
		}
		topics = append(topics, &topic)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}

	// total count for the same filter (use CountDocuments)
	total, err := tr.TopicCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count topics: %w", err)
	}

	return topics, total, nil
}

func (tr *TopicRepository) CheckTopicExists(ctx context.Context, topicKey string) (bool, error) {
	var count int64
	filter := bson.M{"topic_key": topicKey}
	count, err := tr.TopicCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check if topic exists: %w", err)
	}
	return count > 0, nil
}
