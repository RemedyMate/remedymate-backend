package controllers

import (
	"net/http"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/interfaces"

	"github.com/gin-gonic/gin"
)

type AdminRedFlagController struct {
	uc interfaces.AdminRedFlagUsecase
}

func NewAdminRedFlagController(uc interfaces.AdminRedFlagUsecase) *AdminRedFlagController {
	return &AdminRedFlagController{uc: uc}
}

func (c *AdminRedFlagController) List(ctx *gin.Context) {
	items, err := c.uc.List(ctx.Request.Context())
	if err != nil {
		HandleHTTPError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"items": items})
}

func (c *AdminRedFlagController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	item, err := c.uc.Get(ctx.Request.Context(), id)
	if err != nil {
		HandleHTTPError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *AdminRedFlagController) Create(ctx *gin.Context) {
	var in dto.CreateRedFlagDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	actor := ctx.GetString("userID")
	item, err := c.uc.Create(ctx.Request.Context(), in, actor)
	if err != nil {
		HandleHTTPError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, item)
}

func (c *AdminRedFlagController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var in dto.UpdateRedFlagDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	actor := ctx.GetString("userID")
	item, err := c.uc.Update(ctx.Request.Context(), id, in, actor)
	if err != nil {
		HandleHTTPError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *AdminRedFlagController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	actor := ctx.GetString("userID")
	if err := c.uc.Delete(ctx.Request.Context(), id, actor); err != nil {
		HandleHTTPError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
