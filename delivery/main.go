package main

import (
	"log"
	"os"

	"github.com/RemedyMate/remedymate-backend/config"
	"github.com/RemedyMate/remedymate-backend/delivery/controllers"
	"github.com/RemedyMate/remedymate-backend/delivery/routers"
	"github.com/RemedyMate/remedymate-backend/infrastructure/auth"
	"github.com/RemedyMate/remedymate-backend/infrastructure/database"
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
	database.ConnectMongo() // commented because we don't need it currently

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

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo)
	oauthUsecase := usecase.NewOAuthUsecase(oauthService, oauthRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo, passwordService, jwtService)

	// Initialize controllers
	oauthController := controllers.NewOAuthController(oauthUsecase)
	authController := controllers.NewAuthController(authUsecase, userUsecase) // Added userUsecase
	userController := controllers.NewUserController(userUsecase)              // Re-added for profile management

	// Setup router
	r := routers.SetupRouter(oauthController, authController, userController) // Added userController back

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("✅ OAuth endpoints: /api/v1/auth/oauth/*")
	log.Printf("✅ Login endpoint: /api/v1/auth/login")
	log.Printf("✅ Protected endpoints: /api/v1/auth/* (with JWT)")
	log.Fatal(r.Run(":" + port))
}
