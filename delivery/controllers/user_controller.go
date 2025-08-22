package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/RemedyMate/remedymate-backend/domain/dto"
	"github.com/RemedyMate/remedymate-backend/domain/entities"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase interfaces.IUserUsecase
}

func NewUserController(usecase interfaces.IUserUsecase) *UserController{
	return &UserController{UserUsecase: usecase}
}

//POST /users/register
func (uc *UserController) Register(c *gin.Context){
	var input dto.RegisterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// Map DTO -> Entity
    user := entities.User{
        Username:        input.Username,
        Email:           input.Email,
        Password:        input.Password,
        PersonalInfo: entities.PersonalInfo{
            FirstName: input.PersonalInfo.FirstName,
            LastName:  input.PersonalInfo.LastName,
            Age:       input.PersonalInfo.Age,
            Gender:    input.PersonalInfo.Gender,
        },
        HealthConditions: input.HealthConditions,
        IsVerified:       false,            // default
        IsProfileFull:    false,            // default
        OAuthProviders: []entities.OAuthProvider{
            {Provider: "google", ID: ""},
        },
        RefreshToken: "",
        IsActive:     false,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
        LastLogin:    time.Now(),
    }

	if err := uc.UserUsecase.RegisterUser(context.Background(), user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}