package repository

import (
	"context"
	"log"
	"time"

	AppError "remedymate-backend/domain/AppError"
	"remedymate-backend/domain/entities"
	"remedymate-backend/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserProfileRepository provides persistence for user onboarding profiles (AppUser).
type UserProfileRepository struct {
	collection *mongo.Collection
}

// NewUserProfileRepository constructs a new UserProfileRepository using the global DB client.
func NewUserProfileRepository() *UserProfileRepository {
	coll := database.Client.Database(database.GetDatabaseName()).Collection("app_users")

	// ensure index on userId (unique)
	ctx := context.Background()
	idxModels := []mongo.IndexModel{
		{Keys: bson.M{"userId": 1}, Options: options.Index().SetUnique(true)},
	}
	if _, err := coll.Indexes().CreateMany(ctx, idxModels); err != nil {
		// non-fatal: log and continue
		log.Printf("error creating app_users indexes: %v", err)
	}

	return &UserProfileRepository{collection: coll}
}

// UpsertDemographics creates or updates the demographics section for a user.
func (r *UserProfileRepository) UpsertDemographics(ctx context.Context, userID string, demographics *entities.Demographics) error {
	filter := bson.M{"userId": userID}
	update := bson.M{
		"$set": bson.M{
			"demographics":                 demographics,
			"status.demographicsCompleted": true,
			"status.lastUpdated":           time.Now(),
			"updatedAt":                    time.Now(),
		},
		"$setOnInsert": bson.M{
			"userId":    userID,
			"createdAt": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("error upserting demographics for user %s: %v", userID, err)
		return AppError.ErrInternalServer
	}
	return nil
}

// UpsertMedicalHistory creates or updates the medical history section for a user.
func (r *UserProfileRepository) UpsertMedicalHistory(ctx context.Context, userID string, history *entities.MedicalHistory) error {
	filter := bson.M{"userId": userID}
	update := bson.M{
		"$set": bson.M{
			"medical_history":                history,
			"status.medicalHistoryCompleted": true,
			"status.lastUpdated":             time.Now(),
			"updatedAt":                      time.Now(),
		},
		"$setOnInsert": bson.M{
			"userId":    userID,
			"createdAt": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("error upserting medical history for user %s: %v", userID, err)
		return AppError.ErrInternalServer
	}
	return nil
}

// UpsertLifestyle creates or updates the lifestyle section for a user.
func (r *UserProfileRepository) UpsertLifestyle(ctx context.Context, userID string, lifestyle *entities.Lifestyle) error {
	filter := bson.M{"userId": userID}
	update := bson.M{
		"$set": bson.M{
			"lifestyle":                 lifestyle,
			"status.lifestyleCompleted": true,
			"status.lastUpdated":        time.Now(),
			"updatedAt":                 time.Now(),
		},
		"$setOnInsert": bson.M{
			"userId":    userID,
			"createdAt": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("error upserting lifestyle for user %s: %v", userID, err)
		return AppError.ErrInternalServer
	}
	return nil
}

// UpsertConsent creates or updates the consent envelope for a user.
func (r *UserProfileRepository) UpsertConsent(ctx context.Context, userID string, consent *entities.ConsentEnvelope) error {
	filter := bson.M{"userId": userID}
	update := bson.M{
		"$set": bson.M{
			"consent":                 consent,
			"status.consentCompleted": true,
			"status.lastUpdated":      time.Now(),
			"updatedAt":               time.Now(),
		},
		"$setOnInsert": bson.M{
			"userId":    userID,
			"createdAt": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("error upserting consent for user %s: %v", userID, err)
		return AppError.ErrInternalServer
	}
	return nil
}

// GetUserProfile retrieves the full AppUser profile for a user.
func (r *UserProfileRepository) GetUserProfile(ctx context.Context, userID string) (*entities.AppUser, error) {
	var appUser entities.AppUser
	err := r.collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&appUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, AppError.ErrUserNotFound
		}
		log.Printf("error finding user profile for %s: %v", userID, err)
		return nil, AppError.ErrInternalServer
	}
	return &appUser, nil
}
