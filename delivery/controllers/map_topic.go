package controllers

import (
	"log"
	"net/http"

	"remedymate-backend/domain/dto"
	"remedymate-backend/usecase"

	"github.com/gin-gonic/gin"
)

type RemedyHandler struct {
	remedyUsecase *usecase.RemedyUsecase
}

func NewRemedyHandler(uc *usecase.RemedyUsecase) *RemedyHandler {
	return &RemedyHandler{remedyUsecase: uc}
}

// MapTopic is the Gin handler function for the POST /map-topic endpoint.
func (h *RemedyHandler) MapTopic(c *gin.Context) {
	var req dto.MapTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	topicKey, err := h.remedyUsecase.MapTopic(ctx, req.UserInput)
	if err != nil {
		log.Printf("Error from usecase: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map symptom to topic."})
		return
	}

	resp := dto.MapTopicResponse{TopicKey: topicKey}
	c.JSON(http.StatusOK, resp)
}
