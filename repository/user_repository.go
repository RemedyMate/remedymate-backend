package repository

import (
	"context"
	"fmt"

	"github.com/RemedyMate/remedymate-backend/domain/entities"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
	"github.com/RemedyMate/remedymate-backend/infrastructure/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	UserCollection *mongo.Collection
}

func NewUserRepository() interfaces.IUserRepository{
	userColl := database.Client.Database("remedymate").Collection("users")

	indexModels := []mongo.IndexModel{
        {Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true)},
        {Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)},
        {Keys: bson.M{"email": 1, "isVerified": 1}},
    }	

	_, err := userColl.Indexes().CreateMany(context.Background(), indexModels)
    if err != nil {
        fmt.Println("Error creating indexes:", err)
    } 

	return &UserRepository{UserCollection: userColl}
}

func (r *UserRepository) InsertUser(ctx context.Context, user entities.User) error{
	_, err := r.UserCollection.InsertOne(ctx, user)
    return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string)(*entities.User, error){
	var user entities.User
	err := r.UserCollection.FindOne(ctx, bson.M{"email":email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}