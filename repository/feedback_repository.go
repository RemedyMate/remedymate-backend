package repository

import (
	"context"
	"time"

	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FeedbackRepositoryImpl struct {
	coll *mongo.Collection
}

func NewFeedbackRepository() interfaces.FeedbackRepository {
	c := database.Client.Database("remedymate").Collection("feedbacks")
	_, _ = c.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "createdAt", Value: -1}}},
		{Keys: bson.D{{Key: "language", Value: 1}}},
		{Keys: bson.D{{Key: "isDeleted", Value: 1}}},
	})
	return &FeedbackRepositoryImpl{coll: c}
}

func (r *FeedbackRepositoryImpl) List(ctx context.Context, limit, offset int, language string) ([]entities.Feedback, error) {
	filter := bson.M{"isDeleted": bson.M{"$ne": true}}
	if language != "" {
		filter["language"] = language
	}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// Initialize empty slice to ensure JSON marshals as [] instead of null
	out := make([]entities.Feedback, 0)

	for cur.Next(ctx) {
		var f entities.Feedback
		if err := cur.Decode(&f); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, cur.Err()
}

func (r *FeedbackRepositoryImpl) Count(ctx context.Context, language string) (int64, error) {
	filter := bson.M{"isDeleted": bson.M{"$ne": true}}
	if language != "" {
		filter["language"] = language
	}
	return r.coll.CountDocuments(ctx, filter)
}

func (r *FeedbackRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Feedback, error) {
	var f entities.Feedback
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "isDeleted": bson.M{"$ne": true}}).Decode(&f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FeedbackRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	now := time.Now()
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"isDeleted": true, "deletedAt": now}})
	return err
}

func (r *FeedbackRepositoryImpl) Create(ctx context.Context, f *entities.Feedback) error {
	f.ID = primitive.NewObjectID().Hex()
	f.CreatedAt = time.Now()
	_, err := r.coll.InsertOne(ctx, f)
	return err
}
