package repository

import (
	"context"
	"log"
	"time"

	AppError "remedymate-backend/domain/AppError"
	"remedymate-backend/domain/entities"
	"remedymate-backend/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConsentRepository provides persistence for append-only consent records.
type ConsentRepository struct {
	collection *mongo.Collection
}

// NewConsentRepository constructs a new ConsentRepository using the global DB client.
func NewConsentRepository() *ConsentRepository {
	coll := database.Client.Database(database.GetDatabaseName()).Collection("consent_records")

	// ensure index on userId + createdAt for queries
	ctx := context.Background()
	idxModels := []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}}, Options: options.Index()},
	}
	if _, err := coll.Indexes().CreateMany(ctx, idxModels); err != nil {
		// non-fatal: log and continue
		log.Printf("error creating consent_records indexes: %v", err)
	}

	return &ConsentRepository{collection: coll}
}

// CreateConsent inserts a new consent record for a user (append-only).
func (r *ConsentRepository) CreateConsent(ctx context.Context, userID string, envelope *entities.ConsentEnvelope) error {
	doc := bson.M{
		"_id":       primitive.NewObjectID().Hex(),
		"userId":    userID,
		"consent":   envelope,
		"createdAt": time.Now(),
	}
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		log.Printf("error inserting consent record: %v", err)
		return AppError.ErrInternalServer
	}
	return nil
}

// ListConsentByUser returns consent records for a user ordered by createdAt desc.
func (r *ConsentRepository) ListConsentByUser(ctx context.Context, userID string) ([]*entities.ConsentEnvelope, error) {
	filter := bson.M{"userId": userID}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("error finding consent records for user %s: %v", userID, err)
		return nil, AppError.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var results []*entities.ConsentEnvelope
	for cursor.Next(ctx) {
		var doc struct {
			ID      string                   `bson:"_id"`
			UserID  string                   `bson:"userId"`
			Consent entities.ConsentEnvelope `bson:"consent"`
			Created time.Time                `bson:"createdAt"`
		}
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("error decoding consent record: %v", err)
			continue
		}
		results = append(results, &doc.Consent)
	}
	if err := cursor.Err(); err != nil {
		return nil, AppError.ErrInternalServer
	}
	return results, nil
}

// AppendRevocation appends a revocation/change record for a given consent type.
// This method creates a new consent record entry with Granted=false and optional revocation metadata.
func (r *ConsentRepository) AppendRevocation(ctx context.Context, userID, consentType, reason string) error {
	item := entities.ConsentItem{
		Type:      consentType,
		Granted:   false,
		Version:   "", // or latest
		Timestamp: time.Now(),
		Reason:    reason,
	}
	envelope := &entities.ConsentEnvelope{
		AuditTrail:  []entities.ConsentItem{item},
		LastUpdated: time.Now(),
	}
	return r.CreateConsent(ctx, userID, envelope)
}
