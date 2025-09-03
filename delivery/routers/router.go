package routers

import (
	"remedymate-backend/delivery/controllers"
	"remedymate-backend/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all application routes

func SetupRouter(
	// oauthController *controllers.OAuthController,
	authController *controllers.AuthController,
	// userController *controllers.UserController,
	remedyMateController *controllers.RemedyMateController,
	conversationController *controllers.ConversationController) *gin.Engine {

	r := gin.Default()

	// API version 1
	v1 := r.Group("/api/v1")
	{
		// Remedy route which comprises /triage, /map_topic, and /compose
		v1.POST("/remedy", remedyMateController.GetRemedy)
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.Refresh)

			// Protected routes (require authentication)
			protected := auth.Group("/")
			protected.Use(middleware.AuthMiddleware())
			{
				protected.POST("/logout", authController.Logout)
				// protected.POST("/change-password", authController.ChangePassword)
			}
		}

		// Protected API routes (require authentication)
		protectedAPI := v1.Group("/")
		protectedAPI.Use(middleware.AuthMiddleware())
		{
			// superadmin routes
			protectedAPI.POST("/register", middleware.SuperAdminMiddleware(), authController.Register)
			// 	// User profile routes
			// 	users := protectedAPI.Group("/users")
			// 	{
			// 		users.GET("/profile", userController.GetProfile)
			// 		users.PUT("/profile", userController.UpdateProfile)
			// 		users.PATCH("/profile", userController.EditProfile)
			// 		users.DELETE("/profile", userController.DeleteProfile)
			// 	}
		}
	}

	// Conversation routes (public access, no authentication required)
	conversation := v1.Group("/conversation")
	{
		// Unified conversation endpoint (handles both start and continue)
		conversation.POST("/", conversationController.HandleConversation)
	}

	return r
}
