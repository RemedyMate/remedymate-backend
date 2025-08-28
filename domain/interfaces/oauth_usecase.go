package interfaces

import (
	"context"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
)

// IOAuthUsecase defines the contract for OAuth business logic
type IOAuthUsecase interface {
	// GetAuthURL generates the OAuth authorization URL for a specific provider
	GetAuthURL(ctx context.Context, provider string) (*dto.OAuthURLResponseDTO, error)

	// HandleCallback processes the OAuth callback and authenticates the user
	HandleCallback(ctx context.Context, provider string, callback dto.OAuthCallbackDTO) (*dto.OAuthResponseDTO, error)

	// RefreshToken refreshes an expired access token
	RefreshToken(ctx context.Context, refreshToken string) (*dto.OAuthResponseDTO, error)

	// ValidateToken validates a JWT token and returns user information
	ValidateToken(ctx context.Context, token string) (*entities.User, error)
}
