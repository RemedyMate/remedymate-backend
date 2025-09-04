package routers

import (
	"remedymate-backend/delivery/controllers"
	"remedymate-backend/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all application routes

func SetupRouter(oauthController *controllers.OAuthController,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	remedyMateController *controllers.RemedyMateController,
	conversationController *controllers.ConversationController) *gin.Engine {

	r := gin.Default()

	// Add CORS middleware for OAuth callbacks
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API version 1
	v1 := r.Group("/api/v1")
	{
		// Remedy route which comprises /triage, /map_topic, and /compose
		v1.POST("/remedy", remedyMateController.GetRemedy)
		// Authentication routes
		auth := v1.Group("/auth")
		{
			// Local authentication (no middleware required)
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)

			// OAuth routes (no middleware required)
			oauth := auth.Group("/oauth")
			{
				oauth.GET("/:provider/url", oauthController.GetAuthURL)
				oauth.GET("/:provider/callback", oauthController.HandleCallback)
				oauth.POST("/:provider/callback", oauthController.HandleCallback)
				oauth.POST("/validate", oauthController.ValidateToken)
			}

			// Protected routes (require authentication)
			protected := auth.Group("/")
			protected.Use(middleware.AuthMiddleware())
			{
				protected.POST("/logout", authController.Logout)
				protected.POST("/change-password", authController.ChangePassword)
			}
		}

		// Protected API routes (require authentication)
		protectedAPI := v1.Group("/")
		protectedAPI.Use(middleware.AuthMiddleware())
		{
			// User profile routes
			users := protectedAPI.Group("/users")
			{
				users.GET("/profile", userController.GetProfile)
				users.PUT("/profile", userController.UpdateProfile)
				users.PATCH("/profile", userController.EditProfile)
				users.DELETE("/profile", userController.DeleteProfile)
			}
		}
	}

	// Conversation routes (public access, no authentication required)
	conversation := v1.Group("/conversation")
	{
		// Unified conversation endpoint (handles both start and continue)
		conversation.POST("/", conversationController.HandleConversation)
		conversation.GET("/offline-topics", conversationController.GetOfflineHealthTopics)
	}

	return r
}
