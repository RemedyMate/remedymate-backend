package main

import (
	"log"
	"os"

	"remedymate-backend/config"
	"remedymate-backend/delivery/controllers"
	"remedymate-backend/delivery/routers"
	"remedymate-backend/domain/dto"
	"remedymate-backend/infrastructure/auth"
	"remedymate-backend/infrastructure/content"
	"remedymate-backend/infrastructure/conversation"
	"remedymate-backend/infrastructure/database"
	"remedymate-backend/infrastructure/guidance"
	"remedymate-backend/infrastructure/llm"
	"remedymate-backend/infrastructure/triage"
	"remedymate-backend/repository"
	"remedymate-backend/usecase"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect to MongoDB
	database.ConnectMongo()

	// Load OAuth configuration
	oauthConfig := config.LoadOAuthConfig()
	if err := oauthConfig.ValidateConfig(); err != nil {
		log.Fatalf("OAuth configuration error: %v", err)
	}

	// Initialize services
	jwtService := auth.NewJWTService()
	oauthService := auth.NewOAuthService(oauthConfig, jwtService)
	passwordService := auth.NewPasswordService()

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	oauthRepo := repository.NewOAuthRepository(database.GetCollection("users"))
	conversationRepo := repository.NewConversationRepository(database.GetCollection("conversation"))
	remedyRepo := repository.NewGeminiRemedyRepo(os.Getenv("GEMINI_API_KEY"), os.Getenv("GEMINI_MODEL"))

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo)
	oauthUsecase := usecase.NewOAuthUsecase(oauthService, oauthRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo, passwordService, jwtService)

	// Initialize RemedyMate services
	contentService := content.NewContentService("./data")

	// Initialize Gemini LLM client
	gemKey := os.Getenv("GEMINI_API_KEY")
	if gemKey == "" {
		log.Fatal("❌ GEMINI_API_KEY not set.")
	}

	llmConfig := dto.LLMConfig{
		APIKey:      gemKey,
		Model:       os.Getenv("GEMINI_MODEL"),
		MaxTokens:   150,
		Temperature: 0.1,
		Timeout:     30,
	}

	geminiClient := llm.NewGeminiClient(llmConfig)
	log.Printf("✅ Using Gemini LLM client (model=%s)", llmConfig.Model)

	triageService := triage.NewTriageService(contentService, geminiClient)
	guidanceComposer := guidance.NewGuidanceComposerService(contentService, geminiClient)

	conversationService := conversation.NewConversationService(geminiClient)

	// Initialize RemedyMate usecase
	remedyMateUsecase := usecase.NewRemedyMateUsecase(
		triageService,
		contentService,
		guidanceComposer,
	)
	remedyUsecase := usecase.NewRemedyUsecase(remedyRepo)

	// Initialize Conversation usecase
	conversationUsecase := usecase.NewConversationUsecase(
		conversationService,
		conversationRepo,
	)

	// Initialize controllers
	oauthController := controllers.NewOAuthController(oauthUsecase)

	authController := controllers.NewAuthController(authUsecase, userUsecase) // Added userUsecase
	userController := controllers.NewUserController(userUsecase)              // Re-added for profile management
	remedyHandler := controllers.NewRemedyHandler(remedyUsecase)
	remedyMateController := controllers.NewRemedyMateController(remedyMateUsecase)
	conversationController := controllers.NewConversationController(conversationUsecase)

	// Setup router
	r := routers.SetupRouter(oauthController, authController, userController, remedyHandler, remedyMateController, conversationController) // Added userController back

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("✅ OAuth endpoints: /api/v1/auth/oauth/*")
	log.Printf("✅ Login endpoint: /api/v1/auth/login")
	log.Printf("✅ Protected endpoints: /api/v1/auth/* (with JWT)")
	log.Printf("✅ Conversation endpoints: /api/v1/conversation/*")
	log.Fatal(r.Run(":" + port))
}
