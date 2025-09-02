package auth

// import (
// 	"context"
// 	"crypto/rand"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"time"

// 	"remedymate-backend/config"
// 	"remedymate-backend/domain/entities"

// 	"golang.org/x/oauth2"
// )

// // OAuthService handles OAuth provider interactions
// type OAuthService struct {
// 	config     *config.OAuthConfig
// 	jwtService *JWTService
// }

// // NewOAuthService creates a new OAuth service instance
// func NewOAuthService(config *config.OAuthConfig, jwtService *JWTService) *OAuthService {
// 	return &OAuthService{
// 		config:     config,
// 		jwtService: jwtService,
// 	}
// }

// // GenerateState generates a random state parameter for CSRF protection
// func (o *OAuthService) GenerateState() (string, error) {
// 	bytes := make([]byte, 32)
// 	if _, err := rand.Read(bytes); err != nil {
// 		return "", err
// 	}
// 	return hex.EncodeToString(bytes), nil
// }

// // GetAuthURL generates the OAuth authorization URL for a specific provider
// func (o *OAuthService) GetAuthURL(provider string) (string, string, error) {
// 	var config *oauth2.Config
// 	var state string
// 	var err error

// 	switch provider {
// 	case "google":
// 		config = o.config.Google
// 	case "facebook":
// 		config = o.config.Facebook
// 	default:
// 		return "", "", fmt.Errorf("unsupported provider: %s", provider)
// 	}

// 	// Generate state for CSRF protection
// 	state, err = o.GenerateState()
// 	if err != nil {
// 		return "", "", fmt.Errorf("failed to generate state: %w", err)
// 	}

// 	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
// 	return authURL, state, nil
// }

// // ExchangeCodeForToken exchanges an authorization code for an access token
// func (o *OAuthService) ExchangeCodeForToken(ctx context.Context, provider, code string) (*oauth2.Token, error) {
// 	var config *oauth2.Config

// 	switch provider {
// 	case "google":
// 		config = o.config.Google
// 	case "facebook":
// 		config = o.config.Facebook
// 	default:
// 		return nil, fmt.Errorf("unsupported provider: %s", provider)
// 	}

// 	return config.Exchange(ctx, code)
// }

// // GetUserInfo retrieves user information from the OAuth provider
// func (o *OAuthService) GetUserInfo(ctx context.Context, provider string, token *oauth2.Token) (*entities.User, error) {
// 	log.Printf("🔄 Getting user info from %s provider...", provider)

// 	var userInfo map[string]interface{}
// 	var err error

// 	switch provider {
// 	case "google":
// 		userInfo, err = o.getGoogleUserInfo(ctx, token.AccessToken)
// 	case "facebook":
// 		userInfo, err = o.getFacebookUserInfo(ctx, token.AccessToken)
// 	default:
// 		return nil, fmt.Errorf("unsupported provider: %s", provider)
// 	}

// 	if err != nil {
// 		log.Printf("❌ Failed to get user info from %s: %v", provider, err)
// 		return nil, err
// 	}

// 	log.Printf("✅ Successfully retrieved user info from %s: %v", provider, userInfo)
// 	user := o.mapToUser(userInfo, provider)
// 	log.Printf("👤 Mapped user info: %s (%s)", user.Username, user.Email)

// 	return user, nil
// }

// // getGoogleUserInfo retrieves user information from Google
// func (o *OAuthService) getGoogleUserInfo(ctx context.Context, accessToken string) (map[string]interface{}, error) {
// 	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Authorization", "Bearer "+accessToken)

// 	client := &http.Client{Timeout: 10 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("Google API returned status: %d", resp.StatusCode)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var userInfo map[string]interface{}
// 	if err := json.Unmarshal(body, &userInfo); err != nil {
// 		return nil, err
// 	}

// 	return userInfo, nil
// }

// // getFacebookUserInfo retrieves user information from Facebook
// func (o *OAuthService) getFacebookUserInfo(ctx context.Context, accessToken string) (map[string]interface{}, error) {
// 	req, err := http.NewRequestWithContext(ctx, "GET", "https://graph.facebook.com/me?fields=id,name,email,first_name,last_name", nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	q := req.URL.Query()
// 	q.Add("access_token", accessToken)
// 	req.URL.RawQuery = q.Encode()

// 	client := &http.Client{Timeout: 10 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("Facebook API returned status: %d", resp.StatusCode)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var userInfo map[string]interface{}
// 	if err := json.Unmarshal(body, &userInfo); err != nil {
// 		return nil, err
// 	}

// 	return userInfo, nil
// }

// // mapToUser converts OAuth provider user info to our User entity
// func (o *OAuthService) mapToUser(userInfo map[string]interface{}, provider string) *entities.User {
// 	log.Printf("🔄 Mapping %s user info to User entity...", provider)

// 	now := time.Now()

// 	user := &entities.User{
// 		IsVerified:    true,  // OAuth users are pre-verified
// 		IsProfileFull: false, // They may need to complete profile
// 		IsActive:      true,
// 		CreatedAt:     now,
// 		UpdatedAt:     now,
// 		LastLogin:     now,
// 		OAuthProviders: []entities.OAuthProvider{
// 			{
// 				Provider: provider,
// 				ID:       fmt.Sprintf("%v", userInfo["id"]),
// 			},
// 		},
// 	}

// 	// Map common fields
// 	if email, ok := userInfo["email"].(string); ok {
// 		user.Email = email
// 		log.Printf("👤 Mapped email: %s", email)
// 	}

// 	// Provider-specific mapping
// 	switch provider {
// 	case "google":
// 		if givenName, ok := userInfo["given_name"].(string); ok {
// 			user.PersonalInfo.FirstName = givenName
// 			log.Printf("👤 Mapped Google given_name: %s", givenName)
// 		}
// 		if familyName, ok := userInfo["family_name"].(string); ok {
// 			user.PersonalInfo.LastName = familyName
// 			log.Printf("👤 Mapped Google family_name: %s", familyName)
// 		}
// 		if name, ok := userInfo["name"].(string); ok {
// 			user.Username = name
// 			log.Printf("👤 Mapped Google name: %s", name)
// 		}

// 	case "facebook":
// 		if name, ok := userInfo["name"].(string); ok {
// 			user.Username = name
// 			log.Printf("👤 Mapped Facebook name: %s", name)
// 		}
// 		if firstName, ok := userInfo["first_name"].(string); ok {
// 			user.PersonalInfo.FirstName = firstName
// 			log.Printf("👤 Mapped Facebook first_name: %s", firstName)
// 		}
// 		if lastName, ok := userInfo["last_name"].(string); ok {
// 			user.PersonalInfo.LastName = lastName
// 			log.Printf("👤 Mapped Facebook last_name: %s", lastName)
// 		}
// 	}

// 	log.Printf("✅ User entity mapping completed: %s (%s)", user.Username, user.Email)
// 	return user
// }

// // GetJWTService returns the JWT service instance
// func (o *OAuthService) GetJWTService() *JWTService {
// 	return o.jwtService
// }
