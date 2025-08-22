package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RemedyMate/remedymate-backend/domain/dto"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
	"github.com/gin-gonic/gin"
)

// OAuthController handles OAuth-related HTTP requests
type OAuthController struct {
	oauthUsecase interfaces.IOAuthUsecase
}

// NewOAuthController creates a new OAuth controller instance
func NewOAuthController(oauthUsecase interfaces.IOAuthUsecase) *OAuthController {
	return &OAuthController{
		oauthUsecase: oauthUsecase,
	}
}

// GetAuthURL generates the OAuth authorization URL
// GET /api/v1/auth/oauth/:provider/url
func (oc *OAuthController) GetAuthURL(c *gin.Context) {
	provider := c.Param("provider")

	// Validate provider
	validProviders := map[string]bool{
		"google":   true,
		"facebook": true,
	}

	if !validProviders[provider] {
		c.JSON(http.StatusBadRequest, dto.OAuthErrorDTO{
			Error:       "Invalid provider",
			Description: "Supported providers: google, github, facebook",
			Provider:    provider,
		})
		return
	}

	response, err := oc.oauthUsecase.GetAuthURL(c.Request.Context(), provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.OAuthErrorDTO{
			Error:       "Failed to generate auth URL",
			Description: err.Error(),
			Provider:    provider,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// HandleCallback processes the OAuth callback (both GET and POST)
func (oc *OAuthController) HandleCallback(c *gin.Context) {
	provider := c.Param("provider")
	log.Printf("🚀 OAuth callback initiated for provider: %s", provider)

	// Validate provider
	validProviders := map[string]bool{
		"google":   true,
		"facebook": true,
	}

	if !validProviders[provider] {
		log.Printf("❌ Invalid OAuth provider: %s", provider)
		c.JSON(http.StatusBadRequest, dto.OAuthErrorDTO{
			Error:       "Invalid provider",
			Description: "Supported providers: google, facebook",
			Provider:    provider,
		})
		return
	}

	var callback dto.OAuthCallbackDTO

	// Handle both GET and POST requests
	if c.Request.Method == "GET" {
		// Extract from query parameters for GET requests (OAuth redirects)
		callback.Code = c.Query("code")
		callback.State = c.Query("state")
		log.Printf(" GET request - Code: %s, State: %s", callback.Code[:10]+"...", callback.State[:10]+"...")
	} else {
		// Extract from JSON body for POST requests
		if err := c.ShouldBindJSON(&callback); err != nil {
			log.Printf("❌ Failed to parse POST request body: %v", err)
			c.JSON(http.StatusBadRequest, dto.OAuthErrorDTO{
				Error:       "Invalid request body",
				Description: err.Error(),
				Provider:    provider,
			})
			return
		}
		log.Printf("📥 POST request - Code: %s, State: %s", callback.Code[:10]+"...", callback.State[:10]+"...")
	}

	// Validate required fields
	if callback.Code == "" {
		log.Printf("❌ Missing authorization code for provider: %s", provider)
		c.JSON(http.StatusBadRequest, dto.OAuthErrorDTO{
			Error:       "Missing authorization code",
			Description: "Authorization code is required",
			Provider:    provider,
		})
		return
	}

	log.Printf("🔄 Processing OAuth callback for provider: %s", provider)
	response, err := oc.oauthUsecase.HandleCallback(c.Request.Context(), provider, callback)
	if err != nil {
		log.Printf("❌ OAuth callback failed for provider %s: %v", provider, err)
		c.JSON(http.StatusInternalServerError, dto.OAuthErrorDTO{
			Error:       "Authentication failed",
			Description: err.Error(),
			Provider:    provider,
		})
		return
	}

	log.Printf("✅ OAuth callback successful for provider: %s, User ID: %s", provider, response.User)

	// For GET requests (OAuth redirects), return JSON response
	if c.Request.Method == "GET" {
		log.Printf("📤 Returning JSON response for GET OAuth callback")
		c.JSON(http.StatusOK, gin.H{
			"message":      "OAuth authentication successful",
			"user":         response.User,
			"access_token": response.AccessToken,
			"provider":     provider,
		})
		return
	}

	// For POST requests, return JSON response
	log.Printf("📤 Returning JSON response for POST OAuth callback")
	c.JSON(http.StatusOK, response)
}

// ValidateToken validates a JWT token
// POST /api/v1/auth/oauth/validate
func (oc *OAuthController) ValidateToken(c *gin.Context) {
	var request struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.OAuthErrorDTO{
			Error:       "Invalid request body",
			Description: err.Error(),
		})
		return
	}

	user, err := oc.oauthUsecase.ValidateToken(c.Request.Context(), request.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.OAuthErrorDTO{
			Error:       "Invalid token",
			Description: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"message": "Token is valid",
	})
}

// HandleOAuthRedirect handles the OAuth redirect from providers
func (oc *OAuthController) HandleOAuthRedirect(c *gin.Context) {
	provider := c.Param("provider")

	// Extract code and state from query parameters
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, dto.OAuthErrorDTO{
			Error:       "Missing authorization code",
			Description: "Authorization code is required",
			Provider:    provider,
		})
		return
	}

	callback := dto.OAuthCallbackDTO{
		Code:  code,
		State: state,
	}

	response, err := oc.oauthUsecase.HandleCallback(c.Request.Context(), provider, callback)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.OAuthErrorDTO{
			Error:       "Authentication failed",
			Description: err.Error(),
			Provider:    provider,
		})
		return
	}

	// Redirect to frontend with success
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	redirectURL := fmt.Sprintf("%s/oauth/success?token=%s", frontendURL, response.AccessToken)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// RefreshToken refreshes an expired access token
// POST /api/v1/auth/oauth/refresh
func (oc *OAuthController) RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.OAuthErrorDTO{
			Error:       "Invalid request body",
			Description: err.Error(),
		})
		return
	}

	response, err := oc.oauthUsecase.RefreshToken(c.Request.Context(), request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.OAuthErrorDTO{
			Error:       "Token refresh failed",
			Description: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
