package middleware

import (
	"log"
	"net/http"
	"strings"

	"remedymate-backend/domain/entities"
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
		c.Set("role", claims.Role)

		// Continue to next handler
		c.Next()
	}
}

// SuperAdminMiddleware ensures the user has superadmin role
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("🔒 SuperAdmin middleware processing request: %s %s", c.Request.Method, c.Request.URL.Path)

		role, exists := c.Get("role")
		if !exists || role != entities.RoleSuperAdmin {
			log.Printf("❌ Access denied. User does not have superadmin role")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied. Superadmin role required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}