package controllers

import (
	"net/http"
	"time"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"

	"github.com/gin-gonic/gin"
)

// ProfileController handles user profile and onboarding-related HTTP requests
type ProfileController struct {
	profileUsecase interfaces.IProfileUsecase
}

// NewProfileController creates a new Profile controller instance
func NewProfileController(profileUsecase interfaces.IProfileUsecase) *ProfileController {
	return &ProfileController{
		profileUsecase: profileUsecase,
	}
}

// GetOnboardingStatus retrieves the user's onboarding progress
func (pc *ProfileController) GetOnboardingStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	status, err := pc.profileUsecase.GetOnboardingStatus(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"percent_complete": status.PercentComplete,
		"next_step":        status.NextStep,
		"sections":         status.Sections,
	})
}

// PostConsent records user consents
func (pc *ProfileController) PostConsent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req dto.ConsentPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Convert to ConsentEnvelope
	envelope := &entities.ConsentEnvelope{
		Baseline:    []entities.ConsentItem{},
		Optional:    []entities.ConsentItem{},
		LastUpdated: time.Now(),
	}
	for _, item := range req.Consents {
		consentItem := entities.ConsentItem{
			Type:      item.Type,
			Granted:   item.Granted,
			Version:   item.Version,
			Timestamp: time.Now(),
		}
		if item.Purpose == "baseline" {
			envelope.Baseline = append(envelope.Baseline, consentItem)
		} else {
			envelope.Optional = append(envelope.Optional, consentItem)
		}
	}

	err := pc.profileUsecase.PostConsent(c.Request.Context(), userID.(string), envelope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Consent recorded successfully",
	})
}

// PatchConsent updates or revokes consents
func (pc *ProfileController) PatchConsent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req dto.ConsentPatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// For each consent where granted=false, append revocation
	for _, consent := range req.Consents {
		if !consent.Granted {
			err := pc.profileUsecase.PatchConsent(c.Request.Context(), userID.(string), consent.Type, consent.RevocationReason)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Consents updated successfully",
	})
}

// PutDemographics updates user demographics
func (pc *ProfileController) PutDemographics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var demographics entities.Demographics
	if err := c.ShouldBindJSON(&demographics); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err := pc.profileUsecase.PutDemographics(c.Request.Context(), userID.(string), &demographics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Demographics updated successfully",
	})
}

// PutMedicalHistory updates user medical history
func (pc *ProfileController) PutMedicalHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var history entities.MedicalHistory
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err := pc.profileUsecase.PutMedicalHistory(c.Request.Context(), userID.(string), &history)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Medical history updated successfully",
	})
}

// PutLifestyle updates user lifestyle data
func (pc *ProfileController) PutLifestyle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var lifestyle entities.Lifestyle
	if err := c.ShouldBindJSON(&lifestyle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err := pc.profileUsecase.PutLifestyle(c.Request.Context(), userID.(string), &lifestyle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lifestyle updated successfully",
	})
}
