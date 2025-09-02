package repository

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"remedymate-backend/domain/entities"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // OAuthRepository implements IOAuthRepository interface
// type OAuthRepository struct {
// 	collection *mongo.Collection
// }

// // NewOAuthRepository creates a new OAuth repository instance
// func NewOAuthRepository(collection *mongo.Collection) *OAuthRepository {
// 	return &OAuthRepository{
// 		collection: collection,
// 	}
// }

// // FindByOAuthProvider finds a user by their OAuth provider and ID
// func (r *OAuthRepository) FindByOAuthProvider(ctx context.Context, provider, providerID string) (*entities.User, error) {
// 	filter := bson.M{
// 		"oauthProviders": bson.M{
// 			"$elemMatch": bson.M{
// 				"provider": provider,
// 				"id":       providerID,
// 			},
// 		},
// 	}

// 	var user entities.User
// 	err := r.collection.FindOne(ctx, filter).Decode(&user)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, nil // User not found
// 		}
// 		return nil, err
// 	}

// 	return &user, nil
// }

// // FindByEmail finds a user by email address
// func (r *OAuthRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
// 	filter := bson.M{"email": email}

// 	var user entities.User
// 	err := r.collection.FindOne(ctx, filter).Decode(&user)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	return &user, nil
// }

// // InsertUser creates a new user in the database
// func (r *OAuthRepository) InsertUser(ctx context.Context, user entities.User) error {
// 	log.Printf("💾 Inserting new user into database: %s (%s)", user.Username, user.Email)

// 	if user.ID == "" {
// 		user.ID = primitive.NewObjectID().Hex()
// 		log.Printf("🆔 Generated new ObjectID: %s", user.ID)
// 	}

// 	user.CreatedAt = time.Now()
// 	user.UpdatedAt = time.Now()

// 	_, err := r.collection.InsertOne(ctx, user)
// 	if err != nil {
// 		log.Printf("❌ Failed to insert user: %v", err)
// 		return err
// 	}

// 	log.Printf("✅ Successfully inserted user into database: %s (%s)", user.Username, user.ID)
// 	return nil
// }

// // UpdateUser updates an existing user
// func (r *OAuthRepository) UpdateUser(ctx context.Context, user entities.User) error {
// 	log.Printf("🔄 Updating user in database: %s (%s)", user.Username, user.ID)

// 	user.UpdatedAt = time.Now()

// 	filter := bson.M{"_id": user.ID}
// 	update := bson.M{"$set": user}

// 	_, err := r.collection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		log.Printf("❌ Failed to update user: %v", err)
// 		return err
// 	}

// 	log.Printf("✅ Successfully updated user in database: %s (%s)", user.Username, user.ID)
// 	return nil
// }

// // UpsertOAuthProvider adds or updates an OAuth provider for a user
// func (r *OAuthRepository) UpsertOAuthProvider(ctx context.Context, userID string, oauthProvider entities.OAuthProvider) error {
// 	filter := bson.M{"_id": userID}

// 	// Check if provider already exists
// 	update := bson.M{
// 		"$addToSet": bson.M{
// 			"oauthProviders": oauthProvider,
// 		},
// 		"$set": bson.M{
// 			"updatedAt": time.Now(),
// 		},
// 	}

// 	_, err := r.collection.UpdateOne(ctx, filter, update)
// 	return err
// }

// // FindByID finds a user by their database ID
// func (r *OAuthRepository) FindByID(ctx context.Context, userID string) (*entities.User, error) {
// 	filter := bson.M{"_id": userID}

// 	var user entities.User
// 	err := r.collection.FindOne(ctx, filter).Decode(&user)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	return &user, nil
// }
