package controllers

import (
	"errors"
	"log"
	"net/http"

	derrors "remedymate-backend/domain/AppError"
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
	var req dto.RemedyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	resp, err := rmc.remedyMateUsecase.GetRemedy(c.Request.Context(), req)
	if err != nil {
		log.Printf("GetRemedy error: %v\n", err)
		if errors.Is(err, derrors.ErrNoTopicMapped) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Mapping failed", Details: err.Error()})
			return
		}
		if errors.Is(err, derrors.ErrTopicNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Topic not found", Details: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to get remedy", Details: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
