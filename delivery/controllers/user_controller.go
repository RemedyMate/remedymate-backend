package controllers

import (
	"log"
	"net/http"

	"github.com/RemedyMate/remedymate-backend/domain/dto"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase interfaces.IUserUsecase
}

func NewUserController(usecase interfaces.IUserUsecase) *UserController {
	return &UserController{UserUsecase: usecase}
}

// GetProfile retrieves user profile
// GET /api/v1/users/profile
func (uc *UserController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	log.Printf("👤 Profile request for user: %s", userID)

	profile, err := uc.UserUsecase.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		log.Printf("❌ Get profile failed for user %s: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("✅ Profile retrieved for user: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"profile": profile,
	})
}

// UpdateProfile updates user profile (basic update)
// PUT /api/v1/users/profile
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var updateData dto.UpdateProfileDTO
	if err := c.ShouldBindJSON(&updateData); err != nil {
		log.Printf("❌ Invalid update profile request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	log.Printf("🔄 Update profile request for user: %s", userID)

	profile, err := uc.UserUsecase.UpdateProfile(c.Request.Context(), userID.(string), updateData)
	if err != nil {
		log.Printf("❌ Update profile failed for user %s: %v", userID, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("✅ Profile updated for user: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": profile,
	})
}

// EditProfile edits user profile (comprehensive edit)
// PATCH /api/v1/users/profile
func (uc *UserController) EditProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var editData dto.EditProfileDTO
	if err := c.ShouldBindJSON(&editData); err != nil {
		log.Printf("❌ Invalid edit profile request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✏️ Edit profile request for user: %s", userID)

	profile, err := uc.UserUsecase.EditProfile(c.Request.Context(), userID.(string), editData)
	if err != nil {
		log.Printf("❌ Edit profile failed for user %s: %v", userID, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("✅ Profile edited for user: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile edited successfully",
		"profile": profile,
	})
}

// DeleteProfile soft deletes user profile
// DELETE /api/v1/users/profile
func (uc *UserController) DeleteProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var deleteData dto.DeleteProfileDTO
	if err := c.ShouldBindJSON(&deleteData); err != nil {
		log.Printf("❌ Invalid delete profile request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	log.Printf("🗑️ Delete profile request for user: %s", userID)

	if err := uc.UserUsecase.DeleteProfile(c.Request.Context(), userID.(string), deleteData); err != nil {
		log.Printf("❌ Delete profile failed for user %s: %v", userID, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("✅ Profile deleted for user: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile deleted successfully",
	})
}
