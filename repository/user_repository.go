package repository

import (
	"context"
	"fmt"
	"log"

	AppError "remedymate-backend/domain/AppError"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	UserCollection       *mongo.Collection
	UserStatusCollection *mongo.Collection
}

func NewUserRepository() interfaces.IUserRepository {
	userColl := database.Client.Database("remedymate").Collection("users")
	userStatColl := database.Client.Database("remedymate").Collection("user_status")
	userIndexModels := []mongo.IndexModel{
		{Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)},
	}
	userStatusIndexModels := []mongo.IndexModel{
		{Keys: bson.M{"userId": 1}, Options: options.Index().SetUnique(true)},
	}

	// Ensure no incorrect index exists on users collection for userId
	// TODO: Check the use of this code
	ctx := context.Background()
	cursor, err := userColl.Indexes().List(ctx)
	if err == nil {
		for cursor.Next(ctx) {
			var idxDoc bson.M
			if err := cursor.Decode(&idxDoc); err == nil {
				// Check by name first
				if name, ok := idxDoc["name"].(string); ok && name == "userId_1" {
					_, _ = userColl.Indexes().DropOne(ctx, name)
					continue
				}
				// Check by key contents
				if keyDoc, ok := idxDoc["key"].(bson.M); ok {
					if _, exists := keyDoc["userId"]; exists {
						if name, ok := idxDoc["name"].(string); ok {
							_, _ = userColl.Indexes().DropOne(ctx, name)
						}
					}
				}
			}
		}
		_ = cursor.Close(ctx)
	}

	_, err = userColl.Indexes().CreateMany(ctx, userIndexModels)
	if err != nil {
		fmt.Println("Error creating indexes:", err)
	}

	// Create indexes on the correct user status collection
	_, err = userStatColl.Indexes().CreateMany(ctx, userStatusIndexModels)
	if err != nil {
		fmt.Println("Error creating indexes:", err)
	}

	_, err = userStatColl.Indexes().CreateMany(context.Background(), userStatusIndexModels)
	if err != nil {
		fmt.Println("Error creating indexes:", err)
	}

	return &UserRepository{UserCollection: userColl, UserStatusCollection: userStatColl}
}

func (r *UserRepository) CreateUserWithStatus(ctx context.Context, user *entities.User, userStatus *entities.UserStatus) error {
	user.ID = primitive.NewObjectID().Hex()
	userStatus.ID = primitive.NewObjectID().Hex()
	userStatus.UserID = user.ID

	// TODO: create a transaction to ensure both user and user status are created
	_, err := r.UserCollection.InsertOne(ctx, user)
	if err != nil {

		log.Printf("Error inserting user: %v", err)
		return AppError.ErrInternalServer
	}
	_, err = r.UserStatusCollection.InsertOne(ctx, userStatus)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return AppError.ErrInternalServer
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, AppError.ErrUserNotFound
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	err := r.UserCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, AppError.ErrUserNotFound
	}
	return &user, nil
}

// CheckByRole checks if any user exists with the specified role
func (r *UserRepository) CheckByRole(ctx context.Context, role string) (*bool, error) {
	count, err := r.UserCollection.CountDocuments(ctx, bson.M{"role": role})
	if err != nil {
		return nil, AppError.ErrInternalServer
	}
	exists := count > 0
	return &exists, nil
}

// FindByID finds a user by their database ID
func (r *UserRepository) FindByID(ctx context.Context, userID string) (*entities.User, error) {
	var user entities.User
	err := r.UserCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	_, err := r.UserCollection.UpdateOne(ctx, filter, update)
	return err
}

// SoftDeleteUser marks a user as inactive instead of deleting them
func (r *UserRepository) SoftDeleteUser(ctx context.Context, userID string) error {
	filter := bson.M{"userId": userID}
	update := bson.M{
		"$set": bson.M{
			"isActive": false,
		},
	}

	_, err := r.UserStatusCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *UserRepository) GetUserStatus(ctx context.Context, userID string) (*entities.UserStatus, error) {
	var userStatus entities.UserStatus

	// Look up by userId field (unique), not the document _id
	err := r.UserStatusCollection.FindOne(ctx, bson.M{"userId": userID}).Decode(&userStatus)
	if err != nil {
		return nil, AppError.ErrUserStatusNotFound
	}
	return &userStatus, nil
}

func (r *UserRepository) CreateUserStatus(ctx context.Context, userStatus *entities.UserStatus) error {
	_, err := r.UserStatusCollection.InsertOne(ctx, userStatus)
	if err != nil {
		log.Printf("Error inserting user status: %v", err)
		return AppError.ErrInternalServer
	}
	return nil
}

// UpdateUserStatusFields updates specific fields in user_status for a given userId
func (r *UserRepository) UpdateUserStatusFields(ctx context.Context, userID string, fields map[string]interface{}) error {
	filter := bson.M{"userId": userID}
	update := bson.M{"$set": fields}
	_, err := r.UserStatusCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return AppError.ErrInternalServer
	}
	return nil
}

// ActivateByEmail sets a user's status IsActive=true by email
func (r *UserRepository) ActivateByEmail(ctx context.Context, email string) error {
	// Find user by email to get userID
	var user entities.User
	if err := r.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return AppError.ErrUserNotFound
	}

	filter := bson.M{"userId": user.ID}
	update := bson.M{"$set": bson.M{"isActive": true, "isVerified": true}}
	result, err := r.UserStatusCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return AppError.ErrInternalServer
	}
	if result.MatchedCount == 0 {
		return AppError.ErrUserStatusNotFound
	}
	return nil
}
