package routers

import (
	"remedymate-backend/delivery/controllers"
	"remedymate-backend/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all application routes

func SetupRouter(
	authController *controllers.AuthController,
	remedyMateController *controllers.RemedyMateController,
	conversationController *controllers.ConversationController,
	topicController *controllers.TopicController,
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
			auth.POST("/activate", authController.Activate)
			auth.GET("/verify", authController.Verify)

			// Protected routes (require authentication)
			protected := auth.Group("/")
			protected.Use(middleware.AuthMiddleware())
			{
				protected.POST("/logout", authController.Logout)
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

		// Admin routes (auth required; all users are admins per requirement)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		{
			admin.GET("/topics", topicController.ListAllTopicsHandler)
			admin.POST("/topic", topicController.CreateTopicHandler)
			admin.PUT("/topics/:topic_key", topicController.UpdateTopicHandler)
			admin.DELETE("/topics/:topic_key", topicController.DeleteTopicHandler)
			admin.GET("/topic/:topic_key", topicController.GetTopicHandler)
		}
	}

	// Conversation routes (public access, no authentication required)
	conversation := v1.Group("/conversation")
	{
		// Initial chat greeting endpoint
		conversation.POST("/init", conversationController.InitiateChat)
		// Unified conversation endpoint (handles both start and continue)
		conversation.POST("/", conversationController.HandleConversation)
		conversation.GET("/offline-topics", conversationController.GetOfflineHealthTopics)
	}

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
