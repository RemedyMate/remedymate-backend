package middleware

import (
	"log"
	"net/http"
	"strings"

	jwtutil "remedymate-backend/util/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("🔒 Auth middleware processing request: %s %s", c.Request.Method, c.Request.URL.Path)

		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("❌ No Authorization header found")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("❌ Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Use 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			log.Printf("❌ Empty token in Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtutil.ValidateToken(token, true)
		if err != nil {
			log.Printf("❌ JWT token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		// Continue to next handler
		c.Next()
	}
}

// OptionalAuthMiddleware allows requests with or without valid tokens
// func OptionalAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		log.Printf("🔓 Optional auth middleware processing request: %s %s", c.Request.Method, c.Request.URL.Path)

// 		// Get Authorization header
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			log.Printf("ℹ️ No Authorization header, continuing without authentication")
// 			c.Next()
// 			return
// 		}

// 		// Check if header starts with "Bearer "
// 		if !strings.HasPrefix(authHeader, "Bearer ") {
// 			log.Printf("ℹ️ Invalid Authorization header format, continuing without authentication")
// 			c.Next()
// 			return
// 		}

// 		// Extract token
// 		token := strings.TrimPrefix(authHeader, "Bearer ")
// 		if token == "" {
// 			log.Printf("ℹ️ Empty token, continuing without authentication")
// 			c.Next()
// 			return
// 		}

// 		log.Printf("🔍 Attempting to validate optional JWT token...")

// 		// Try to validate token
// 		claims, err := auth.NewJWTService().ValidateToken(token)
// 		if err != nil {
// 			log.Printf("ℹ️ Optional JWT token validation failed: %v", err)
// 			c.Next()
// 			return
// 		}

// 		log.Printf("✅ Optional JWT token validated for user: %s (%s)", claims.Username, claims.UserID)

// 		// Set user information in context if token is valid
// 		c.Set("userID", claims.UserID)
// 		c.Set("username", claims.Username)
// 		c.Set("email", claims.Email)

// 		// Continue to next handler
// 		c.Next()
// 	}
// }
