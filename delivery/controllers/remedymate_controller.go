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
