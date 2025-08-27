package controllers

import (
	"net/http"

	"github.com/RemedyMate/remedymate-backend/domain/dto"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
	"github.com/gin-gonic/gin"
)

type RemedyMateController struct {
	remedyMateUsecase interfaces.RemedyMateUsecase
}

func NewRemedyMateController(remedyMateUsecase interfaces.RemedyMateUsecase) *RemedyMateController {
	return &RemedyMateController{
		remedyMateUsecase: remedyMateUsecase,
	}
}

// handles triage-only requests
func (rmc *RemedyMateController) GetTriage(c *gin.Context) {
	var req dto.TriageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	response, err := rmc.remedyMateUsecase.GetTriage(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Triage failed",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetContent handles content retrieval requests
func (rmc *RemedyMateController) GetContent(c *gin.Context) {
	// Get path parameter
	topicKey := c.Param("topic_key")
	if topicKey == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Topic key is required",
			Details: "Please provide a valid topic key in the URL path",
		})
		return
	}

	// Get query parameter
	language := c.Query("language")
	if language == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Language parameter is required",
			Details: "Please provide a language parameter (e.g., ?language=en or ?language=am)",
		})
		return
	}

	// Get content from usecase
	content, err := rmc.remedyMateUsecase.GetContent(c.Request.Context(), topicKey, language)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "topic '"+topicKey+"' not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Topic not found",
				Details: "The requested topic key does not exist in our approved content library",
			})
			return
		}
		
		// Check if it's a language not available error
		if err.Error() == "language '"+language+"' not available for topic '"+topicKey+"'" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Language not available",
				Details: "The requested language is not available for this topic",
			})
			return
		}

		// Check if it's an unsupported language error
		if err.Error() == "unsupported language: "+language+". Supported languages: en, am" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Unsupported language",
				Details: "Only 'en' and 'am' languages are supported",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to retrieve content",
			Details: err.Error(),
		})
		return
	}

	// Return the content in the exact format specified
	c.JSON(http.StatusOK, gin.H{
		"self_care":      content.SelfCare,
		"otc_categories": content.OTCCategories,
		"seek_care_if":   content.SeekCareIf,
		"disclaimer":     content.Disclaimer,
	})
}

