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
	conversationController *controllers.ConversationController,
	adminRedFlagController *controllers.AdminRedFlagController,
	adminFeedbackController *controllers.AdminFeedbackController,
	feedbackPublicController *controllers.FeedbackPublicController) *gin.Engine {

	r := gin.Default()

	// API version 1
	v1 := r.Group("/api/v1")
	{
		// Remedy route which comprises /triage, /map_topic, and /compose
		v1.POST("/remedy", remedyMateController.GetRemedy)
		// Public feedback route
		v1.POST("/feedbacks", feedbackPublicController.Create)
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.Refresh)

			// Protected routes (require authentication)
			protected := auth.Group("/")
			protected.Use(middleware.AuthMiddleware())
			{
				protected.POST("/logout", authController.Logout)
			}
		}
	}

	// Conversation routes (public access, no authentication required)
	conversation := v1.Group("/conversation")
	{
		// Unified conversation endpoint (handles both start and continue)
		conversation.POST("/", conversationController.HandleConversation)

		// Legacy endpoints (for backward compatibility)
		conversation.POST("/start", conversationController.StartConversation)
		conversation.POST("/answer", conversationController.SubmitAnswer)
		conversation.GET("/:id/report", conversationController.GetReport)
	}

	// Admin routes (auth required; all users are admins per requirement)
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		// Redflags
		admin.GET("/redflags", adminRedFlagController.List)
		admin.POST("/redflags", adminRedFlagController.Create)
		admin.PUT("/redflags/:id", adminRedFlagController.Update)
		admin.GET("/redflags/:id", adminRedFlagController.Get)
		admin.DELETE("/redflags/:id", adminRedFlagController.Delete)

		// Feedbacks
		admin.GET("/feedbacks", adminFeedbackController.List)
		admin.GET("/feedbacks/:id", adminFeedbackController.Get)
		admin.DELETE("/feedbacks/:id", adminFeedbackController.Delete)
	}

	return r
}
