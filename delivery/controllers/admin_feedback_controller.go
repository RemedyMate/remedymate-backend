package controllers

import (
	"net/http"
	"strconv"

	"remedymate-backend/domain/interfaces"

	"github.com/gin-gonic/gin"
)

type AdminFeedbackController struct {
	uc interfaces.AdminFeedbackUsecase
}

func NewAdminFeedbackController(uc interfaces.AdminFeedbackUsecase) *AdminFeedbackController {
	return &AdminFeedbackController{uc: uc}
}

func (c *AdminFeedbackController) List(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	language := ctx.Query("language")
	items, total, err := c.uc.List(ctx.Request.Context(), limit, offset, language)
	if err != nil {
		HandleHTTPError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"items": items, "total": total})
}

func (c *AdminFeedbackController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	item, err := c.uc.Get(ctx.Request.Context(), id)
	if err != nil {
		HandleHTTPError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *AdminFeedbackController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.uc.Delete(ctx.Request.Context(), id); err != nil {
		HandleHTTPError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
