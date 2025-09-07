package repository

import (
	"context"
	"log"
	"time"

	AppError "remedymate-backend/domain/AppError"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ActivationTokenRepository struct {
	collection *mongo.Collection
}

func NewActivationTokenRepository() interfaces.IActivationTokenRepository {
	return &ActivationTokenRepository{collection: database.Client.Database("remedymate").Collection("activation_tokens")}
}

func (r *ActivationTokenRepository) Create(ctx context.Context, token *entities.ActivationToken) error {
	token.ID = primitive.NewObjectID().Hex()
	_, err := r.collection.InsertOne(ctx, token)
	if err != nil {
		return AppError.ErrInternalServer
	}
	return nil
}

func (r *ActivationTokenRepository) FindValidByToken(ctx context.Context, token string) (*entities.ActivationToken, error) {
	var at entities.ActivationToken
	err := r.collection.FindOne(ctx, bson.M{"token": token, "usedAt": bson.M{"$exists": false}, "expiresAt": bson.M{"$gt": time.Now()}}).Decode(&at)
	if err != nil {
		return nil, AppError.ErrInvalidToken
	}
	return &at, nil
}

func (r *ActivationTokenRepository) MarkUsed(ctx context.Context, id string) error {
	now := time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"usedAt": now}})
	if err != nil {
		return AppError.ErrInternalServer
	}
	return nil
}

func (r *ActivationTokenRepository) FindValidActivationTokenByEmail(ctx context.Context, email string) (*entities.ActivationToken, error) {
	var at entities.ActivationToken
	err := r.collection.FindOne(ctx, bson.M{"email": email, "usedAt": bson.M{"$exists": false}, "expiresAt": bson.M{"$gt": time.Now()}}).Decode(&at)
	if err != nil {
		log.Println("Error finding activation token by email:", err)
		return nil, AppError.ErrInvalidActivationToken
	}
	return &at, nil
}
