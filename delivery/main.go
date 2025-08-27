package main

import (
	"log"
	"os"

	"github.com/RemedyMate/remedymate-backend/config"
	"github.com/RemedyMate/remedymate-backend/delivery/controllers"
	"github.com/RemedyMate/remedymate-backend/delivery/routers"
	"github.com/RemedyMate/remedymate-backend/domain/dto"
	"github.com/RemedyMate/remedymate-backend/infrastructure/auth"
	"github.com/RemedyMate/remedymate-backend/infrastructure/content"
	"github.com/RemedyMate/remedymate-backend/infrastructure/database"
	"github.com/RemedyMate/remedymate-backend/infrastructure/guidance"
	"github.com/RemedyMate/remedymate-backend/infrastructure/llm"
	"github.com/RemedyMate/remedymate-backend/infrastructure/triage"
	"github.com/RemedyMate/remedymate-backend/repository"
	"github.com/RemedyMate/remedymate-backend/usecase"
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
	llmRepo := repository.NewGeminiRemedyRepo(os.Getenv("GEMINI_API_KEY"), os.Getenv("GEMINI_MODEL"))

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo)
	oauthUsecase := usecase.NewOAuthUsecase(oauthService, oauthRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo, passwordService, jwtService)
	remedyUsecase := usecase.NewRemedyUsecase(llmRepo)

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

	// Initialize RemedyMate usecase
	remedyMateUsecase := usecase.NewRemedyMateUsecase(
		triageService,
		contentService,
		guidanceComposer,
	)

	// Initialize controllers
	oauthController := controllers.NewOAuthController(oauthUsecase)
	authController := controllers.NewAuthController(authUsecase, userUsecase) // Added userUsecase
	userController := controllers.NewUserController(userUsecase)              // Re-added for profile management
	remedyHandler := controllers.NewRemedyHandler(remedyUsecase)
	remedyMateController := controllers.NewRemedyMateController(remedyMateUsecase)

	// Setup router
	r := routers.SetupRouter(oauthController, authController, userController, remedyHandler, remedyMateController) // Added userController back

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := r.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
