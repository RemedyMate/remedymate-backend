package repository

import (
	"context"
	AppError "remedymate-backend/domain/AppError"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshTokenRepository struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository() interfaces.IRefreshTokenRepository {
	// TODO: create an index to delete expired tokens
	collection := database.Client.Database("remedymate").Collection("refresh_tokens")
	return &RefreshTokenRepository{
		collection: collection,
	}
}

func (r *RefreshTokenRepository) StoreRefreshToken(ctx context.Context, refreshToken *entities.RefreshToken) error {
	id := primitive.NewObjectID().Hex()
	refreshToken.ID = id
	_, err := r.collection.InsertOne(ctx, refreshToken)
	if err != nil {
		return AppError.ErrInternalServer
	}
	return nil
}
