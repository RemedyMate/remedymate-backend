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

type RedFlagRepositoryImpl struct {
	coll *mongo.Collection
}

func NewRedFlagRepository() interfaces.RedFlagRepository {
	c := database.Client.Database("remedymate").Collection("redflags")
	_, _ = c.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "language", Value: 1}, {Key: "level", Value: 1}}},
		{Keys: bson.D{{Key: "isDeleted", Value: 1}}},
		{Keys: bson.D{{Key: "description", Value: "text"}}},
	})
	return &RedFlagRepositoryImpl{coll: c}
}

func (r *RedFlagRepositoryImpl) List(ctx context.Context) ([]entities.RedFlag, error) {
	filter := bson.M{"isDeleted": bson.M{"$ne": true}}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// Initialize empty slice to ensure JSON marshals as [] instead of null
	out := make([]entities.RedFlag, 0)

	for cur.Next(ctx) {
		var rf entities.RedFlag
		if err := cur.Decode(&rf); err != nil {
			return nil, err
		}
		out = append(out, rf)
	}
	return out, cur.Err()
}

func (r *RedFlagRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.RedFlag, error) {
	var rf entities.RedFlag
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "isDeleted": bson.M{"$ne": true}}).Decode(&rf)
	if err != nil {
		return nil, err
	}
	return &rf, nil
}

func (r *RedFlagRepositoryImpl) Create(ctx context.Context, rf *entities.RedFlag) error {
	rf.ID = primitive.NewObjectID().Hex()
	rf.CreatedAt = time.Now()
	rf.UpdatedAt = rf.CreatedAt
	_, err := r.coll.InsertOne(ctx, rf)
	return err
}

func (r *RedFlagRepositoryImpl) Update(ctx context.Context, rf *entities.RedFlag) error {
	rf.UpdatedAt = time.Now()
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": rf.ID}, bson.M{"$set": rf})
	return err
}

func (r *RedFlagRepositoryImpl) SoftDelete(ctx context.Context, id string, deletedBy string) error {
	now := time.Now()
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"isDeleted": true, "deletedAt": now, "deletedBy": deletedBy}})
	return err
}
