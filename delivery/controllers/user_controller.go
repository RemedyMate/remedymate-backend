package controllers

import (
	"log"
	"net/http"
	"strconv"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/interfaces"

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
			"error": "Invalid request body",
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

// EditProfile was removed; use UpdateProfile (PUT) for updates.

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

// GetUserProfilesPaginated retrieves user profiles with pagination, filtering, and sorting (superadmin only)
// GET /api/v1/admin/users/profiles/paginated
func (uc *UserController) GetUserProfilesPaginated(c *gin.Context) {
	log.Printf("👥 Get user profiles with pagination request")

	// Parse query parameters
	var params dto.UserProfilesQueryParams

	// Parse pagination parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}
	if params.Page == 0 {
		params.Page = 1
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = limit
		}
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	// Parse other parameters
	params.Search = c.Query("search")
	params.Status = c.Query("status")
	params.Role = c.Query("role")
	params.SortBy = c.Query("sort_by")
	params.Order = c.Query("order")

	// Default sort order
	if params.Order != "asc" && params.Order != "desc" {
		params.Order = "desc"
	}

	response, err := uc.UserUsecase.GetUserProfilesWithPagination(c.Request.Context(), params)
	if err != nil {
		log.Printf("❌ Get user profiles with pagination failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("✅ User profiles with pagination retrieved successfully")
	c.JSON(http.StatusOK, response)
}
