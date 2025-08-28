package interfaces

import (
	"context"

	"remedymate-backend/domain/entities"
)

// IOAuthRepository defines the contract for OAuth-related data operations
type IOAuthRepository interface {
	// FindByOAuthProvider finds a user by their OAuth provider and ID
	FindByOAuthProvider(ctx context.Context, provider, providerID string) (*entities.User, error)

	// FindByEmail finds a user by email address
	FindByEmail(ctx context.Context, email string) (*entities.User, error)

	// InsertUser creates a new user in the database
	InsertUser(ctx context.Context, user entities.User) error

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user entities.User) error

	// UpsertOAuthProvider adds or updates an OAuth provider for a user
	UpsertOAuthProvider(ctx context.Context, userID string, oauthProvider entities.OAuthProvider) error

	// FindByID finds a user by their database ID
	FindByID(ctx context.Context, userID string) (*entities.User, error)
}
