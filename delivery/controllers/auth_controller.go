package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	authUsecase interfaces.IAuthUsecase
	userUsecase interfaces.IUserUsecase // Added for registration
}

// NewAuthController creates a new Auth controller instance
func NewAuthController(authUsecase interfaces.IAuthUsecase, userUsecase interfaces.IUserUsecase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
		userUsecase: userUsecase,
	}
}

func toPtr(input string) *string {
	return &input
}

// Register creates a new user account
func (ac *AuthController) Register(c *gin.Context) {
	var input dto.RegisterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("invalid registration request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	// Map DTO -> Entity
	user := entities.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		PersonalInfo: &entities.PersonalInfo{
			FirstName: &input.PersonalInfo.FirstName,
			LastName:  &input.PersonalInfo.LastName,
			Age:       &input.PersonalInfo.Age,
			Gender:    &input.PersonalInfo.Gender,
		},
		Role:      entities.RoleAdmin,
		CreatedBy: toPtr(c.GetString("userID")), // Set creator if available
		UpdatedBy: toPtr(c.GetString("userID")),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastLogin: time.Now(),
	}

	if err := ac.userUsecase.RegisterUser(context.Background(), user); err != nil {
		HandleHTTPError(c, err)
		return
	}

	log.Printf("✅ Registration successful for email: %s", input.Email)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

// Login authenticates a user with email and password
// POST /api/v1/auth/login
func (ac *AuthController) Login(c *gin.Context) {
	var loginData dto.LoginDTO
	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Printf("❌ Invalid login request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	log.Printf("🔐 Login request for email: %s", loginData.Email)

	response, err := ac.authUsecase.Login(c.Request.Context(), loginData)
	if err != nil {
		log.Printf("❌ Login failed for email %s: %v", loginData.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("✅ Login successful for email: %s", loginData.Email)
	c.JSON(http.StatusOK, response)
}

func (ac *AuthController) Refresh(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("invalid token refresh request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	response, err := ac.authUsecase.RefreshToken(c.Request.Context(), request.RefreshToken)
	if err != nil {
		HandleHTTPError(c, err)
		return
	}

	log.Printf("✅ Token refreshed successfully")
	c.JSON(http.StatusOK, response)
}

// Logout logs out a user
// POST /api/v1/auth/logout
func (ac *AuthController) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	log.Printf("🚪 Logout request for user: %s", userID)

	if err := ac.authUsecase.Logout(c.Request.Context(), userID.(string)); err != nil {
		log.Printf("❌ Logout failed for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Logout failed",
		})
		return
	}

	log.Printf("✅ Logout successful for user: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// ChangePassword changes a user's password
// POST /api/v1/auth/change-password
// func (ac *AuthController) ChangePassword(c *gin.Context) {
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error": "User not authenticated",
// 		})
// 		return
// 	}

// 	var request struct {
// 		OldPassword string `json:"old_password" binding:"required"`
// 		NewPassword string `json:"new_password" binding:"required,min=6"`
// 	}

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		log.Printf("❌ Invalid change password request: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "Invalid request body",
// 			"details": err.Error(),
// 		})
// 		return
// 	}

// 	log.Printf("🔐 Password change request for user: %s", userID)

// 	if err := ac.authUsecase.ChangePassword(c.Request.Context(), userID.(string), request.OldPassword, request.NewPassword); err != nil {
// 		log.Printf("❌ Password change failed for user %s: %v", userID, err)
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	log.Printf("✅ Password changed successfully for user: %s", userID)
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Password changed successfully",
// 	})
// }
