package repository

import (
	"context"
	AppError "remedymate-backend/domain/AppError"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshTokenRepository struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository() interfaces.IRefreshTokenRepository {
	collection := database.Client.Database("remedymate").Collection("refresh_tokens")

	return &RefreshTokenRepository{
		collection: collection,
	}
}

func (r *RefreshTokenRepository) StoreRefreshToken(ctx context.Context, refreshToken *entities.RefreshToken) error {
	_, err := r.collection.InsertOne(ctx, refreshToken)
	if err != nil {
		return AppError.ErrInternalServer
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, tokenId string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": tokenId})
	if err != nil {
		return AppError.ErrInternalServer
	}
	if result.DeletedCount == 0 {
		return AppError.ErrRefreshTokenNotFound
	}

	return nil
}
