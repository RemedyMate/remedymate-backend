package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationRepositoryImpl struct {
	collection *mongo.Collection
}

// NewConversationRepository creates a new conversation repository
func NewConversationRepository(collection *mongo.Collection) interfaces.ConversationRepository {
	// Create indexes for better performance
	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetExpireAfterSeconds(24 * 60 * 60), // Expire after 24 hours
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
	})

	if err != nil {
		// Log error but don't fail - indexes might already exist
		log.Printf("Failed to create conversation indexes: %v", err)
	}

	return &ConversationRepositoryImpl{
		collection: collection,
	}
}

// CreateConversation creates a new conversation
func (cr *ConversationRepositoryImpl) CreateConversation(ctx context.Context, conversation *entities.Conversation) error {
	conversation.CreatedAt = time.Now()
	conversation.UpdatedAt = time.Now()
	conversation.Status = entities.ConversationStatusActive
	conversation.CurrentStep = 1

	// Generate ObjectID if not provided
	if conversation.ID == "" {
		conversation.ID = primitive.NewObjectID().Hex()
	}

	_, err := cr.collection.InsertOne(ctx, conversation)
	if err != nil {
		return err
	}

	return nil
}

func (cr *ConversationRepositoryImpl) GetOfflineHealthTopics(ctx context.Context) ([]entities.HealthTopic, error) {
	// Use the correct collection for health topics
	healthTopicsCollection := cr.collection.Database().Collection("health_topics")

	filter := bson.M{"status": "active"}
	cursor, err := healthTopicsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var healthTopics []entities.HealthTopic
	if err := cursor.All(ctx, &healthTopics); err != nil {
		return nil, err
	}

	return healthTopics, nil
}

// GetConversation retrieves a conversation by ID
func (cr *ConversationRepositoryImpl) GetConversation(ctx context.Context, conversationID string) (*entities.Conversation, error) {
	var conversation entities.Conversation

	filter := bson.M{"_id": conversationID}
	err := cr.collection.FindOne(ctx, filter).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("conversation not found: %s", conversationID)
		}
		return nil, err
	}

	return &conversation, nil
}

// UpdateConversation updates an existing conversation
func (cr *ConversationRepositoryImpl) UpdateConversation(ctx context.Context, conversation *entities.Conversation) error {
	conversation.UpdatedAt = time.Now()

	filter := bson.M{"_id": conversation.ID}
	update := bson.M{"$set": conversation}

	_, err := cr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// AddAnswer adds an answer to a conversation
func (cr *ConversationRepositoryImpl) AddAnswer(ctx context.Context, conversationID string, answer entities.Answer) error {
	answer.AnsweredAt = time.Now()

	filter := bson.M{"_id": conversationID}
	update := bson.M{
		"$push": bson.M{"answers": answer},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err := cr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// UpdateConversationStatus updates the status of a conversation
func (cr *ConversationRepositoryImpl) UpdateConversationStatus(ctx context.Context, conversationID string, status entities.ConversationStatus) error {
	filter := bson.M{"_id": conversationID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	if status == entities.ConversationStatusComplete {
		now := time.Now()
		update["$set"].(bson.M)["completed_at"] = now
	}

	_, err := cr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// SetFinalReport sets the final health report for a conversation
func (cr *ConversationRepositoryImpl) SetFinalReport(ctx context.Context, conversationID string, report *entities.HealthReport) error {
	report.GeneratedAt = time.Now()

	filter := bson.M{"_id": conversationID}
	update := bson.M{
		"$set": bson.M{
			"final_report": report,
			"updated_at":   time.Now(),
		},
	}

	_, err := cr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// DeleteExpiredConversations deletes conversations that have expired
func (cr *ConversationRepositoryImpl) DeleteExpiredConversations(ctx context.Context, maxAgeHours int) error {
	cutoffTime := time.Now().Add(-time.Duration(maxAgeHours) * time.Hour)

	filter := bson.M{
		"created_at": bson.M{"$lt": cutoffTime},
		"status":     bson.M{"$in": []entities.ConversationStatus{entities.ConversationStatusActive}},
	}

	_, err := cr.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
