package main

import (
	"os"

	"github.com/RemedyMate/remedymate-backend/delivery/controllers"
	"github.com/RemedyMate/remedymate-backend/delivery/routers"
	"github.com/RemedyMate/remedymate-backend/infrastructure/database"
	"github.com/RemedyMate/remedymate-backend/repository"
	"github.com/RemedyMate/remedymate-backend/usecase"
)

func main() {
	//Load environment adn connect to mongo
	database.ConnectMongo()

	//Initialize repository and usecase
	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo)
	userController := controllers.NewUserController(userUsecase)

	//Setup router
	r := routers.SetupRouter(userController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}