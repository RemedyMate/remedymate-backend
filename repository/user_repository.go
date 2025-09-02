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

	_, err := userColl.Indexes().CreateMany(context.Background(), userIndexModels)
	if err != nil {
		fmt.Println("Error creating indexes:", err)
	}

	_, err = userColl.Indexes().CreateMany(context.Background(), userStatusIndexModels)
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
	// session, err := database.Client.StartSession()
	// if err != nil {
	// 	return fmt.Errorf("failed to start session: %w", err)
	// }
	// defer session.EndSession(ctx)

	_, err := r.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return AppError.ErrInternalServer
	}
	_, err = r.UserStatusCollection.InsertOne(ctx, userStatus)
	if err != nil {
		return AppError.ErrInternalServer
	}

	// callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
	// 	return nil, nil
	// }
	// _, err = session.WithTransaction(ctx, callback)
	// if err != nil {
	// 	return err
	// }
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
	err := r.UserStatusCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&userStatus)
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
