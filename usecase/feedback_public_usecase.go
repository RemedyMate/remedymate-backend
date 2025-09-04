package usecase

import (
	"context"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
)

type PublicFeedbackUsecaseImpl struct {
	repo interfaces.FeedbackRepository
}

func NewPublicFeedbackUsecase(repo interfaces.FeedbackRepository) interfaces.PublicFeedbackUsecase {
	return &PublicFeedbackUsecaseImpl{repo: repo}
}

func (uc *PublicFeedbackUsecaseImpl) Create(ctx context.Context, in dto.CreateFeedbackDTO) (*entities.Feedback, error) {
	f := &entities.Feedback{
		SessionID: in.SessionID,
		TopicKey:  in.TopicKey,
		Language:  in.Language,
		Rating:    in.Rating,
		Message:   in.Message,
	}
	if err := uc.repo.Create(ctx, f); err != nil { return nil, err }
	return f, nil
}
