package controllers

import (
	"log"
	"net/http"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/interfaces"

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

// GetRemedy handles the entire flow: triage, mapping, content retrieval
func (rmc *RemedyMateController) GetRemedy(c *gin.Context) {
	// PROCESS: Triage
	log.Print("Triage process started")
	var req dto.RemedyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	response, err := rmc.remedyMateUsecase.GetTriage(c.Request.Context(), req.Text, req.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Triage failed",
			Details: err.Error(),
		})
		return
	}
	log.Printf("Triage result: Level=%s, RedFlags=%v", response.Level, response.RedFlags)

	// PROCESS: Mapping
	log.Print("Mapping process started")
	if response.Level == "RED" {
		c.JSON(http.StatusOK, response.Message)
		return
	}
	topicKey, err := rmc.remedyMateUsecase.MapTopic(c.Request.Context(), req.Text)
	if err != nil {
		log.Printf("Error from usecase: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map symptom to topic."})
		return
	}

	// PROCESS: Content Retrieval
	log.Printf("Content retrieval process started for topic key: %s", topicKey)
	if topicKey == "" {
		log.Print("the topic key is empty")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Mapping failed",
			Details: "No topic key could be mapped from the provided symptoms",
		})
		return
	}

	// Get content from usecase
	content, err := rmc.remedyMateUsecase.GetContent(c.Request.Context(), topicKey, req.Language)
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
		if err.Error() == "language '"+req.Language+"' not available for topic '"+topicKey+"'" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Language not available",
				Details: "The requested language is not available for this topic",
			})
			return
		}

		// Check if it's an unsupported language error
		if err.Error() == "unsupported language: "+req.Language+". Supported languages: en, am" {
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

// ComposeGuidance handles guidance composition requests
// func (rmc *RemedyMateController) ComposeGuidance(c *gin.Context) {
// 	var req dto.ComposeRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
// 			Error:   "Invalid request format",
// 			Details: err.Error(),
// 		})
// 		return
// 	}

// 	response, err := rmc.remedyMateUsecase.ComposeGuidance(c.Request.Context(), req)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
// 			Error:   "Guidance composition failed",
// 			Details: err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, response)
// }
