package dto

// OAuthCallbackDTO represents the callback data from OAuth providers
type OAuthCallbackDTO struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state"` // Optional CSRF protection
}

// OAuthResponseDTO represents the response after successful OAuth authentication
type OAuthResponseDTO struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
	Message      string      `json:"message"`
}

// OAuthURLResponseDTO represents the OAuth authorization URL
type OAuthURLResponseDTO struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"` // For CSRF protection
}

// OAuthErrorDTO represents OAuth error responses
type OAuthErrorDTO struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
	Provider    string `json:"provider,omitempty"`
}
