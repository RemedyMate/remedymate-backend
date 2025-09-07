package controllers

import (
	"errors"
	AppError "remedymate-backend/domain/AppError"

	"github.com/gin-gonic/gin"
)

func HandleHTTPError(c *gin.Context, err error) {
	// user authentication
	switch {
	case errors.Is(err, AppError.ErrUserNotAuthenticated):
		c.JSON(401, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrDuplicateUsername):
		c.JSON(409, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrInvalidInput):
		c.JSON(400, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrTopicAlreadyExists):
		c.JSON(409, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrInvalidActivationToken):
		c.JSON(400, gin.H{"error": err.Error()})
	// topic-related
	case errors.Is(err, AppError.ErrTopicNotFound):
		c.JSON(404, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrLanguageNotAvailable):
		c.JSON(400, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrUnsupportedLanguage):
		c.JSON(400, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrNoTopicMapped):
		c.JSON(404, gin.H{"error": err.Error()})

	// user-related
	case errors.Is(err, AppError.ErrUserNotFound):
		c.JSON(404, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrUserStatusNotFound):
		c.JSON(404, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrEmailAlreadyExist):
		c.JSON(409, gin.H{"error": err.Error()}) // conflict
	case errors.Is(err, AppError.ErrIncorrectPassword):
		c.JSON(401, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrInactiveAccount):
		c.JSON(403, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrInvalidToken):
		c.JSON(401, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrActivationFailed):
		c.JSON(400, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrVerificationFailed):
		c.JSON(400, gin.H{"error": err.Error()})
	case errors.Is(err, AppError.ErrEmailSendFailed):
		c.JSON(500, gin.H{"error": err.Error()})

	// token-related
	case errors.Is(err, AppError.ErrRefreshTokenNotFound):
		c.JSON(404, gin.H{"error": err.Error()})

	// server errors
	case errors.Is(err, AppError.ErrInternalServer):
		c.JSON(500, gin.H{"error": err.Error()})

	// fallback
	default:
		c.JSON(500, gin.H{"error": "Internal Server Error"})
	}
}
