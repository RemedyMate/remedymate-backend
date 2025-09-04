package controllers

import (
	"net/http"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/interfaces"

	"github.com/gin-gonic/gin"
)

type FeedbackPublicController struct {
	uc interfaces.PublicFeedbackUsecase
}

func NewFeedbackPublicController(uc interfaces.PublicFeedbackUsecase) *FeedbackPublicController {
	return &FeedbackPublicController{uc: uc}
}

func (c *FeedbackPublicController) Create(ctx *gin.Context) {
	var in dto.CreateFeedbackDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	f, err := c.uc.Create(ctx.Request.Context(), in)
	if err != nil { HandleHTTPError(ctx, err); return }
	ctx.JSON(http.StatusCreated, f)
}
