package routers

import (
	"github.com/RemedyMate/remedymate-backend/delivery/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(UserController *controllers.UserController) *gin.Engine{
	r := gin.Default()

	userGroup := r.Group("/api/v1/auth")
	{
		userGroup.POST("/register", UserController.Register)
		//login, profile, password reset routes wiill be added here
	}

	return r
}