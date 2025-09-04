package interfaces

import (
	"context"
	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
)

type PublicFeedbackUsecase interface {
	Create(ctx context.Context, in dto.CreateFeedbackDTO) (*entities.Feedback, error)
}
